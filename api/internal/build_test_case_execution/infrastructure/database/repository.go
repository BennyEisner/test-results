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

// getChartQuery constructs the SQL query for a given chart type and context.
func (r *SQLBuildTestCaseExecutionRepository) getChartQuery(chartType string, projectID int64, suiteID, buildID *int64) (string, string, string, []interface{}, int) {
	var baseQuery, groupBy, orderBy string
	args := []interface{}{projectID}
	paramIndex := 2

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
		if suiteID == nil {
			baseQuery = `
                SELECT
                    ts.name as label,
                    AVG(b.duration) as value
                FROM builds b
                JOIN test_suites ts ON b.test_suite_id = ts.id
                WHERE ts.project_id = $1
            `
			groupBy = "GROUP BY ts.name"
			orderBy = "ORDER BY value DESC"
		} else {
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
		}
	case "line", "pass-fail-trend":
		baseQuery = `
            SELECT
                DATE(b.created_at)::text as date,
                SUM(CASE WHEN btce.status = 'passed' THEN 1 ELSE 0 END) as passed,
                SUM(CASE WHEN btce.status = 'failed' THEN 1 ELSE 0 END) as failed,
                SUM(CASE WHEN btce.status = 'skipped' THEN 1 ELSE 0 END) as skipped
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
                    e.status as label,
                    COUNT(e.id) as value
                FROM build_test_case_executions e
                WHERE e.build_id = $1
                GROUP BY e.status
            `
			args = []interface{}{*buildID}
			paramIndex = 2
			groupBy = ""
			orderBy = ""
		} else if suiteID == nil {
			baseQuery = `
                SELECT
                    ts.name as label,
                    (SUM(CASE WHEN e.status = 'passed' THEN 1 ELSE 0 END) * 100.0 / COUNT(e.id)) as value
                FROM build_test_case_executions e
                JOIN builds b ON e.build_id = b.id
                JOIN test_suites ts ON b.test_suite_id = ts.id
                WHERE ts.project_id = $1
            `
			groupBy = "GROUP BY ts.name"
			orderBy = "ORDER BY value DESC"
		} else {
			baseQuery = `
                SELECT
                    b.id::text as label,
                    (SUM(CASE WHEN e.status = 'passed' THEN 1 ELSE 0 END) * 100.0 / COUNT(e.id)) as value
                FROM build_test_case_executions e
                JOIN builds b ON e.build_id = b.id
                WHERE b.test_suite_id = $1
            `
			args = []interface{}{*suiteID}
			paramIndex = 2
			groupBy = "GROUP BY b.id"
			orderBy = "ORDER BY b.id DESC"
		}
	}
	return baseQuery, groupBy, orderBy, args, paramIndex
}

// GetChartData returns data for a chart
func (r *SQLBuildTestCaseExecutionRepository) GetChartData(ctx context.Context, projectID int64, chartType string, suiteID, buildID *int64, limit *int) (*dashboardModels.DataChartDTO, error) {
	limitVal := 15 // A more reasonable default limit
	if limit != nil {
		limitVal = *limit
	}

	baseQuery, groupBy, orderBy, args, paramIndex := r.getChartQuery(chartType, projectID, suiteID, buildID)
	if baseQuery == "" {
		return nil, fmt.Errorf("unknown chart type: %s", chartType)
	}

	var conditions string
	if chartType == "build-duration" && suiteID != nil {
		conditions = fmt.Sprintf(" WHERE rn <= $%d", paramIndex)
		args = append(args, limitVal)
		paramIndex++
		conditions += fmt.Sprintf(" AND suite_id = $%d", paramIndex)
		args = append(args, *suiteID)
		paramIndex++
	} else if buildID == nil && chartType != "test-case-pass-rate" {
		if suiteID != nil {
			conditions += fmt.Sprintf(" AND ts.id = $%d", paramIndex)
			args = append(args, *suiteID)
			paramIndex++
		}
	}

	query := fmt.Sprintf("%s %s %s %s", baseQuery, conditions, groupBy, orderBy)

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
	var skippedData []int
	var values []float64
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
			values = append(values, value)
		case "line", "pass-fail-trend":
			var date string
			var passed, failed, skipped int
			if err := rows.Scan(&date, &passed, &failed, &skipped); err != nil {
				return nil, fmt.Errorf("failed to scan chart data: %w", err)
			}
			labels = append(labels, date)
			passedData = append(passedData, passed)
			failedData = append(failedData, failed)
			skippedData = append(skippedData, skipped)
		}
	}

	log.Printf("GetChartData query returned %d labels", len(labels))

	datasets, xAxisLabel, yAxisLabel := r.getChartStyling(chartType, passedData, failedData, skippedData, values, labels)

	return &dashboardModels.DataChartDTO{
		Labels:     labels,
		Datasets:   datasets,
		XAxisLabel: xAxisLabel,
		YAxisLabel: yAxisLabel,
	}, nil
}

