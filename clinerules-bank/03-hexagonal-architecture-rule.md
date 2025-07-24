# Hexagonal Architecture Rule

This rule enforces the principles of Hexagonal Architecture (Ports and Adapters) within the Go backend. It ensures a clean separation of concerns between the core application logic and infrastructure-specific details.

## 1. Core Principle: Separation of Concerns

The backend code is organized into three main layers:

-   **`domain`**: The core of the application. It contains the business models, interfaces (ports), and domain-specific errors. It has no dependencies on any other layer.
-   **`application`**: The layer that orchestrates the business logic. It implements the `domain` ports and contains the application services. It depends only on the `domain` layer.
-   **`infrastructure`**: The outermost layer that contains all the details about how the application interacts with the outside world. This includes HTTP handlers, database repositories, and other external services. It depends on the `application` and `domain` layers.

## 2. HTTP Handler Responsibilities

HTTP handlers, located in the `infrastructure/http` directory, have a strict set of responsibilities:

1.  **Parse the Request**: Extract data from the HTTP request, including URL parameters, query parameters, and the request body.
2.  **Validate the Request**: Perform basic validation on the incoming data. For more complex validation, delegate to the application service.
3.  **Call the Application Service**: Invoke the appropriate method on the application service, passing the validated data.
4.  **Handle Errors**: Check for errors returned from the application service and translate them into appropriate HTTP status codes and error responses.
5.  **Format the Response**: If the service call is successful, format the returned data into the appropriate response format (e.g., JSON) and write it to the HTTP response.

**Crucially, HTTP handlers MUST NOT contain any business logic.** All business logic must reside within the `application` and `domain` layers.

## 3. Example Workflow

```
flowchart TD
    A[HTTP Request] --> B{HTTP Handler};
    B --> C{Application Service};
    C --> D{Domain Logic};
    D --> E[Database/External Service];
    E --> D;
    D --> C;
    C --> B;
    B --> F[HTTP Response];
```

## 4. Directory Structure

Each feature or component within the `internal` directory should follow this structure:

```
<component>/
├── application/
│   └── service.go
├── domain/
│   ├── errors/
│   │   └── errors.go
│   ├── models/
│   │   └── models.go
│   └── ports/
│       └── ports.go
└── infrastructure/
    ├── database/
    │   └── repository.go
    └── http/
        └── http_handler.go
