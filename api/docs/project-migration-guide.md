# Project Domain Migration Guide

This guide explains how the Projects domain has been migrated to use hexagonal architecture and how to use both the old and new implementations.

## Overview

The Projects domain has been successfully migrated to hexagonal architecture with the following components:

- **Domain Layer**: Core business entities and interfaces
- **Application Layer**: Business logic and use cases
- **Infrastructure Layer**: Database and HTTP adapters
- **Dependency Injection**: Clean wiring of all components

## API Endpoints

### Old Implementation (v1)
The existing API endpoints continue to work unchanged:

- `GET /api/projects` - Get all projects
- `POST /api/projects` - Create a new project
- `GET /api/projects/{id}` - Get project by ID
- `PATCH /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project

### New Implementation (v2)
New hexagonal architecture endpoints are available at:

- `GET /api/v2/projects` - Get all projects (hexagonal)
- `POST /api/v2/projects` - Create a new project (hexagonal)
- `GET /api/v2/projects/{id}` - Get project by ID (hexagonal)
- `PATCH /api/v2/projects/{id}` - Update project (hexagonal)
- `DELETE /api/v2/projects/{id}` - Delete project (hexagonal)

## Key Differences

### 1. **Better Error Handling**
The hexagonal implementation provides more specific error types:

```json
// Old implementation
{
  "error": "Database error: project not found"
}

// New implementation
{
  "error": "project not found"
}
```

### 2. **Domain-Specific Errors**
The new implementation uses domain-specific error codes:

- `PROJECT_NOT_FOUND` - Project doesn't exist
- `INVALID_INPUT` - Invalid input parameters
- `DUPLICATE_PROJECT` - Project with same name already exists

### 3. **Consistent Response Format**
All responses follow a consistent format with proper HTTP status codes.

## Testing the Migration

### 1. **Start the Server**
```bash
cd api
go run cmd/server/main.go
```

### 2. **Test Old Endpoints**
```bash
# Get all projects (old implementation)
curl http://localhost:8080/api/projects

# Create a project (old implementation)
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Project"}'
```

### 3. **Test New Endpoints**
```bash
# Get all projects (new hexagonal implementation)
curl http://localhost:8080/api/v2/projects

# Create a project (new hexagonal implementation)
curl -X POST http://localhost:8080/api/v2/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Project V2"}'
```

## Code Structure

### Domain Layer (`internal/domain/ports.go`)
```go
// Domain model
type Project struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

// Input port (service interface)
type ProjectService interface {
    GetProjectByID(ctx context.Context, id int64) (*Project, error)
    CreateProject(ctx context.Context, name string) (*Project, error)
    // ... other methods
}

// Output port (repository interface)
type ProjectRepository interface {
    GetByID(ctx context.Context, id int64) (*Project, error)
    Create(ctx context.Context, p *Project) error
    // ... other methods
}
```

### Application Layer (`internal/application/project_service.go`)
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

### Infrastructure Layer

#### Database Adapter (`internal/infrastructure/database/project_repository.go`)
```go
type SQLProjectRepository struct {
    db *sql.DB
}

func (r *SQLProjectRepository) Create(ctx context.Context, p *domain.Project) error {
    query := `INSERT INTO projects (name) VALUES ($1) RETURNING id`
    return r.db.QueryRowContext(ctx, query, p.Name).Scan(&p.ID)
}
```

#### HTTP Adapter (`internal/infrastructure/http/project_handler.go`)
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

## Benefits of the Migration

### 1. **Testability**
- Domain logic can be tested without external dependencies
- Easy to mock repositories and services
- Clear separation of concerns

### 2. **Maintainability**
- Clear boundaries between layers
- Business logic is isolated from infrastructure concerns
- Easy to understand and modify

### 3. **Flexibility**
- Easy to swap implementations (e.g., different databases)
- Can add new adapters without changing core logic
- Framework-agnostic domain layer

### 4. **Error Handling**
- Domain-specific error types
- Consistent error responses
- Better debugging and monitoring

## Running Tests

### Unit Tests
```bash
# Test the application layer
go test ./internal/application -v

# Test the infrastructure layer
go test ./internal/infrastructure -v
```

### Integration Tests
```bash
# Test the entire hexagonal implementation
go test ./internal/... -v -tags=integration
```

## Migration Strategy

### Phase 1: Parallel Implementation ✅
- Both old and new implementations coexist
- New endpoints available at `/api/v2/`
- No breaking changes to existing clients

### Phase 2: Gradual Migration
- Update clients to use new endpoints
- Monitor performance and error rates
- Validate business logic correctness

### Phase 3: Cleanup
- Remove old implementation
- Update all clients to use new endpoints
- Clean up deprecated code

## Monitoring and Validation

### 1. **Compare Responses**
Ensure both implementations return the same data:

```bash
# Compare old vs new
curl http://localhost:8080/api/projects > old_response.json
curl http://localhost:8080/api/v2/projects > new_response.json
diff old_response.json new_response.json
```

### 2. **Error Scenarios**
Test error handling in both implementations:

```bash
# Test invalid input
curl -X POST http://localhost:8080/api/v2/projects \
  -H "Content-Type: application/json" \
  -d '{"name": ""}'

# Test duplicate project
curl -X POST http://localhost:8080/api/v2/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "Duplicate Project"}'
curl -X POST http://localhost:8080/api/v2/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "Duplicate Project"}'
```

### 3. **Performance Comparison**
Monitor response times and resource usage:

```bash
# Benchmark old implementation
ab -n 1000 -c 10 http://localhost:8080/api/projects

# Benchmark new implementation
ab -n 1000 -c 10 http://localhost:8080/api/v2/projects
```

## Next Steps

1. **Extend to Other Domains**: Apply the same pattern to Build, TestSuite, etc.
2. **Add Validation**: Implement domain validation rules
3. **Add Events**: Consider domain events for complex workflows
4. **Add CQRS**: Consider Command Query Responsibility Segregation
5. **Add Monitoring**: Add metrics and logging for the new implementation

## Troubleshooting

### Common Issues

1. **Import Errors**: Ensure all dependencies are installed
   ```bash
   go mod tidy
   ```

2. **Database Connection**: Verify database configuration
   ```bash
   export POSTGRES_DSN="host=localhost port=5432 user=postgres password=password dbname=test_results sslmode=disable"
   ```

3. **Port Conflicts**: Ensure port 8080 is available
   ```bash
   lsof -i :8080
   ```

### Getting Help

- Check the logs for detailed error messages
- Review the hexagonal architecture documentation
- Run the tests to identify specific issues
- Compare old vs new implementation behavior

The migration provides a solid foundation for a maintainable, testable, and scalable application while maintaining backward compatibility. 