func (r *SQLBuildTestCaseExecutionRepository) getChartStyling(chartType string, passedData, failedData, skippedData []int, values []float64, labels []string) ([]dashboardModels.DatasetDTO, string, string) {
	var xAxisLabel, yAxisLabel string
	datasets := []dashboardModels.DatasetDTO{}

	backgroundColors, borderColors := r.getDynamicColors(chartType, values, labels)

	switch chartType {
	case "bar":
		xAxisLabel = "Test Cases"
		yAxisLabel = "Number of Executions"
		datasets = append(datasets, dashboardModels.DatasetDTO{
			Label:           "Executions",
			Data:            passedData,
			BackgroundColor: backgroundColors,
			BorderColor:     borderColors,
		})
	case "build-duration":
		xAxisLabel = "Build ID"
		yAxisLabel = "Duration (seconds)"
		datasets = append(datasets, dashboardModels.DatasetDTO{
			Label:           "Duration (s)",
			Data:            passedData,
			BackgroundColor: backgroundColors,
			BorderColor:     borderColors,
		})
	case "test-case-pass-rate":
		xAxisLabel = "Test Cases"
		yAxisLabel = "Pass Rate (%)"
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
			BackgroundColor: []string{"#57F064"},
			BorderColor:     []string{"#57F064"},
		}, dashboardModels.DatasetDTO{
			Label:           "Failed",
			Data:            failedData,
			BackgroundColor: []string{"#EB4A4A"},
			BorderColor:     []string{"#EB4A4A"},
		}, dashboardModels.DatasetDTO{
			Label:           "Skipped",
			Data:            skippedData,
			BackgroundColor: []string{"#808080"},
			BorderColor:     []string{"#808080"},
		})
	}
	return datasets, xAxisLabel, yAxisLabel
}

func (r *SQLBuildTestCaseExecutionRepository) getDynamicColors(chartType string, values []float64, labels []string) ([]string, []string) {
	var backgroundColors []string
	var borderColors []string

	if len(values) == 0 {
		return []string{"#3B82F6"}, []string{"#1D4ED8"}
	}

	switch chartType {
	case "test-case-pass-rate":
		// Check if this is a build-level view by looking at the labels
		isBuildLevel := false
		for _, label := range labels {
			if label == "passed" || label == "failed" || label == "skipped" {
				isBuildLevel = true
				break
			}
		}

		if isBuildLevel {
			// Build-level view: assign static colors for pass/fail/skipped
			for _, label := range labels {
				if label == "passed" {
					backgroundColors = append(backgroundColors, "#57F064") // Green
					borderColors = append(borderColors, "#57F064")
				} else if label == "failed" {
					backgroundColors = append(backgroundColors, "#EB4A4A") // Red
					borderColors = append(borderColors, "#EB4A4A")
				} else if label == "skipped" {
					backgroundColors = append(backgroundColors, "#808080") // Grey
					borderColors = append(borderColors, "#808080")
				} else {
					// Fallback for other statuses
					backgroundColors = append(backgroundColors, "#E9EE5C") // Yellow
					borderColors = append(borderColors, "#E9EE5C")
				}
			}
		} else {
			// Project/Suite-level view: use percentage-based colors
			for _, v := range values {
				color := getColorForPercentage(v / 100.0)
				backgroundColors = append(backgroundColors, color)
				borderColors = append(borderColors, color)
			}
		}
	case "build-duration":
		minVal, maxVal := values[0], values[0]
		for _, v := range values {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}

		for _, v := range values {
			// Normalize the value to a 0-1 range (inverted, so shorter is better)
			var percentage float64
			if maxVal-minVal == 0 {
				percentage = 1.0 // All values are the same, so default to green
			} else {
				percentage = 1.0 - (v-minVal)/(maxVal-minVal)
			}
			color := getColorForPercentage(percentage)
			backgroundColors = append(backgroundColors, color)
			borderColors = append(borderColors, color)
		}
	default:
		// Default colors for other chart types
		return []string{"#3B82F6"}, []string{"#1D4ED8"}
	}

	return backgroundColors, borderColors
}

// getColorForPercentage generates a color from a multi-point gradient based on a percentage (0.0 to 1.0)
func getColorForPercentage(p float64) string {
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}

	// New Spectrum: Red (0.0) -> Orange (0.5) -> Yellow (0.75) -> Green (1.0)
	// Colors in RGB:
	red := [3]int{235, 74, 74}     // #EB4A4A
	orange := [3]int{235, 130, 74} // #EB824A
	yellow := [3]int{255, 248, 82} // #E9EE5C
	green := [3]int{112, 221, 122} //#70DD7A
	var r, g, b int

	if p < 0.5 {
		// Interpolate between Red and Orange
		t := p * 2
		r = int(float64(red[0]) + t*(float64(orange[0])-float64(red[0])))
		g = int(float64(red[1]) + t*(float64(orange[1])-float64(red[1])))
		b = int(float64(red[2]) + t*(float64(orange[2])-float64(red[2])))
	} else if p < 0.75 {
		// Interpolate between Orange and Yellow
		t := (p - 0.5) * 4
		r = int(float64(orange[0]) + t*(float64(yellow[0])-float64(orange[0])))
		g = int(float64(orange[1]) + t*(float64(yellow[1])-float64(orange[1])))
		b = int(float64(orange[2]) + t*(float64(yellow[2])-float64(orange[2])))
	} else {
		// Interpolate between Yellow and Green
		t := (p - 0.75) * 4
		r = int(float64(yellow[0]) + t*(float64(green[0])-float64(yellow[0])))
		g = int(float64(yellow[1]) + t*(float64(green[1])-float64(yellow[1])))
		b = int(float64(yellow[2]) + t*(float64(green[2])-float64(yellow[2])))
	}

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
