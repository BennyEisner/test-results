# Technical Context

## Technology Stack
- **Backend**: Go, Gorilla Mux, Goth (for OAuth)
- **Frontend**: React, TypeScript, Vite, Tanstack Table, Chart.js, React Router, Axios, Bootstrap
- **Database**: PostgreSQL
- **CLI**: Go
- **Containerization**: Docker, Docker Compose
- **Web Server**: Nginx (for serving the frontend)

## Development Setup
- The project is set up as a monorepo.
- The entire stack can be run locally using `docker-compose up`.
- The backend API is written in Go and manages its dependencies with Go Modules (`go.mod`).
- The frontend is a React application built with Vite and uses `npm` for dependency management.
- The database schema is managed with SQL scripts located in the `/db` directory.

## Build and Deployment
- The frontend is built using `npm run build`, which creates a production-ready bundle.
- The Go backend is built into a binary.
- Both the frontend and backend are containerized using Dockerfiles.
- Nginx is used as a reverse proxy and to serve the static frontend files in the production Docker environment.

## Technical Constraints
- The project relies on the specific versions of dependencies listed in `go.mod` and `package.json`.
- The database schema is designed for PostgreSQL.
- The CLI is designed to parse JUnit and ReadyAPI XML formats.
