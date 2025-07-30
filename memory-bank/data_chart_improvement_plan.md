# Data Chart Improvement Plan

## 1. Overview

This plan outlines the steps to enhance the data charts in the dashboard, focusing on improving readability and aesthetics. The key areas of improvement are:

*   **Descriptive Labels:** Adding labels to the axes and data points to provide more context.
*   **Color Scheme:** Implementing a more effective color scheme to improve data differentiation and visual appeal.

## 2. Architecture and Implementation Strategy

The implementation will follow the existing architectural patterns of the application, including the Backend for Frontend (BFF) and centralized state management.

### 2.1. Backend (Go)

The backend will be responsible for providing the frontend with the necessary data to render the charts, including the labels and colors.

*   **Data Models:** The `Chart` and `ChartData` models in `api/internal/dashboard/domain/models/models.go` will be updated to include fields for labels and colors.
*   **API Endpoints:** The dashboard API endpoints in `api/internal/dashboard/infrastructure/http/http_handler.go` will be updated to return the new data.
*   **Database Queries:** The database queries in `api/internal/dashboard/infrastructure/database/repository.go` will be updated to retrieve the new data.

### 2.2. Frontend (React)

The frontend will be responsible for rendering the charts with the new labels and colors.

*   **`DataChart` Component:** The `DataChart` component in `frontend/src/components/widgets/DataChart.tsx` will be refactored to support the new features. This will involve:
    *   Updating the component to accept the new data from the backend.
    *   Using a charting library (e.g., Chart.js, Recharts) to render the charts.
    *   Configuring the charting library to display the labels and colors correctly.
*   **`dashboardApi.ts`:** The `dashboardApi.ts` service in `frontend/src/services/dashboardApi.ts` will be updated to fetch the new data from the backend.
*   **`dashboard.ts`:** The `dashboard.ts` types in `frontend/src/types/dashboard.ts` will be updated to include the new data fields.

## 3. Detailed Implementation Steps

### 3.1. Backend

1.  **Update Data Models:**
    *   Add `XAxisLabel`, `YAxisLabel`, and `DataPointLabel` fields to the `Chart` model.
    *   Add a `Color` field to the `ChartDataPoint` model.
2.  **Update Database Queries:**
    *   Modify the queries to retrieve the new label and color data from the database.
3.  **Update API Endpoints:**
    *   Update the API endpoints to return the new data in the JSON response.

### 3.2. Frontend

1.  **Update Types:**
    *   Update the `Chart` and `ChartDataPoint` types in `frontend/src/types/dashboard.ts` to include the new fields.
2.  **Update API Service:**
    *   Update the `getChartData` function in `frontend/src/services/dashboardApi.ts` to fetch the new data.
3.  **Refactor `DataChart` Component:**
    *   Install a charting library (e.g., `npm install chart.js react-chartjs-2`).
    *   Replace the existing chart rendering logic with the new charting library.
    *   Configure the chart to use the new labels and colors from the backend.

## 4. Color Scheme

The new color scheme will be designed to be both aesthetically pleasing and effective at differentiating data. The following is a proposed color palette:

*   **Primary Color:** `#3B82F6` (blue)
*   **Secondary Color:** `#10B981` (green)
*   **Accent Color:** `#F59E0B` (amber)
*   **Neutral Colors:** `#6B7280` (gray), `#9CA3AF` (light gray)

These colors will be used consistently across all charts to ensure a cohesive look and feel.

## 5. Timeline

The implementation will be carried out in the following order:

1.  **Backend:** 1-2 days
2.  **Frontend:** 2-3 days
3.  **Testing and Deployment:** 1 day

This plan provides a clear roadmap for improving the data charts in the dashboard. By following these steps, we can create a more informative and visually appealing experience for our users.
