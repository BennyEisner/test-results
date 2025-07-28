# Active Context

## Current Work Focus
The primary focus is on ensuring the dashboard is fully functional and documenting the recent fixes.

## Recent Changes
1.  **`useSmartRefresh` Hook:** Fixed a logic flaw that caused excessive API calls. The hook now correctly respects the `refreshOn` parameter, preventing unnecessary data fetching.
2.  **`useSmartRefresh` Hook:** Fixed a logic flaw that prevented static widgets from loading on the initial render. The hook now correctly triggers a data fetch when the component mounts.
3.  **`ComponentConfigModal`:** Corrected an issue where the `limit` parameter was being handled as a string instead of a number, which caused the limit functionality to fail.

## Current State
-   **Dashboard Performance:** The excessive API calls on the dashboard have been resolved.
-   **Static Charts:** The loading issue has been resolved, and all static widgets now render correctly.
-   **Dynamic Chart Limit:** The `limit` option is now fully functional.

## Root Cause Analysis
-   **Excessive API Calls:** The `useSmartRefresh` hook's `useEffect` dependency array was not correctly using the `refreshOn` parameter, causing it to re-fetch data on every context change.
-   **Static Charts Loading Issue:** The `useSmartRefresh` hook was not triggering a data fetch on the initial render for static components because the trigger logic did not account for the initial state.
-   **Limit Functionality Issue:** The `ComponentConfigModal` was not converting the `limit` value from a string to a number, causing the API to ignore the parameter.

## Next Steps
- ✅ Update the `progress.md` and `systemPatterns.md` files in the memory bank to reflect the recent fixes.
- The dashboard performance issue has been resolved. The excessive API calls were caused by the `useSmartRefresh` hook not properly respecting the `refreshOn` parameter.
