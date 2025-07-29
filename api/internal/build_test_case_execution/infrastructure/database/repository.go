package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/models"
	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/ports"
	dashboardModels "github.com/BennyEisner/test-results/internal/dashboard/domain/models"
)

// SQLBuildTestCaseExecutionRepository implements BuildTestCaseExecutionRepository
type SQLBuildTestCaseExecutionRepository struct {
	db *sql.DB
}

// NewSQLBuildTestCaseExecutionRepository creates a new SQL repository
func NewSQLBuildTestCaseExecutionRepository(db *sql.DB) ports.BuildTestCaseExecutionRepository {
	return &SQLBuildTestCaseExecutionRepository{db: db}
}

// GetByID retrieves a build test case execution by ID
func (r *SQLBuildTestCaseExecutionRepository) GetByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error) {
	query := `SELECT id, build_id, test_case_id, status, execution_time, created_at 
			  FROM build_test_case_executions WHERE id = $1`

	var execution models.BuildTestCaseExecution
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&execution.ID,
		&execution.BuildID,
		&execution.TestCaseID,
		&execution.Status,
		&execution.ExecutionTime,
		&execution.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get execution by ID: %w", err)
	}

	return &execution, nil
}

// GetAllByBuildID retrieves all build test case executions for a build
func (r *SQLBuildTestCaseExecutionRepository) GetAllByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error) {
	query := `SELECT e.id, e.build_id, e.test_case_id, tc.name, tc.classname, 
			  e.status, e.execution_time, e.created_at
			  FROM build_test_case_executions e
			  JOIN test_cases tc ON e.test_case_id = tc.id
			  WHERE e.build_id = $1`

	rows, err := r.db.QueryContext(ctx, query, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions by build ID: %w", err)
	}
	defer rows.Close()

	var executions []*models.BuildExecutionDetail
	for rows.Next() {
		var execution models.BuildExecutionDetail
		err := rows.Scan(
			&execution.ExecutionID,
			&execution.BuildID,
			&execution.TestCaseID,
			&execution.TestCaseName,
			&execution.ClassName,
			&execution.Status,
			&execution.ExecutionTime,
			&execution.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}
		executions = append(executions, &execution)
	}

	return executions, nil
}

