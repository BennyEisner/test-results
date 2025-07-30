package models

// StatusBadgeDTO represents the data for a status badge widget.
type StatusBadgeDTO struct {
	Status string `json:"status"`
	Text   string `json:"text"`
}

// MetricCardDTO represents the data for a metric card widget.
type MetricCardDTO struct {
	Title      string `json:"title"`
	Value      string `json:"value"`
	Change     string `json:"change,omitempty"`
	ChangeType string `json:"changeType,omitempty"`
}

// DataChartDTO represents the data for a chart widget.
type DataChartDTO struct {
	Labels     []string     `json:"labels"`
	Datasets   []DatasetDTO `json:"datasets"`
	XAxisLabel string       `json:"xAxisLabel"`
	YAxisLabel string       `json:"yAxisLabel"`
}

type WidgetOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type AvailableWidgetsDTO struct {
	Metrics []WidgetOption `json:"metrics"`
	Charts  []WidgetOption `json:"charts"`
}

// DatasetDTO represents a dataset for a chart.
type DatasetDTO struct {
	Label           string   `json:"label"`
	Data            []int    `json:"data"`
	BackgroundColor []string `json:"backgroundColor,omitempty"`
	BorderColor     []string `json:"borderColor,omitempty"`
}
