# Technical Context

## Backend

- **Language:** Go (version 1.21+)
- **Framework:** Standard library, with a hexagonal architecture.
- **Dependencies:** `golangci-lint` for linting. Other dependencies are managed with Go modules (`go.mod`).
- **Database:** PostgreSQL

## Frontend

- **Framework:** React
- **Language:** TypeScript
- **Build Tool:** Vite
- **Package Manager:** npm
- **State Management:** React Context for authentication.

## Infrastructure

- **Containerization:** Docker and Docker Compose for local development environment.
- **CI/CD:** Jenkins is mentioned as a user of API keys.

## Authentication

- **OAuth2:** GitHub is used for local development. Okta is planned for production.
- **API Keys:** Used for CLI tools and CI/CD systems.
