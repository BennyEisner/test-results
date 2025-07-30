# System Patterns

## Backend for Frontend (BFF)

The system employs a Backend for Frontend (BFF) pattern, where the backend API is specifically designed to serve the needs of the frontend. This pattern simplifies the frontend code by offloading complex data processing and aggregation to the backend.

### Key Characteristics

*   **Data Aggregation:** The backend aggregates data from multiple sources and provides it to the frontend in a single, easy-to-use format.
*   **Simplified Frontend:** The frontend is responsible for rendering the UI and handling user interactions, while the backend handles the business logic.
*   **Improved Performance:** The BFF pattern can improve performance by reducing the number of API calls the frontend needs to make.

## Hexagonal Architecture

The backend is built using a hexagonal architecture, which separates the core business logic from the infrastructure concerns. This makes the code more modular, testable, and maintainable.

### Key Characteristics

*   **Ports and Adapters:** The core business logic is exposed through a set of ports, and the infrastructure concerns are implemented as adapters that plug into these ports.
*   **Dependency Inversion:** The dependencies flow from the outside in, which means that the core business logic does not depend on the infrastructure concerns.
*   **Testability:** The hexagonal architecture makes it easy to test the core business logic in isolation from the infrastructure concerns.

## Centralized State Management

The frontend uses a centralized state management pattern, where the state of the application is stored in a single, global store. This makes it easy to manage the state of the application and to share data between components.

### Key Characteristics

*   **Single Source of Truth:** The global store is the single source of truth for the application's state.
*   **Predictable State Changes:** The state of the application can only be changed by dispatching actions, which makes the state changes predictable and easy to debug.
*   **Easy Data Sharing:** The global store makes it easy to share data between components, without having to pass props down through the component tree.
