# Progress

## What Works

- **GitHub OAuth2 Authentication**: The end-to-end authentication flow with GitHub is now fully functional. A series of complex issues, including session management conflicts, provider name extraction errors, incorrect credential usage, CORS misconfiguration, and a context key mismatch, have been diagnosed and resolved.
- **API Key Management**: The API key management functionality is now fully operational. A bug that caused a `401 Unauthorized` error when creating, listing, or deleting API keys has been resolved. The issue was traced to an incorrect context key being used in the backend handlers, which has been corrected. Additionally, a related frontend bug that occurred when a user had no API keys has been fixed.
- **Protected Routes**: The `ProtectedRoute` component effectively guards protected routes, ensuring that only authenticated users can access them.
- **Routing**: The routing logic has been improved to handle authentication state changes correctly. Unauthenticated users are now redirected to a dedicated `/login` route, while authenticated users are redirected to the dashboard.
- **API Base Path**: The `authApi` service has been updated to use an absolute path for authentication requests, which resolves the 400 error that occurred when initiating login from the home page.
- **Proxy Configuration**: The Vite proxy has been updated to correctly rewrite the path for API requests, ensuring that they are properly routed to the backend.
- **Dynamic Dashboard**: A flexible, grid-based dashboard system is in place. It allows for dynamic rendering of various components (widgets) like tables and charts.
- **Dashboard Editing**: The dashboard supports an editing mode that allows users to add, remove, and rearrange components via drag-and-drop.
- **Component Registry**: A component registry dynamically renders widgets and defines their properties, including configuration options for static components.
- **Enhanced Authentication Flow**: The login process is now more robust with centralized redirection, loading and error states, and improved user feedback on the login page.
- **Analytical Dashboard Redesign**: The application's UI has been transformed into a statistics-focused dashboard. This includes:
    - A professional, data-centric design with a clear visual hierarchy.
    - Reusable widget components (`MetricCard`, `StatusBadge`, `DataChart`).
    - A semantic color scheme for intuitive status indication.
    - A flexible grid-based layout using `react-grid-layout`.
- **Dashboard Backend**: The backend for the dashboard has been implemented, following the hexagonal architecture pattern. This includes a new `dashboard` domain with its own services, repositories, and HTTP handlers. The backend provides endpoints for fetching status, metrics, and chart data.
- **Frontend Dashboard**: The frontend dashboard provides a flexible and interactive interface for visualizing project data. It uses a widget-based system with a `ComponentRegistry` that dynamically renders components based on a layout configuration.

## What's Left to Build

- **User Roles and Permissions**: The system currently lacks a role-based access control (RBAC) system, which is necessary for managing user permissions and restricting access to certain features.
- **Two-Factor Authentication (2FA)**: To enhance security, 2FA should be implemented to provide an additional layer of protection for user accounts.
- **Audit Logs**: There is no auditing mechanism to track user activities, such as login attempts, API key creation, or other sensitive operations.

## Known Issues

- **[CRITICAL REGRESSION]** **Dashboard Charts Broken:** All dashboard charts are currently broken due to a series of backend changes that misaligned the API contract between the frontend and backend. All chart-related API calls are failing with a 400 Bad Request error.
- **Hardcoded Widgets**: The `GetAvailableWidgets` function in the dashboard service returns a hardcoded list of widgets. This should be replaced with a dynamic, configuration-based approach.
- **Frontend Error Handling**: The dashboard displays placeholder messages when required data is not available, but more robust error handling is needed to handle API errors and other unexpected issues.
- **No Rate Limiting**: The authentication endpoints lack rate limiting, which makes them vulnerable to brute-force attacks.
- **Improved Error Handling**: The login flow now has improved error handling, displaying messages to the user on failure. However, more comprehensive error handling across the application is still needed.
- **[RESOLVED]** **Data Not Rendering from API:** Fixed the issue where the `DataChart` component was not displaying data from the API. The problem was that the backend `GetChartData` function was not handling the `bar` and `line` chart types that were being requested by the frontend. Modified the service to return hardcoded data for these chart types.
- **[RESOLVED]** **Grid Layout Race Condition:** Fixed a race condition in the dashboard grid layout system where layout updates were attempted before initialization was complete. Added proper initialization checks to prevent the "Attempted to update grid layout before initialization is complete" error.
- **[RESOLVED]** **`build-duration` Chart Error:** Fixed an issue where the `build-duration` chart was failing due to a data type mismatch. The repository now correctly scans the `duration` as a `float64` and casts it to an `int`.
- **Redundant API Calls:** The dashboard makes multiple, seemingly unnecessary, calls to save the layout on initial load.
