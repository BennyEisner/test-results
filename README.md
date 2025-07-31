# Test Results Reporting and Analysis Platform

## 1. Overview

This project is a comprehensive, full-stack application designed to provide a centralized system for ingesting, storing, and analyzing test results from various CI/CD pipelines. It empowers development and QA teams with a powerful dashboard to track test performance, identify trends, and improve software quality.

### Key Features

*   **Multi-format Test Result Ingestion**: A robust Go-based CLI accepts test results in standard formats like **JUnit** and **ReadyAPI**, making it easy to integrate with CI systems like GitHub Actions, Travis CI, and Jenkins.
*   **RESTful API**: A high-performance Go backend exposes a RESTful API for data ingestion, querying, and management. It is built using a clean, maintainable Hexagonal Architecture.
*   **Relational Data Store**: Test results, builds, suites, and projects are stored in a **PostgreSQL** database, providing a reliable and scalable data foundation.
*   **Interactive Dashboard**: A dynamic **React** frontend provides an interactive dashboard for visualizing test data, exploring build histories, and analyzing test suite performance over time.
*   **Authentication**: Secure access to the platform is managed through an authentication system supporting both **OAuth2** (e.g., GitHub) for web users and **API keys** for CLI and automated clients.

---

## 2. Architecture

The system is designed as a set of containerized services that work together to provide a seamless experience, from data submission to visualization.

### Backend Architecture (API)

The Go backend follows the **Hexagonal Architecture**, ensuring a clean separation of concerns and making the system easier to test and maintain.

*   **Domain Core**: Contains the core business logic and models, with no dependencies on external technologies.
*   **Application Layer**: Orchestrates the business logic by using the domain core.
*   **Infrastructure Layer**: Contains all external-facing components, such as:
    *   **HTTP Handlers**: Expose the API endpoints.
    *   **Database Repositories**: Implement data persistence using PostgreSQL.
    *   **Authentication Middleware**: Secures the API.
*   **Backend for Frontend (BFF)**: The API is designed as a BFF, providing data tailored specifically for the needs of the React frontend.

### Frontend Architecture

The frontend single-page application built with **React** and **TypeScript**.

*   **Routing**: **React Router** is used for all client-side navigation.
*   **State Management**: Application-wide state, such as authentication status and dashboard configurations, is managed using React's **Context API** (`AuthContext`, `DashboardContext`).
*   **UI Components**: The UI is built with a combination of custom components and libraries, including:
    *   **Tanstack Table** for data-rich tables.
    *   **Chart.js** for data visualization.
    *   **React-Grid-Layout** for customizable dashboard widgets.
*   **Styling**: A combination of CSS Modules, SCSS, and global stylesheets are used for styling.

### Data Model

The PostgreSQL database schema is designed to store test data hierarchically:

*   `projects` -> `test_suites` -> `builds` -> `test_cases` -> `build_test_case_executions`
*   Authentication tables (`auth_users`, `auth_api_keys`, `auth_sessions`) manage user and client access.

### Containerization

The entire application is containerized using **Docker** and orchestrated with **Docker Compose**. This provides a consistent, isolated development environment and simplifies deployment. The `docker-compose.yml` file defines three main services: `db`, `api`, and `frontend`.

---

## 3. Technology Stack

| Category      | Technologies                                                              |
|---------------|---------------------------------------------------------------------------|
| **Backend**   | Go, Gorilla/Mux, Goth, Cobra, `sqlx`                                        |
| **Frontend**  | React, TypeScript, Vite, React Router, Tanstack Table, Chart.js, SCSS       |
| **Database**  | PostgreSQL                                                                |
| **CLI**       | Go, Cobra                                                                 |
| **DevOps**    | Docker, Docker Compose, Nginx                                             |


