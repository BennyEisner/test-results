# Data Chart Build-Level Investigation Plan

## 1. Objective

The primary objective of this investigation is to identify and resolve the root cause of the "Test Case Pass Rate" chart failing to display data when filtered by a specific build. The chart works correctly at the project and suite levels, but not at the build level.

## 2. Summary of Changes Made

To address a related race condition and "stuck in loading" state, the following changes were implemented:

1.  **`useSmartRefresh` Hook:**
    *   Integrated an `AbortController` to cancel in-flight `fetch` requests when a new request is initiated. This prevents race conditions where old data could overwrite new data.
2.  **`dashboardApi` Service:**
    *   The `getChartData` function was updated to accept an `AbortSignal` and pass it to the underlying `api.get` call.
3.  **`DataChart` Component:**
    *   The `fetcher` function passed to `useSmartRefresh` was memoized using `useCallback`. This prevents the hook from re-triggering unnecessarily on every render, which was causing requests to be aborted prematurely and leaving the UI in a "loading" state.

While these changes fixed the UI stability issues, the original data loading problem at the build level persists.

## 3. Investigation Strategy

The investigation will proceed in a layered approach, starting from the frontend and moving systematically to the backend.

### Phase 1: Frontend Verification

This phase focuses on ensuring the frontend is correctly requesting and handling the data.

1.  **Verify API Request Parameters:**
    *   **Action:** Add logging in the `DataChart` component to inspect the exact parameters being passed to the `dashboardApi.getChartData` function, especially the `buildId`.
    *   **Expected Outcome:** Confirm that a valid `buildId` is being passed when a build is selected.

2.  **Inspect Network Response:**
    *   **Action:** Use the browser's developer tools to inspect the network request for the chart data when a build is selected. Check the response payload.
    *   **Expected Outcome:**
        *   If the response payload is empty or contains no data, the issue is likely in the backend.
        *   If the response payload contains valid data, the issue is in the frontend's data handling or rendering logic.

3.  **Analyze Frontend Data Transformation:**
    *   **Action:** If the network response contains data, add logging in the `transformData` function within `DataChart.tsx` to inspect the data before and after transformation.
    *   **Expected Outcome:** Confirm that the data is being transformed correctly into the format expected by the charting library.

### Phase 2: API & Backend Investigation

If Phase 1 indicates a backend issue, this phase will focus on the API endpoint and the database query.

1.  **Review API Handler (`api/internal/dashboard/infrastructure/http/http_handler.go`):**
    *   **Action:** Inspect the `GetChartData` handler. Add logging to check the `build_id` received from the request and the parameters passed to the `dashboardService`.
    *   **Expected Outcome:** Confirm that the `build_id` is being correctly parsed and passed to the service layer.

2.  **Examine Application Service (`api/internal/dashboard/application/service.go`):**
    *   **Action:** Review the `GetChartData` method in the `DashboardService`. Add logging to inspect the parameters passed to the repository.
    *   **Expected Outcome:** Confirm that the service layer is correctly calling the repository with the `build_id`.

3.  **Analyze Database Repository (`api/internal/dashboard/infrastructure/database/repository.go`):**
    *   **Action:** This is the most likely location of the issue. Carefully review the SQL query for the "Test Case Pass Rate" chart (`pass-fail-trend`). Pay close attention to how the `build_id` is used in the `WHERE` clause and any `JOIN`s.
    *   **Expected Outcome:** Identify any logical errors in the SQL query that would cause it to return no data when a `build_id` is provided. It's possible the query is structured in a way that only works for project or suite-level aggregation.

4.  **Direct Database Query:**
    *   **Action:** Manually execute the SQL query from the repository directly against the database with a known `build_id` that should have data.
    *   **Expected Outcome:**
        *   If the query returns no data, the SQL logic is flawed.
        *   If the query returns data, the issue may be in how the Go application is handling the database connection or processing the results.

## 4. Next Steps

Based on the outcome of this investigation, the next step will be to implement a fix in the appropriate location (frontend, API handler, service, or repository). The fix will then be tested thoroughly to ensure the chart works correctly at all levels (project, suite, and build). After the fix is confirmed, all temporary logging will be removed.
