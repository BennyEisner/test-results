# Active Context

## Current Focus: Analytical Dashboard Implementation

The current focus is on the implementation of a new analytical, statistics-focused dashboard. This involves creating a cohesive visual system, developing reusable widget components, and ensuring a professional and data-centric user experience.

## Key Findings

### Authentication Flow

The authentication process is handled via an OAuth2 flow managed by `AuthContext` and the `authApi` service.

-   **`LoginPage.tsx`**: Presents login options (GitHub, Okta) and triggers the `login` function from `AuthContext`.
-   **`AuthContext.tsx`**: Manages the user's authentication state (`user`, `isAuthenticated`, `isLoading`). It initiates the OAuth2 flow by redirecting to the backend and fetches the current user's data on load.
-   **`authApi.ts`**: A service layer that handles all HTTP requests to the backend's `/auth` endpoints, including login, logout, and user data retrieval.
-   **`auth_middleware.go`**: The authentication middleware is responsible for validating user sessions and API keys. It attaches an `AuthContext` to the request, which is then used by the HTTP handlers to access user information.
-   **`http_handler.go`**: The HTTP handlers for the `/auth` endpoints have been updated to consistently use the `middleware.GetAuthContext(r)` helper function to retrieve the `AuthContext`. This resolves a critical bug where handlers were using an incorrect context key, leading to `401 Unauthorized` errors.

### Dashboard Architecture

The dashboard has been redesigned with a focus on analytics and data visualization.

-   **`DashboardContainer.tsx`**: The core component for rendering the dashboard grid using `react-grid-layout`. It manages the layout of widgets and passes down necessary context, such as `projectId` and `suiteId`.
-   **`ComponentRegistry.tsx`**: A factory for dashboard widgets that dynamically renders components based on a `type` string. It has been updated to support new widget types: `MetricCard`, `StatusBadge`, and `DataChart`.
-   **Widget Components**: A new set of reusable widget components has been created in `frontend/src/components/widgets/`:
    -   `MetricCard.tsx`: Displays a single metric with a title, value, and trend indicator.
    -   `StatusBadge.tsx`: A badge for displaying status information with semantic coloring.
    -   `DataChart.tsx`: A versatile chart component for visualizing data.
-   **Styling**: A dedicated CSS file, `frontend/src/styles/dashboard.css`, has been created to provide a consistent and professional look and feel for the dashboard, including a semantic color scheme.
-   **Types**: The `frontend/src/types/dashboard.ts` file has been updated to include the new widget types and configuration options.

### Recent Authentication Enhancements

**Post-Login Redirection Improvements:**
-   **Centralized Redirection Logic**: Moved redirection logic from `AuthContext.tsx` to `App.tsx` using a new `AppRoutes` component that observes authentication state changes.
-   **Automatic Dashboard Redirect**: Authenticated users visiting `/login` are now automatically redirected to `/dashboard`.
-   **Clean Separation of Concerns**: `AuthContext` now focuses solely on state management, while routing logic is handled in the appropriate component.

**Login Page Enhancements:**
-   **Error State Management**: Added comprehensive error handling to `AuthContext` with `error` state and `clearError` function.
-   **Loading States**: Login buttons now show loading spinners and are disabled during authentication process.
-   **User Feedback**: Error messages are displayed prominently on the login page when authentication fails.
-   **Already Authenticated Handling**: Users who are already logged in see a clear "Go to Dashboard" button instead of login options.

## Next Steps

-   Continue to refine the dashboard by adding more widget types and configuration options.
-   Implement the `DashboardContext` to provide global state management for dashboard-related selections.
-   Enhance the data visualization capabilities of the `DataChart` component.
