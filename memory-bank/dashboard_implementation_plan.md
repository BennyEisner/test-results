# Validated Plan for Enhanced Dashboard Implementation

This plan is designed to be fully compatible with the existing architecture and coding standards.

---

### **Phase 1: Backend API Development & Testing**

This phase focuses on creating a new `dashboard` domain to encapsulate all data-fetching logic for the new widgets, ensuring a clean separation of concerns as per the hexagonal architecture.

1.  **Create New `dashboard` Domain:**
    *   **File Creation:**
        *   `api/internal/dashboard/domain/models/models.go`: Define DTOs (`StatusBadgeDTO`, `MetricCardDTO`, `DataChartDTO`) for clear data contracts.
        *   `api/internal/dashboard/domain/ports/ports.go`: Define the `DashboardService` interface.
        *   `api/internal/dashboard/application/service.go`: Implement the `DashboardService`, which will orchestrate calls to other repositories.
        *   `api/internal/dashboard/infrastructure/http/http_handler.go`: Create the HTTP handler for the new dashboard endpoints.
    *   **Dependency Injection (`api/internal/shared/container/container.go`):** Register the new `DashboardService` and `DashboardHandler`.
    *   **Routing (`api/cmd/server/main.go`):** Register the new API routes:
        *   `GET /api/projects/{projectID}/dashboard/status`
        *   `GET /api/projects/{projectID}/dashboard/metric/{metricType}`
        *   `GET /api/projects/{projectID}/dashboard/chart/{chartType}`
        *   `GET /api/dashboard/available-widgets`

2.  **Extend Existing Repositories:**
    *   **`api/internal/build/domain/ports/ports.go`**: Add `GetLatestBuildStatus` to the `BuildRepository` interface.
    *   **`api/internal/build/infrastructure/database/repository.go`**: Implement `GetLatestBuildStatus` with the required SQL query.
    *   **`api/internal/build_test_case_execution/domain/ports/ports.go`**: Add `GetMetric` and `GetChartData` to the `BuildTestCaseExecutionRepository` interface.
    *   **`api/internal/build_test_case_execution/infrastructure/database/repository.go`**: Implement the new methods with aggregation queries (e.g., `COUNT`, `AVG`, `GROUP BY`).

3.  **Backend Unit Testing:**
    *   Write unit tests for the new `DashboardService` to ensure the business logic is correct.
    *   Write unit tests for the `DashboardHandler` to verify request handling and response formatting.

---

### **Phase 2: Frontend Integration & Testing**

This phase connects the UI components to the new backend endpoints via a dedicated context, ensuring a clean data flow.

1.  **Establish Data Flow:**
    *   **`frontend/src/services/dashboardApi.ts`**: Create this new file to manage all API calls to the dashboard endpoints.
    *   **`frontend/src/types/dashboard.ts`**: Add TypeScript types corresponding to the backend DTOs.
    *   **`frontend/src/context/DashboardContext.tsx`**: Fully implement this context to:
        *   Fetch data using `dashboardApi.ts`.
        *   Manage loading and error states gracefully.
        *   Provide data and state to consumer components.

2.  **Connect Widgets to Live Data:**
    *   Refactor `StatusBadge.tsx`, `MetricCard.tsx`, and `DataChart.tsx` to remove mock data and consume live data from the `DashboardContext`.

3.  **Frontend Unit Testing:**
    *   Write unit tests for the `DashboardContext` to verify its state management logic.
    *   Update tests for the widget components to ensure they render correctly with data from the context.

---

### **Phase 3: Dynamic Configuration & Finalization**

This phase empowers users to customize their dashboards, with changes persisted on the backend.

1.  **Implement Widget Configuration:**
    *   **`frontend/src/components/dashboard/ComponentConfigModal.tsx`**: Enhance the modal to allow users to select the `metricType` or `chartType` for each widget.
    *   **`frontend/src/components/dashboard/DashboardContainer.tsx`**: Update this component to pass the user's widget configuration to the `ComponentRegistry`.
    *   **`frontend/src/components/dashboard/ComponentRegistry.tsx`**: Modify the registry to pass the configuration down to the individual widget instances.

2.  **Persist Configuration:**
    *   **`api/internal/user_config/application/service.go`**: Update the `SaveUserConfig` service to handle the additional widget configuration data within the `layouts` JSON blob.

3.  **Documentation and Code Comments:**
    *   **Memory Bank Update:**
        *   `memory-bank/progress.md`: Update the "What Works" section and feature checklist.
        *   `memory-bank/activeContext.md` & `memory-bank/systemPatterns.md`: Document the new `dashboard` domain and its data flows.
    *   **Code Comments:** Add comments to explain complex SQL queries and business logic in the new services.
