# Active Context

## Current Focus: Authentication Flow Enhancement

The immediate focus has shifted from component analysis to implementing improvements to the authentication flow and login page based on the comprehensive analysis completed. The goal is to create a more robust, user-friendly authentication experience with proper error handling and centralized redirection logic.

## Key Findings

### Authentication Flow

The authentication process is handled via an OAuth2 flow managed by `AuthContext` and the `authApi` service.

-   **`LoginPage.tsx`**: Presents login options (GitHub, Okta) and triggers the `login` function from `AuthContext`.
-   **`AuthContext.tsx`**: Manages the user's authentication state (`user`, `isAuthenticated`, `isLoading`). It initiates the OAuth2 flow by redirecting to the backend and fetches the current user's data on load.
-   **`authApi.ts`**: A service layer that handles all HTTP requests to the backend's `/auth` endpoints, including login, logout, and user data retrieval.
-   **`auth_middleware.go`**: The authentication middleware is responsible for validating user sessions and API keys. It attaches an `AuthContext` to the request, which is then used by the HTTP handlers to access user information.
-   **`http_handler.go`**: The HTTP handlers for the `/auth` endpoints have been updated to consistently use the `middleware.GetAuthContext(r)` helper function to retrieve the `AuthContext`. This resolves a critical bug where handlers were using an incorrect context key, leading to `401 Unauthorized` errors.

### Dashboard Architecture

The dashboard is a dynamic, grid-based system that allows for flexible component layouts.

-   **`DashboardContainer.tsx`**: The core component for rendering the dashboard grid using `react-grid-layout`. It receives a layout configuration and renders the specified components. It supports an "editing" mode for drag-and-drop and resizing.
-   **`ComponentRegistry.tsx`**: A crucial component that acts as a factory for dashboard widgets. It dynamically renders components based on a `type` string and passes the necessary props. It also defines the metadata for each available widget, including configuration options.
-   **`DashboardEditor.tsx`**: Provides the UI for adding new widgets to the dashboard, including modals for selecting and configuring widgets.
-   **`DashboardContext.tsx`**: Defines the context for sharing dashboard-related state, such as the selected project and suite IDs.

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

## Immediate Considerations

-   The `DashboardContext` is defined but not yet implemented, meaning there is no global state management for the dashboard's project and suite selections.
-   The dashboard components rely on `projectId`, `suiteId`, and `buildId` being passed down as props or from context. The mechanism for selecting these IDs and providing them to the dashboard needs to be clearly understood.
-   Consider implementing logout redirection logic in the centralized routing system for consistency.