// Create creates a new build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Create(ctx context.Context, execution *models.BuildTestCaseExecution) error {
	query := `INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time)
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		execution.BuildID,
		execution.TestCaseID,
		execution.Status,
		execution.ExecutionTime,
	).Scan(&execution.ID, &execution.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create execution: %w", err)
	}

	return nil
}

// Update updates an existing build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Update(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error) {
	query := `UPDATE build_test_case_executions 
			  SET build_id = $1, test_case_id = $2, status = $3, execution_time = $4
			  WHERE id = $5 RETURNING id, build_id, test_case_id, status, execution_time, created_at`

	var updatedExecution models.BuildTestCaseExecution
	err := r.db.QueryRowContext(ctx, query,
		execution.BuildID,
		execution.TestCaseID,
		execution.Status,
		execution.ExecutionTime,
		id,
	).Scan(
		&updatedExecution.ID,
		&updatedExecution.BuildID,
		&updatedExecution.TestCaseID,
		&updatedExecution.Status,
		&updatedExecution.ExecutionTime,
		&updatedExecution.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	return &updatedExecution, nil
}

// Delete deletes a build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM build_test_case_executions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete execution: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("execution not found")
	}

	return nil
}

// GetMetric returns a metric for a project
func (r *SQLBuildTestCaseExecutionRepository) GetMetric(ctx context.Context, projectID int64, metricType string) (*dashboardModels.MetricCardDTO, error) {
	var query string
	var args []interface{}

	switch metricType {
	case "pass_rate":
		query = `
			SELECT 
				COALESCE(SUM(CASE WHEN btce.status = 'passed' THEN 1 ELSE 0 END) * 100.0 / COUNT(btce.id), 0)
			FROM build_test_case_executions btce
			JOIN builds b ON btce.build_id = b.id
			JOIN test_suites ts ON b.test_suite_id = ts.id
			WHERE ts.project_id = $1
		`
		args = append(args, projectID)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}

	var value float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return &dashboardModels.MetricCardDTO{Title: "Pass Rate", Value: "0%"}, nil
		}
		return nil, fmt.Errorf("failed to get metric: %w", err)
	}

	return &dashboardModels.MetricCardDTO{
		Title: "Pass Rate",
		Value: fmt.Sprintf("%.2f%%", value),
	}, nil
}

// GetChartData returns data for a chart
func (r *SQLBuildTestCaseExecutionRepository) GetChartData(ctx context.Context, projectID int64, chartType string, suiteID *int64, buildID *int64, limit *int) (*dashboardModels.DataChartDTO, error) {
	var baseQuery string
	var groupBy string
	var orderBy string
	args := []interface{}{projectID}
	paramIndex := 2

	limitVal := 15 // A more reasonable default limit
	if limit != nil {
		limitVal = *limit
	}

	switch chartType {
	case "bar":
		baseQuery = `
            SELECT
                tc.name as label,
                COUNT(btce.id) as value
            FROM build_test_case_executions btce
            JOIN test_cases tc ON btce.test_case_id = tc.id
            JOIN builds b ON btce.build_id = b.id
            JOIN test_suites ts ON b.test_suite_id = ts.id
            WHERE ts.project_id = $1
        `
		groupBy = "GROUP BY tc.name"
		orderBy = "ORDER BY value DESC"
	case "build-duration":
		baseQuery = `
            WITH ranked_builds AS (
                SELECT
                    b.id::text as label,
                    b.duration as value,
                    ts.id as suite_id,
                    ROW_NUMBER() OVER(PARTITION BY ts.id ORDER BY b.created_at DESC) as rn
                FROM builds b
                JOIN test_suites ts ON b.test_suite_id = ts.id
                WHERE ts.project_id = $1
            )
            SELECT label, value FROM ranked_builds
        `
		orderBy = "ORDER BY suite_id, rn"
	case "line", "pass-fail-trend":
		baseQuery = `
            SELECT
                DATE(b.created_at)::text as date,
                SUM(CASE WHEN btce.status = 'passed' THEN 1 ELSE 0 END) as passed,
                SUM(CASE WHEN btce.status = 'failed' THEN 1 ELSE 0 END) as failed
            FROM build_test_case_executions btce
            JOIN builds b ON btce.build_id = b.id
            JOIN test_suites ts ON b.test_suite_id = ts.id
            WHERE ts.project_id = $1
        `
		groupBy = "GROUP BY DATE(b.created_at)"
		orderBy = "ORDER BY DATE(b.created_at)"
	case "test-case-pass-rate":
		if buildID != nil {
			baseQuery = `
                SELECT
                    tc.name as label,
                    CASE WHEN e.status = 'passed' THEN 100.0 ELSE 0.0 END as value
                FROM build_test_case_executions e
                JOIN test_cases tc ON e.test_case_id = tc.id
                JOIN builds b ON e.build_id = b.id
                WHERE b.id = $1
            `
			args = []interface{}{*buildID}
			paramIndex = 2
		} else {
			baseQuery = `
                SELECT
                    tc.name as label,
                    (SUM(CASE WHEN e.status = 'passed' THEN 1 ELSE 0 END) * 100.0 / COUNT(e.id)) as value
                FROM build_test_case_executions e
                JOIN test_cases tc ON e.test_case_id = tc.id
                JOIN builds b ON e.build_id = b.id
                JOIN test_suites ts ON b.test_suite_id = ts.id
                WHERE ts.project_id = $1
            `
			groupBy = "GROUP BY tc.name"
			orderBy = "ORDER BY value DESC"
		}
	default:
		return nil, fmt.Errorf("unknown chart type: %s", chartType)
	}

	var conditions string
	if chartType == "build-duration" {
		conditions = fmt.Sprintf(" WHERE rn <= $%d", paramIndex)
		args = append(args, limitVal)
		paramIndex++
		if suiteID != nil {
			conditions += fmt.Sprintf(" AND suite_id = $%d", paramIndex)
			args = append(args, *suiteID)
			paramIndex++
		}
	} else if buildID == nil {
		if suiteID != nil {
			conditions += fmt.Sprintf(" AND ts.id = $%d", paramIndex)
			args = append(args, *suiteID)
			paramIndex++
		}
	}

	query := fmt.Sprintf("%s %s %s %s", baseQuery, conditions, groupBy, orderBy)

	// Apply limit for other chart types that need it
	if chartType != "build-duration" && chartType != "test-case-pass-rate" {
		query += fmt.Sprintf(" LIMIT $%d", paramIndex)
		args = append(args, limitVal)
		paramIndex++
	}

	log.Printf("Executing GetChartData query: %s", query)
	log.Printf("With args: %v", args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data: %w", err)
	}
	defer rows.Close()

	var labels []string
	var passedData []int
	var failedData []int

	datasets := []dashboardModels.DatasetDTO{}

	for rows.Next() {
		switch chartType {
		case "bar", "build-duration", "test-case-pass-rate":
			var label string
			var value float64
			if err := rows.Scan(&label, &value); err != nil {
				return nil, fmt.Errorf("failed to scan chart data: %w", err)
			}
			labels = append(labels, label)
			passedData = append(passedData, int(value))
		case "line", "pass-fail-trend":
			var date string
			var passed int
			var failed int
			if err := rows.Scan(&date, &passed, &failed); err != nil {
				return nil, fmt.Errorf("failed to scan chart data: %w", err)
			}
			labels = append(labels, date)
			passedData = append(passedData, passed)
			failedData = append(failedData, failed)
		}
	}

	log.Printf("GetChartData query returned %d labels", len(labels))

	// Define color schemes for different chart types
	var xAxisLabel, yAxisLabel string
	var backgroundColors, borderColors []string

	switch chartType {
	case "bar":
		xAxisLabel = "Test Cases"
		yAxisLabel = "Number of Executions"
		backgroundColors = []string{"#3B82F6"}
		borderColors = []string{"#1D4ED8"}
		datasets = append(datasets, dashboardModels.DatasetDTO{
			Label:           "Executions",
			Data:            passedData,
			BackgroundColor: backgroundColors,
			BorderColor:     borderColors,
		})
	case "build-duration":
		xAxisLabel = "Build ID"
		yAxisLabel = "Duration (seconds)"
		backgroundColors = []string{"#10B981"}
		borderColors = []string{"#059669"}
		datasets = append(datasets, dashboardModels.DatasetDTO{
			Label:           "Duration (s)",
			Data:            passedData,
			BackgroundColor: backgroundColors,
			BorderColor:     borderColors,
		})
	case "test-case-pass-rate":
		xAxisLabel = "Test Cases"
		yAxisLabel = "Pass Rate (%)"
		backgroundColors = []string{"#F59E0B"}
		borderColors = []string{"#D97706"}
		datasets = append(datasets, dashboardModels.DatasetDTO{
			Label:           "Pass Rate (%)",
			Data:            passedData,
			BackgroundColor: backgroundColors,
			BorderColor:     borderColors,
		})
	case "line", "pass-fail-trend":
		xAxisLabel = "Date"
		yAxisLabel = "Number of Tests"
		datasets = append(datasets, dashboardModels.DatasetDTO{
			Label:           "Passed",
			Data:            passedData,
			BackgroundColor: []string{"#10B981"},
			BorderColor:     []string{"#059669"},
		}, dashboardModels.DatasetDTO{
			Label:           "Failed",
			Data:            failedData,
			BackgroundColor: []string{"#EF4444"},
			BorderColor:     []string{"#DC2626"},
		})
	}

	return &dashboardModels.DataChartDTO{
		Labels:     labels,
		Datasets:   datasets,
		XAxisLabel: xAxisLabel,
		YAxisLabel: yAxisLabel,
	}, nil
}
