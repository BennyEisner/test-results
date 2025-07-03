# API Service

This is the backend service for the Test Results platform. It is responsible for accepting test result submissions from the CLI and storing them in a database, as well as exposing endpoints for viewing and querying the results.

## Features

* RESTful API using `net/http` and `http.ServeMux`
* JSON-based request and response bodies
* Unit, integration, and end-to-end test structure
* Structured logging middleware
* Clean architecture with separation of concerns (`cmd`, `internal`, `middleware`, `routes`, `config`)
* Interface-based testing with uber-go/mock

## Go Task Runner

### Install

```shell
brew install go-task
```

### Usage

#### Install dependencies

```shell
task deps
```

#### Generate mocks

```shell
task mocks
```

#### Run lint checks

```shell
task lint
```

#### Run tests

```shell
task test
```

#### Run full CI pipeline

```shell
task ci
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
│   ├── models/           # Domain models and interfaces
│   │   └── mocks/        # Generated mocks for interfaces
│   ├── service/          # Business logic
│   │   └── mocks/        # Generated mocks for service interfaces
│   ├── db/               # Database implementations
│   ├── middleware/       # HTTP middleware
│   └── utils/            # Utility functions
├── routes/               # HTTP route definitions
│   └── router.go
├── tests/                # Tests organized by type
│   ├── e2e/
│   ├── integration/
│   └── unit/
├── Dockerfile            # Docker image for deployment
├── .dockerignore         # Exclude unnecessary files from Docker context
├── Taskfile.yml          # Task definitions
├── .golangci.yml         # Linter configuration
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

## Testing with uber-go/mock

This project uses [uber-go/mock](https://github.com/uber-go/mock) for interface-based testing, providing a clean and type-safe alternative to database mocking.

### Why uber-go/mock?

- **No database setup required** - Pure interface mocking
- **Faster tests** - No SQL parsing or database connections
- **Type-safe** - Compile-time checking of method signatures
- **Auto-generated** - Mocks automatically generated from interfaces
- **Always up-to-date** - Regenerates when interfaces change

### Testing Workflow

#### 1. Define Interfaces

Interfaces are already defined in the codebase:

```go
// internal/models/project_repository.go
type ProjectRepository interface {
    GetByID(ctx context.Context, id int64) (*Project, error)
    GetAll(ctx context.Context) ([]*Project, error)
    Create(ctx context.Context, p *Project) error
    // ... other methods
}

// internal/service/project_service.go
type ProjectServiceInterface interface {
    GetProjectByID(id int64) (*models.Project, error)
    CreateProject(name string) (*models.Project, error)
    // ... other methods
}
```

#### 2. Generate Mocks

```shell
task mocks
```

This generates mocks for all interfaces:
- `internal/models/mocks/project_repository_mock.go`
- `internal/service/mocks/project_service_mock.go`
- And many more...

#### 3. Write Tests Using Mocks

```go
package service

import (
    "testing"
    "github.com/BennyEisner/test-results/internal/models"
    mock_models "github.com/BennyEisner/test-results/internal/models/mocks"
    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"
)

func TestProjectService_GetProjectByID(t *testing.T) {
    // Create mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // Create mock repository
    mockRepo := mock_models.NewMockProjectRepository(ctrl)
    
    // Create service with mock
    service := NewProjectService(mockRepo)

    t.Run("success", func(t *testing.T) {
        expectedProject := &models.Project{ID: 1, Name: "Test Project"}
        
        // Set up mock expectations
        mockRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(expectedProject, nil)
        
        // Call the service method
        result, err := service.GetProjectByID(1)
        
        // Assertions
        assert.NoError(t, err)
        assert.Equal(t, expectedProject, result)
    })
}
```

#### 4. Run Tests

```shell
# Run all tests
task test

# Run specific test
go test ./internal/service -v -run TestProjectService_GetProjectByID

# Run all mock-based tests
go test ./internal/service -v -run ".*WithMock"
```

### Mock Generation Details

The `task mocks` command generates mocks for:

- **Repository interfaces**: `ProjectRepository`, `SearchRepository`
- **Service interfaces**: All service interfaces in the `service` package
- **Custom interfaces**: Any other interfaces defined in the codebase

Generated mocks are placed in:
- `internal/models/mocks/` - For repository mocks
- `internal/service/mocks/` - For service mocks

### Key Mock Features

#### Basic Expectations
```go
mockRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(expectedProject, nil)
```

#### Callback Functions
```go
mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
    func(ctx context.Context, p *models.Project) error {
        p.ID = 1  // Set ID on the passed project
        return nil
    })
```

#### Error Scenarios
```go
mockRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, errors.New("not found"))
```

#### Multiple Calls
```go
mockRepo.EXPECT().GetAll(gomock.Any()).Return([]*models.Project{}, nil)
mockRepo.EXPECT().Count(gomock.Any()).Return(0, nil)
```

### Testing Best Practices

1. **Use interfaces** - Define clear interfaces for all dependencies
2. **Generate mocks** - Always regenerate mocks after interface changes
3. **Test business logic** - Focus on testing service logic, not database operations
4. **Use table-driven tests** - For multiple test scenarios
5. **Assert expectations** - Verify that expected methods were called

### Example Test Structure

```go
func TestService_Method(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mock := NewMockInterface(ctrl)
    service := NewService(mock)

    tests := []struct {
        name     string
        setup    func()
        input    interface{}
        expected interface{}
        wantErr  bool
    }{
        {
            name: "success case",
            setup: func() {
                mock.EXPECT().Method(gomock.Any()).Return(expected, nil)
            },
            // ... test data
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            // ... test logic
        })
    }
}
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

