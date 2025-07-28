package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/dashboard/domain/ports"
)

type DashboardHandler struct {
	service ports.DashboardService
}

func NewDashboardHandler(service ports.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectID")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetStatus(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetStatusResponse{Status: status}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DashboardHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectID")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	metricType := r.PathValue("metricType")
	metric, err := h.service.GetMetric(r.Context(), projectID, metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetMetricResponse{Metric: metric}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DashboardHandler) GetChartData(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectID")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	chartType := r.PathValue("chartType")

	suiteIDStr := r.URL.Query().Get("suite_id")
	var suiteID *int64
	if suiteIDStr != "" {
		id, err := strconv.ParseInt(suiteIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid suite ID", http.StatusBadRequest)
			return
		}
		suiteID = &id
	}

	buildIDStr := r.URL.Query().Get("build_id")
	var buildID *int64
	if buildIDStr != "" {
		id, err := strconv.ParseInt(buildIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid build ID", http.StatusBadRequest)
			return
		}
		buildID = &id
	}

	limitStr := r.URL.Query().Get("limit")
	var limit *int
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = &l
	}

	chartData, err := h.service.GetChartData(r.Context(), projectID, chartType, suiteID, buildID, limit)
	if err != nil {
		if err.Error() == fmt.Sprintf("unknown chart type: %s", chartType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := GetChartDataResponse{ChartData: chartData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DashboardHandler) GetAvailableWidgets(w http.ResponseWriter, r *http.Request) {
	widgets, err := h.service.GetAvailableWidgets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(widgets)
}
