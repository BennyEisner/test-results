package application

import (
	"context"

	buildPorts "github.com/BennyEisner/test-results/internal/build/domain/ports"
	buildExecPorts "github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/ports"
	"github.com/BennyEisner/test-results/internal/dashboard/domain/models"
)

type DashboardServiceImpl struct {
	buildRepo     buildPorts.BuildRepository
	buildExecRepo buildExecPorts.BuildTestCaseExecutionRepository
}

func NewDashboardService(buildRepo buildPorts.BuildRepository, buildExecRepo buildExecPorts.BuildTestCaseExecutionRepository) *DashboardServiceImpl {
	return &DashboardServiceImpl{
		buildRepo:     buildRepo,
		buildExecRepo: buildExecRepo,
	}
}

func (s *DashboardServiceImpl) GetStatus(ctx context.Context, projectID int64) (*models.StatusBadgeDTO, error) {
	status, err := s.buildRepo.GetLatestBuildStatus(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &models.StatusBadgeDTO{Status: status}, nil
}

func (s *DashboardServiceImpl) GetMetric(ctx context.Context, projectID int64, metricType string) (*models.MetricCardDTO, error) {
	return s.buildExecRepo.GetMetric(ctx, projectID, metricType)
}

func (s *DashboardServiceImpl) GetChartData(ctx context.Context, projectID int64, chartType string, suiteID *int64, buildID *int64) (*models.DataChartDTO, error) {
	return s.buildExecRepo.GetChartData(ctx, projectID, chartType, suiteID, buildID)
}

func (s *DashboardServiceImpl) GetAvailableWidgets(ctx context.Context) (*models.AvailableWidgetsDTO, error) {
	// In the future, this could be fetched from a dynamic source (e.g., config file, database)
	return &models.AvailableWidgetsDTO{
		Metrics: []models.WidgetOption{
			{Value: "passing-rate", Label: "Passing Rate"},
			{Value: "execution-time", Label: "Execution Time"},
		},
		Charts: []models.WidgetOption{
			{Value: "build-duration", Label: "Build Duration"},
			{Value: "pass-fail-trend", Label: "Pass/Fail Trend"},
			{Value: "test-case-pass-rate", Label: "Test Case Pass Rate"},
		},
	}, nil
}
