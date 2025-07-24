# System Patterns

## System Architecture
The application follows a classic three-tier architecture with a frontend, a backend, and a database. It is designed as a monorepo, which simplifies development and deployment.

- **Frontend**: A single-page application (SPA) built with React. It communicates with the backend via a RESTful API.
- **Backend**: A Go-based REST API that implements a hexagonal architecture. This separates the core application logic from external concerns like the database and HTTP handlers.
- **Database**: A PostgreSQL database that serves as the single source of truth for all test data.
- **CLI**: A Go-based command-line tool that acts as a client to the backend API, allowing for programmatic submission of test results.

## Design Patterns
- **Hexagonal Architecture (Ports and Adapters)**: The Go backend is structured around this pattern. The core application logic is isolated in the `internal` directory, with `domain`, `application`, and `infrastructure` layers. This promotes separation of concerns and testability.
- **Repository Pattern**: The backend uses repositories to abstract the data access logic, decoupling the application from the specific database implementation.
- **Dependency Injection**: The backend uses a container to manage and inject dependencies, which promotes loose coupling and testability.
- **RESTful API**: The backend exposes a RESTful API for the frontend and CLI to consume.
- **Component-Based UI**: The frontend is built with React, using a component-based architecture to create a modular and reusable UI. The new dashboard design introduces a widget-based system with a `ComponentRegistry` that dynamically renders components based on a layout configuration.
- **Widget-Based Dashboard**: The dashboard is composed of reusable widgets such as `MetricCard`, `StatusBadge`, and `DataChart`. This approach allows for flexible and customizable dashboard layouts.
- **Global Navigation**: A `BreadcrumbNavbar` has been implemented to provide consistent navigation and context across the application. It includes a project dropdown for easy project switching.
- **Semantic Color Scheme**: The UI now uses a semantic color palette to convey meaning and status. Colors are used consistently for errors (red), warnings (amber), success (green), and informational data (blue).
- **State Management**: The frontend uses React Context for managing global state, such as authentication status and dashboard layouts.

- **Dashboard Backend Component**: The dashboard backend is a key component responsible for providing data to the frontend for visualization. It is designed to be stateless and follows the hexagonal architecture pattern. It exposes a set of API endpoints for fetching dashboard data, such as project status, metrics, and chart data. The component relies on the `build` and `build_test_case_execution` modules to access the necessary data from the database.
- **Frontend Dashboard Component**: The frontend dashboard provides a flexible and interactive interface for visualizing project data. It uses a widget-based system with a `ComponentRegistry` that dynamically renders components based on a layout configuration. The dashboard's state is managed through a combination of React component state and the `useDashboardLayouts` hook.
- **Dynamic Chart Data**: The system now supports dynamic chart data, allowing users to visualize different data sets in the dashboard. The backend has been updated to support different chart types, and the frontend has been updated to handle the dynamic data.
- **[BROKEN]** **API Contract Mismatch**: A recent change to the backend API has broken the contract with the frontend. The `GetChartData` endpoint was modified to accept `project_id` as a query parameter instead of a path parameter, causing all chart data requests to fail with a 400 Bad Request error. This has rendered the dashboard's chart functionality completely non-functional.

## Component Relationships and Data Flow
1.  A user or CI/CD pipeline uses the **CLI** to submit a test result file (JUnit or ReadyAPI XML) to the **backend API**.
2.  The **backend API** receives the request, parses the data, and stores it in the **PostgreSQL database**.
3.  A user accesses the **React frontend** in their browser.
4.  The **frontend** makes API calls to the **backend** to fetch test data.
5.  The **backend** retrieves the data from the **database** and returns it to the **frontend**.
6.  The **frontend** renders the data in a user-friendly dashboard with tables and charts.

## Key Architectural Decisions
- **Monorepo**: The decision to use a monorepo simplifies cross-cutting concerns and dependency management between the frontend, backend, and CLI.
- **Hexagonal Architecture**: This choice for the backend architecture ensures a clean separation of concerns and makes the application more maintainable and testable.
- **Containerization**: Using Docker and Docker Compose for the entire stack simplifies the development setup and ensures consistency between environments.
