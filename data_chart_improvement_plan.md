# Data Chart Improvement Plan

## 1. Executive Summary

The current implementation of the "Test Case Pass Rate" chart in the dashboard fails to display data when filtered by a specific build. This plan outlines the necessary changes to the backend to resolve this issue. The core of the problem lies in the backend's SQL query for this chart, which returns an empty result set for single-build views.

The proposed solution involves creating a dedicated query for the build-level view of the "Test Case Pass Rate" chart.

## 2. Problem Analysis

*   **Root Cause:** The SQL query for the `test-case-pass-rate` chart is designed for aggregation over multiple builds (at the project or suite level). When filtered to a single build, it fails to return any data, causing the frontend to render an empty chart.
*   **User Experience Issue:** The chart is not providing the expected information at the build-detail level, which hinders the user's ability to analyze test results for a specific build.

## 3. Proposed Solution

I will implement a targeted backend fix to resolve this issue.

*   **Backend Query Refactoring:**
    *   **File to Modify:** `api/internal/build_test_case_execution/infrastructure/database/repository.go`
    *   **Function to Modify:** `GetChartData`
    *   **Change:** I will introduce a special case for `chartType == "test-case-pass-rate"` when `buildID` is not `nil`.
    *   **New Query:** A simplified and more direct query will be used for this case:
        ```sql
        SELECT
            tc.name as label,
            CASE WHEN e.status = 'passed' THEN 100.0 ELSE 0.0 END as value
        FROM build_test_case_executions e
        JOIN test_cases tc ON e.test_case_id = tc.id
        WHERE e.build_id = $1
        ORDER BY tc.name
        ```
    *   **Benefit:** This will immediately fix the "empty chart" issue by ensuring that data is returned for the build-level view. The frontend will then be able to render the chart with the correct data.

## 4. Implementation Steps

1.  Modify `api/internal/build_test_case_execution/infrastructure/database/repository.go` to implement the backend query refactoring as described above.

## 5. Expected Outcome

After implementing the backend fix, the "Test Case Pass Rate" chart will correctly display data when a build is selected. It will show a bar for each test case in the build, with a value of 100 if it passed and 0 if it failed. This resolves the immediate bug and provides the user with the correct information for their analysis.
