# Hexagonal Architecture and Development Workflow

This document outlines the Hexagonal (Ports and Adapters) Architecture used in this project. Adherence to these principles is mandatory for all development to ensure a clean, maintainable, and scalable codebase.

## 1. Core Principles of Our Hexagonal Architecture

Our architecture is divided into three distinct layers, with a strict dependency rule: **dependencies must always point inwards**.

```
+-----------------------------------------------------------------+
|                                                                 |
|  +-----------------+      +-----------------+      +----------+ |
|  |   HTTP Handler  |----->| Application Svc |----->|  Domain  | |
|  +-----------------+      +-----------------+      +----------+ |
|        (Input Adapter)          (Use Case)           (Core)     |
|                                     ^                           |
|                                     |                           |
|  +-----------------+                |                           |
|  | DB Repository   |----------------+                           |
|  +-----------------+                                            |
|        (Output Adapter)                                         |
|                                                                 |
+------------------ Infrastructure Layer -------------------------+
```

### Layers and Responsibilities

1.  **Domain (The Core)**
    *   **Contains:** Business logic, entities, value objects, and **ports** (Go interfaces).
    *   **Responsibility:** Represents the heart of the application. It is completely independent of any external technology or framework. It knows nothing about the database, the API, or any other service.
    *   **Example:** `build/domain/models/models.go`, `build/domain/ports/ports.go`

2.  **Application (The Use Cases)**
    *   **Contains:** Application services that orchestrate the business logic.
    *   **Responsibility:** Implements the application's use cases by coordinating the domain objects. It depends on the interfaces (ports) defined in the domain layer but has no knowledge of their concrete implementations.
    *   **Example:** `build/application/service.go`

3.  **Infrastructure (The Adapters)**
    *   **Contains:** Concrete implementations of the ports defined in the domain. This includes database repositories, HTTP handlers, and clients for external services.
    *   **Responsibility:** Acts as the bridge between the application core and the outside world. It handles all the technical details of interacting with databases, message queues, APIs, etc.
    *   **Example:** `build/infrastructure/database/repository.go`, `build/infrastructure/http/http_handler.go`

## 2. Workflow for Adding New Features

When adding a new feature, you must follow this "inside-out" approach to respect the architecture:

1.  **Define the Domain:** Start by defining the necessary business models and ports (interfaces) in the `domain` directory. Ask yourself: "What does the core of my application need to do?"
2.  **Implement the Application Service:** Create the application service in the `application` directory. This service will use the ports you defined in the domain to execute the business logic.
3.  **Create Infrastructure Adapters:** In the `infrastructure` directory, create the concrete implementations of your ports.
    *   For database operations, create a new repository in `infrastructure/database`.
    *   For a new API endpoint, create a new handler in `infrastructure/http`.
4.  **Wire the Components:** In `internal/shared/container/container.go`, register your new services and repositories to make them available through dependency injection.

## 3. Addressing Issues and Bugs

Our primary goal when fixing bugs is to address the **root cause**, not just the symptoms. This ensures the stability and integrity of the system.

1.  **Investigate First:** Before writing any code, perform a thorough investigation.
    *   **Check the Logs:** The application logs are the first place to look for errors.
    *   **Verify the Database:** Connect to the database and inspect the schema and data directly. Ensure the tables, columns, and relationships match what the code expects.
    *   **Trace the Request:** Follow the flow of data from the HTTP handler (infrastructure) to the application service and down to the domain to pinpoint exactly where the failure occurs.

2.  **Fix the Foundation:** The solution must be implemented at the correct layer.
    *   If a query is failing due to a schema mismatch, the fix is in the **database repository** (`infrastructure/database`) or, if necessary, a new **database migration**.
    *   If business logic is incorrect, the fix is in the **domain models** or **application service**.
    *   **Never** put business logic in an HTTP handler or database-specific code in a domain entity.

3.  **Avoid Quick Fixes:** Do not implement "quick hacks" that violate architectural boundaries. A fix that solves the immediate problem but introduces technical debt is not acceptable. Always opt for the solution that strengthens the architecture.
