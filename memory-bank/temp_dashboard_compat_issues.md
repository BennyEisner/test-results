# Dashboard Compatibility Issues

## Status: PARTIALLY RESOLVED - NEW ISSUE IDENTIFIED

This document tracks compatibility issues discovered during the dashboard refactor and their resolution status.

## Issues Resolved

### 1. Data Model Mismatch (CRITICAL)
**Status:** RESOLVED
**Description:** The frontend and backend had a mismatch in the user configuration data model. The backend was returning an array of configs with `layouts` as a JSON string, while the frontend expected a single object with a parsed `layouts` array.
**Resolution:**
- The backend service `GetUserConfigs` was updated to return a single `*models.UserConfig` object.
- The corresponding HTTP handler was updated to match.
- The frontend hook `useDashboardLayouts.ts` was updated to correctly parse the `layouts` JSON string from the API response and handle the `active_layout_id` field.

### 2. API Response Structure Mismatch
**Status:** RESOLVED
**Description:** The API handler returned `active_layout_id`, but the frontend was expecting `activeId`.
**Resolution:** The frontend was updated to correctly use `active_layout_id`.

### 3. Service Interface Inconsistency
**Status:** RESOLVED
**Description:** The `UserConfigService` interface was defined to return a slice of configs (`[]*models.UserConfig`), which was inconsistent with the application's logic of one configuration per user.
**Resolution:** The interface was changed to return a single `*models.UserConfig`.

## New Issue Identified

### 1. Chart Data Fetching Fails with 400 Error
**Status:** ACTIVE DEBUGGING
**Location:** `/api/dashboard/projects/{id}/chart/bar`
**Description:** After resolving the initial data model issues, a new problem has surfaced. When a project is selected in the UI, the API call to fetch chart data fails with a `400 Bad Request`.
**Impact:** Dashboard charts do not load, showing an error state.

## Next Steps

The immediate priority is to investigate the `400 Bad Request` error. This will involve:
1.  **Backend Investigation:** Examine the `GetChartData` handler in `api/internal/dashboard/infrastructure/http/http_handler.go` to determine what conditions would lead to a `400` error. This could be due to missing or invalid query parameters.
2.  **Frontend Investigation:** Analyze the data fetching logic in `frontend/src/components/widgets/DataChart.tsx` and the relevant service calls in `frontend/src/services/dashboardApi.ts` to ensure the correct parameters (e.g., `projectId`, `chartType`, any required query params) are being sent.
3.  **Address Redundant API Calls:** Investigate the multiple `saveLayouts` calls on page load to prevent unnecessary backend traffic and potential race conditions.
