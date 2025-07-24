package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/dashboard/domain/models"
)

// DashboardService defines the interface for dashboard business logic.
type DashboardService interface {
	GetStatus(ctx context.Context, projectID int64) (*models.StatusBadgeDTO, error)
	GetMetric(ctx context.Context, projectID int64, metricType string) (*models.MetricCardDTO, error)
	GetChartData(ctx context.Context, projectID int64, chartType string, suiteID *int64, buildID *int64) (*models.DataChartDTO, error)
	GetAvailableWidgets(ctx context.Context) (*models.AvailableWidgetsDTO, error)
}
