# Hexagonal Architecture Implementation

This document explains how the test-results API has been refactored to follow hexagonal architecture (ports and adapters) principles.

## Overview

Hexagonal architecture separates the application into three main layers:

1. **Domain Layer** - Core business logic and entities
2. **Application Layer** - Use cases and application services
3. **Infrastructure Layer** - External concerns (databases, HTTP, etc.)

## Directory Structure

```
api/internal/
├── domain/                    # Domain layer (core business logic)
│   └── ports.go              # Domain models and port interfaces
├── application/              # Application layer (use cases)
│   └── project_service.go    # Application services implementing domain ports
├── infrastructure/           # Infrastructure layer (adapters)
│   ├── database/            # Database adapters (secondary/driven)
│   │   └── project_repository.go
│   ├── http/                # HTTP adapters (primary/driving)
│   │   └── project_handler.go
│   ├── container.go         # Dependency injection container
│   └── router.go            # HTTP router
└── cmd/server/
    └── example_hexagonal.go # Example of hexagonal setup
```

## Key Components

### 1. Domain Layer (`internal/domain/`)

The domain layer contains:
- **Domain Models**: Core business entities (Project, Build, TestSuite, etc.)
- **Input Ports**: Interfaces that define what the application can do
- **Output Ports**: Interfaces that define what the application needs from external systems
- **Domain Errors**: Business-specific error types

```go
// Domain model
type Project struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

// Input port (primary/driving adapter interface)
type ProjectService interface {
    GetProjectByID(ctx context.Context, id int64) (*Project, error)
    CreateProject(ctx context.Context, name string) (*Project, error)
    // ... other methods
}

// Output port (secondary/driven adapter interface)
type ProjectRepository interface {
    GetByID(ctx context.Context, id int64) (*Project, error)
    Create(ctx context.Context, p *Project) error
    // ... other methods
}
```

### 2. Application Layer (`internal/application/`)

The application layer contains use cases that implement the domain input ports:

```go
type ProjectService struct {
    projectRepo domain.ProjectRepository
}

func (s *ProjectService) CreateProject(ctx context.Context, name string) (*domain.Project, error) {
    // Business logic validation
    if name == "" {
        return nil, domain.ErrInvalidInput
    }
    
    // Check for duplicates
    existingProject, err := s.projectRepo.GetByName(ctx, name)
    if err == nil && existingProject != nil {
        return nil, domain.ErrDuplicateProject
    }
    
    // Create project
    project := &domain.Project{Name: name}
    if err := s.projectRepo.Create(ctx, project); err != nil {
        return nil, fmt.Errorf("failed to create project: %w", err)
    }
    
    return project, nil
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

#### Database Adapters (Secondary/Driven)

Implement domain output ports:

```go
type SQLProjectRepository struct {
    db *sql.DB
}

func (r *SQLProjectRepository) Create(ctx context.Context, p *domain.Project) error {
    query := `INSERT INTO projects (name) VALUES ($1) RETURNING id`
    return r.db.QueryRowContext(ctx, query, p.Name).Scan(&p.ID)
}
```

#### HTTP Adapters (Primary/Driving)

Implement domain input ports and handle HTTP concerns:

```go
type ProjectHandler struct {
    projectService domain.ProjectService
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
    var request struct {
        Name string `json:"name"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        respondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    ctx := r.Context()
    project, err := h.projectService.CreateProject(ctx, request.Name)
    if err != nil {
        switch err {
        case domain.ErrInvalidInput:
            respondWithError(w, http.StatusBadRequest, "project name is required")
        case domain.ErrDuplicateProject:
            respondWithError(w, http.StatusConflict, "project already exists")
        default:
            respondWithError(w, http.StatusInternalServerError, "internal error")
        }
        return
    }
    
    respondWithJSON(w, http.StatusCreated, project)
}
```

### 4. Dependency Injection Container

The container wires up all dependencies:

```go
type Container struct {
    DB *sql.DB
    
    // Repositories (Secondary/Driven Adapters)
    ProjectRepository domain.ProjectRepository
    
    // Application Services (Use Cases)
    ProjectService domain.ProjectService
    
    // HTTP Handlers (Primary/Driving Adapters)
    ProjectHandler *http.ProjectHandler
}

func NewContainer(db *sql.DB) *Container {
    container := &Container{DB: db}
    
    // Wire up dependencies
    container.ProjectRepository = database.NewSQLProjectRepository(db)
    container.ProjectService = application.NewProjectService(container.ProjectRepository)
    container.ProjectHandler = http.NewProjectHandler(container.ProjectService)
    
    return container
}
```

## Benefits of Hexagonal Architecture

### 1. **Testability**
- Domain logic can be tested without external dependencies
- Easy to mock repositories and services
- Clear separation of concerns

### 2. **Flexibility**
- Easy to swap implementations (e.g., different databases)
- Can add new adapters without changing core logic
- Framework-agnostic domain layer

### 3. **Maintainability**
- Clear boundaries between layers
- Business logic is isolated from infrastructure concerns
- Easy to understand and modify

### 4. **Scalability**
- Can easily add new use cases
- Can implement different interfaces for different clients
- Clear dependency direction

## Migration Strategy

To migrate the existing codebase to hexagonal architecture:

1. **Extract Domain Models**: Move business entities to `internal/domain/`
2. **Define Ports**: Create interfaces for all external dependencies
3. **Implement Application Services**: Move business logic to application layer
4. **Create Adapters**: Implement ports for databases, HTTP, etc.
5. **Update Dependencies**: Wire everything together with dependency injection

## Example Usage

```go
// In main.go or similar
func main() {
    // Setup database
    db := connectDB()
    
    // Create container with all dependencies
    container := infrastructure.NewContainer(db)
    
    // Create router
    router := infrastructure.NewRouter(container)
    
    // Start server
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    server.ListenAndServe()
}
```

## Testing

With hexagonal architecture, testing becomes much easier:

```go
func TestProjectService_CreateProject(t *testing.T) {
    // Create mock repository
    mockRepo := &MockProjectRepository{}
    
    // Create service with mock
    service := application.NewProjectService(mockRepo)
    
    // Test business logic
    project, err := service.CreateProject(context.Background(), "Test Project")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "Test Project", project.Name)
}
```

## Next Steps

To complete the hexagonal architecture implementation:

1. **Extend Domain Layer**: Add all remaining domain models and ports
2. **Implement All Services**: Create application services for all use cases
3. **Add More Adapters**: Implement repositories for all entities
4. **Add Validation**: Implement domain validation rules
5. **Add Events**: Consider domain events for complex workflows
6. **Add CQRS**: Consider Command Query Responsibility Segregation for complex queries

This architecture provides a solid foundation for a maintainable, testable, and scalable application. 