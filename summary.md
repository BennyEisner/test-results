# Repository Summary

## 1. Project Overview

This project is a full-stack test results tracking system designed to collect, store, and visualize test data from CI/CD pipelines. It consists of a Go-based CLI for ingesting JUnit/ReadyAPI reports, a Go REST API for data processing, a PostgreSQL database for storage, and a React frontend for visualization. The system provides a comprehensive solution for monitoring test outcomes and trends over time.

## 2. Repository Structure

The repository is a monorepo organized into the following main directories:

-   **`api/`**: Contains the Go-based REST API, which follows a hexagonal architecture pattern. It is responsible for handling data from the CLI and serving it to the frontend.
-   **`cli/`**: A Go-based command-line interface used to parse and upload test results to the API.
-   **`frontend/`**: A React application built with Vite that provides a user interface for viewing test results.
-   **`db/`**: Contains SQL schema definitions and migration scripts for the PostgreSQL database.
-   **`docs/`**: Includes project documentation, such as architectural decision records (ADRs) and API documentation.
-   **`scripts/`**: Contains utility scripts for the project.

## 3. Key Components

-   **`api/internal/shared/container/container.go`**: The dependency injection container for the API, which wires together all services, repositories, and handlers.
-   **`api/cmd/server/main.go`**: The main entry point for the Go API server.
-   **`cli/main.go`**: The main entry point for the CLI application.
-   **`frontend/src/App.tsx`**: The main React component that sets up the application's routing.
-   **`db/schema.sql`**: Defines the database schema, including tables for projects, test suites, builds, test cases, and failures.
-   **`docker-compose.yml`**: Configures the services for local development, including the database, API, and frontend.

## 4. Technologies Used

-   **Backend**: Go, `net/http` for the web server, `lib/pq` for PostgreSQL interaction, and `testify` for testing.
-   **Frontend**: React, TypeScript, Vite, React Router for routing, Axios for API requests, and Chart.js for data visualization.
-   **Database**: PostgreSQL.
-   **CLI**: Go with the `cobra` library for building the command-line interface.
-   **Build & Deployment**: Docker and Docker Compose for containerization and local development.

## 5. Notable Features

-   **Hexagonal Architecture**: The API is structured using a hexagonal (ports and adapters) architecture, which separates business logic from infrastructure concerns, improving testability and maintainability.
-   **Monorepo Structure**: The project is organized as a monorepo, which simplifies dependency management and cross-component changes.
-   **Comprehensive Test Coverage**: The project includes unit tests for the API services, demonstrating a commitment to code quality.
-   **CI/CD Integration**: The system is designed to integrate with various CI systems like GitHub Actions, Travis, and Jenkins.

## 6. Testing and Documentation

-   **Testing**: The API has unit tests for its services, using the `testify` library for assertions and mocking. The hexagonal architecture facilitates testing by allowing repositories to be easily mocked.
-   **Documentation**: The `docs/` directory contains detailed documentation, including an explanation of the hexagonal architecture, API documentation, and ADRs. The API also includes Swagger documentation for its endpoints.

## 7. Deployment

The project is designed to be run using Docker and Docker Compose. The `docker-compose.yml` file defines the services for the database, API, and frontend, making it easy to set up a local development environment. The `Makefile` provides convenient commands for running the different parts of the application.
