# Dynamic Chart Data Grouping

This document outlines the dynamic data grouping logic for dashboard charts, specifically for "Build Duration" and "Test Case Pass Rate." The grouping behavior changes based on the context provided (project, suite, or build), ensuring that the visualizations are always meaningful and relevant.

## Context-Aware Grouping

The core principle is to adapt the data aggregation to the user's current view. Here’s how it works for each chart type:

### 1. Build Duration Chart

-   **Chart Type:** `build-duration`
-   **Data Type:** Pie/Doughnut

| Context | Grouping | Description |
| :--- | :--- | :--- |
| **Project Only** | By Test Suite | Shows the **average** build duration for each test suite within the project. |
| **Project and Suite** | By Build | Shows the duration for each individual build within the selected suite. |

### 2. Test Case Pass Rate Chart

-   **Chart Type:** `test-case-pass-rate`
-   **Data Type:** Pie/Doughnut

| Context | Grouping | Description |
| :--- | :--- | :--- |
| **Project Only** | By Test Suite | Shows the **overall** pass rate for each test suite within the project. |
| **Project and Suite** | By Build | Shows the pass rate for each individual build within the selected suite. |
| **Build** | By Test Case | Shows the pass/fail status for each test case within the selected build. |

## Implementation Details

The logic is implemented in the `GetChartData` function in `api/internal/build_test_case_execution/infrastructure/database/repository.go`. This function uses a helper, `getChartQuery`, to select the appropriate SQL query based on the provided `suiteID` and `buildID`.

By centralizing this logic in the backend, the frontend `DataChart` component remains simple and does not need to be aware of the grouping rules. It simply requests the data for a given context, and the backend returns the correctly aggregated data.
