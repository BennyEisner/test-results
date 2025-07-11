# Test Results

## Overview

This project is a full-stack application designed to:

* Accept test results in **JUnit** and **ReadyAPI** formats from various CI systems (GitHub Actions, Travis, Jenkins) via a **Go-based CLI**
* Post the test results to a **Go-based REST API**
* Store and manage test results in a **PostgreSQL** database
* Display test results over time using a **React frontend dashboard**

The project provides experience in:

* Backend development with Go
* Frontend development with React
* Working with RESTful APIs
* Using relational databases
* Unit, integration, and end-to-end testing

---

## System Architecture Diagram

```
+------------------------+        +----------------+        +---------------------+
|   CLI / Build Agent    +------->+     REST API    +------->+  PostgreSQL DB       |
| (junit/readyapi input) |        | (Go-based)      |        | (test results store) |
+------------------------+        +--------+-------+        +-----------+---------+
                                            |
                                            v
                                    +---------------+
                                    | React Frontend|
                                    +---------------+
```

---

## Getting Started Guide

### Install Required Tools

* Go
* Node.js & npm
* Docker

### Clone the Repository

```bash
git clone https://github.com/your-org/fullstack-test-tracker.git
cd fullstack-test-tracker
```

### Start the Project

```bash
docker-compose up -d
cd api && go run main.go
cd ../frontend && npm install && npm start
```

### Authentication Setup (Optional)

If you want to use the authentication system:

```bash
# Run the authentication setup script
./scripts/setup-auth-dev.sh

# Follow the prompts to set up GitHub OAuth2
# See docs/auth-development.md for detailed instructions
```

### Seed Database

Insert test data manually or using a seed script.

### Run CLI

Point the CLI to sample JUnit/ReadyAPI result files and post them to the API.

---

## Contribution Guide

* Use feature branches (see [Feature Branch Workflow](docs/feature-branch.md))
* Format code with `gofmt`, Prettier, etc.
* Write tests for new code
* Use pre-commit hooks if available

For more detailed contribution guidelines, see [Contributing Guide](docs/contributing.md).

---

## Project Layout (Monorepo)

```text
fullstack-test-tracker/
├── api/               # Go REST API
│   ├── main.go
│   ├── handlers/
│   ├── models/
│   └── tests/
├── cli/               # Go CLI
│   ├── main.go
│   ├── parsers/
│   └── uploader/
├── frontend/          # React app
│   ├── public/
│   ├── src/
│   └── tests/
├── db/                # SQL migrations
│   └── schema.sql
├── docker-compose.yml
├── Makefile
├── README.md
└── docs/
```

---

## Milestone Plan

### Week 1: Setup & Orientation

* Install dependencies
* Run services locally
* Understand architecture and DB schema

### Week 2: Backend API

* Create endpoints for tests, suites, runs
* Connect to PostgreSQL
* Unit test handlers and models

### Week 3: CLI Tool

* Parse JUnit & ReadyAPI XML
* Post data to API
* Include metadata flags (CI type, run ID, etc.)

### Week 4: Frontend App

* Display list of test runs
* Filter by CI system and test status
* Drill-down into run details

### Week 5+: Testing & CI

* Add integration and E2E tests
* Set up CI pipelines (GitHub Actions)
* Optional: Add authentication

---

## Story Descriptions
* [ ] Setup local dev environment and verify Docker + Go + Node installation
* [ ] Run Postgres with Docker and connect via psql or GUI
* [ ] Create initial DB schema for test results and test suites
* [ ] Build a basic Go server with a health check endpoint
* [ ] Add API route to submit test run metadata (test name, suite, CI system, timestamp)
* [ ] Implement unit tests for API handlers
* [ ] Write CLI parser for JUnit XML files
* [ ] Send parsed JUnit data to the API using HTTP POST
* [ ] Expand CLI to support ReadyAPI XML format
* [ ] Add metadata flags to CLI for CI source, run ID, and suite ID
* [ ] Build frontend layout with a basic test run list view
* [ ] Fetch and render test data from the API using React Query
* [ ] Add filter UI for CI system and test status
* [ ] Build detail page for a single test run
* [ ] Implement integration tests for API + DB flow
* [ ] Add E2E test
* [ ] Set up GitHub Actions CI pipeline for Go tests and React build

---

## Key Concepts to Learn

* REST API design
* Git workflows
* Relational modeling
* Unit, integration, and E2E testing
* Containerized dev environments with Docker

---

## Documentation

Additional documentation is available in the `docs/` directory:

* [Architectural Design Records (ADR)](docs/ADR.md) - Design decisions and architectural considerations
* [API Documentation](docs/api.md) - Detailed API reference and usage
* [Database Documentation](docs/db.md) - Database schema and migration information
* [Authentication Development Guide](docs/auth-development.md) - Local development setup and usage for authentication system
* [Contributing Guide](docs/contributing.md) - How to contribute to the project
* [Feature Branch Workflow](docs/feature-branch.md) - Git workflow for feature development

## Optional

We can generate a GitHub repo scaffold if needed to jumpstart development.

## Local Development URLs

When running the stack locally with Docker Compose, use the following URLs to access the services:

- **Frontend (UI):** [http://localhost:8088](http://localhost:8088)
  - This serves the web UI via Nginx.

- **Backend API:** [http://localhost:8080](http://localhost:8080)
  - The API root. Endpoints are available under `/api`, e.g.:
    - [http://localhost:8080/api/projects](http://localhost:8080/api/projects)
    - [http://localhost:8080/api/builds](http://localhost:8080/api/builds)

- **Swagger API Documentation:** [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
  - Interactive OpenAPI docs for the backend API.

- **Health Checks:**
  - [http://localhost:8080/readyz](http://localhost:8080/readyz) (readiness)
  - [http://localhost:8080/livez](http://localhost:8080/livez) (liveness)
  - [http://localhost:8080/healthz](http://localhost:8080/healthz) (comprehensive health)

## Using the CLI to Post Results to the API

The CLI can be used to submit test results to the backend API. Make sure the API is running and accessible (see URLs above).

### Example: Post JUnit Results

From the project root, run:

```sh
cd cli
./cli post --file <path-to-junit-xml> --project <project-name> --suite <suite-name> --api-url http://localhost:8080/api
```

- `--file` (required): Path to the JUnit XML file to upload.
- `--project` (required): Name of the project to associate the results with.
- `--suite` (required): Name of the test suite.
- `--api-url` (optional): The base URL of the API (default: `http://localhost:8080/api`).

### Example: Using Dockerized CLI

If you want to run the CLI in a container:

```sh
docker run --rm -v $(pwd)/results:/results cli-image-name post --file /results/junit.xml --project MyProject --suite MySuite --api-url http://host.docker.internal:8080/api
```

### Notes
- Ensure the API is running and accessible at the specified `--api-url`.
- The CLI may require configuration (see `cli/README.md` for more details).
- For more CLI commands and options, run:
  ```sh
  ./cli --help
  ```
