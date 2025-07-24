# Dynamic Chart Data

This document outlines the current implementation of the chart data fetching mechanism and proposes a plan for making it dynamic, allowing users to select the data they want to see in the charts.

## Current Implementation

### Frontend

The `DataChart` component is responsible for fetching and rendering chart data. It takes a `chartType` prop, which it uses to request data from the backend.

```typescript
// frontend/src/components/widgets/DataChart.tsx

const DataChart = ({ projectId, chartType }: DataChartProps) => {
    // ...
    useEffect(() => {
        if (projectId) {
            const fetchChartData = async () => {
                try {
                    setLoading(true);
                    const response = await dashboardApi.getChartData(Number(projectId), chartType);
                    setChartData(response.chart_data);
                } catch (err) {
                    setError(err as Error);
                } finally {
                    setLoading(false);
                }
            };
            fetchChartData();
        }
    }, [projectId, chartType]);
    // ...
};
```

The `dashboardApi.getChartData` function makes a GET request to the `/dashboard/projects/{projectId}/chart/{chartType}` endpoint.

### Backend

The `GetChartData` handler in `api/internal/dashboard/infrastructure/http/http_handler.go` receives the request and calls the `GetChartData` service.

The `GetChartData` service in `api/internal/dashboard/application/service.go` currently returns hardcoded data for the `bar` and `line` chart types.

```go
// api/internal/dashboard/application/service.go

func (s *DashboardServiceImpl) GetChartData(ctx context.Context, projectID int64, chartType string) (*models.DataChartDTO, error) {
	// Hardcoded data for testing
	if chartType == "build-duration" || chartType == "bar" {
		// ... returns hardcoded data
	}

	if chartType == "pass-fail-trend" || chartType == "line" {
		// ... returns hardcoded data
	}

	return s.buildExecRepo.GetChartData(ctx, projectID, chartType)
}
```

## Proposed Changes for Dynamic Data

To allow users to select the data they want to see, we need to make the following changes:

### Frontend

1.  **Update `ComponentConfigModal`:**
    *   Add a new field to the modal that allows users to select the type of data they want to see in the chart (e.g., "Build Duration", "Pass/Fail Trend", "Execution Time").
    *   This field should be a dropdown that is populated with the available chart types from the backend.
    *   When the user selects a chart type, the `chartType` prop of the `DataChart` component should be updated.

2.  **Update `DataChart` component:**
    *   The `DataChart` component will receive the selected `chartType` as a prop and use it to fetch the corresponding data from the backend.

### Backend

1.  **Update `GetChartData` service:**
    *   Remove the hardcoded data.
    *   The `chartType` parameter will now be used to determine which data to fetch from the repository.
    *   The service will call the appropriate repository method based on the `chartType`.

2.  **Update `build_test_case_execution` repository:**
    *   The `GetChartData` function will be modified to accept a `chartType` parameter.
    *   It will use this parameter to construct the correct SQL query to fetch the requested data.

By implementing these changes, we will create a flexible and user-friendly charting system that allows users to visualize the data that is most important to them.
