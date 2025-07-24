package http

import "github.com/BennyEisner/test-results/internal/dashboard/domain/models"

type GetStatusResponse struct {
	Status *models.StatusBadgeDTO `json:"status"`
}

type GetMetricResponse struct {
	Metric *models.MetricCardDTO `json:"metric"`
}

type GetChartDataResponse struct {
	ChartData *models.DataChartDTO `json:"chart_data"`
}

type GetAvailableWidgetsResponse struct {
	Widgets *models.AvailableWidgetsDTO `json:"widgets"`
}
