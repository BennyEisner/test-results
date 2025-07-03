# API Service

This is the backend service for the Test Results platform. It is responsible for accepting test result submissions from the CLI and storing them in a database, as well as exposing endpoints for viewing and querying the results.

## Features

* RESTful API using `net/http` and `http.ServeMux`
* JSON-based request and response bodies
* Unit, integration, and end-to-end test structure
* Structured logging middleware
* Clean architecture with separation of concerns (`cmd`, `internal`, `middleware`, `routes`, `config`)

## Go Task Runner

### Install

```shell
brew install go-task
```

### Usage

#### Install dependencies


task deps

#### Run lint checks

```shell
task lint
```


#### Run tests

```shell
task test
```


## Project Structure

```
api/
├── cmd/
│   └── server/           # Entry point of the application
│       └── main.go
├── config/               # Configuration handling
│   └── config.go
├── internal/             # Private application code
│   ├── handler/          # HTTP handlers
│   ├── model/            # Domain models
│   └── service/          # Business logic
├── middleware/           # HTTP middleware (e.g., logging)
│   └── logging.go
├── pkg/                  # Shared public libraries (empty for now)
├── routes/               # HTTP route definitions
│   └── router.go
├── tests/                # Tests organized by type
│   ├── e2e/
│   ├── integration/
│   └── unit/
├── Dockerfile            # Docker image for deployment
├── .dockerignore         # Exclude unnecessary files from Docker context
└── go.mod                # Go module definition
```

## Getting Started

```sh
cd api
make build    # or go build ./cmd/server
make run      # or go run ./cmd/server
```

## API Endpoints

* `GET /healthz` – Health check
* `GET /readyz` – Readiness check
* `POST /results` – Submit a new test result
* `GET /results` – Query stored results

## Testing

```sh
make test        # Run unit tests
make integration # Run integration tests
make e2e         # Run end-to-end tests
```

## Docker

To build and run the API using Docker:

```sh
docker build -t test-results-api .
docker run -p 8080:8080 test-results-api
```

## License

MIT License. See [LICENSE](../LICENSE) file for details.

## Cyclomatic Complexity Analysis

This project uses [gocyclo](https://github.com/fzipp/gocyclo) to check for functions with high cyclomatic complexity.

### Install gocyclo

```
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
```

Make sure your Go bin directory (e.g., `$HOME/go/bin`) is in your `PATH`.

### Usage

To check for functions with a cyclomatic complexity over 10, run:

```
task cyclo
```

This will report all functions in the `api/` directory with complexity greater than 10.

