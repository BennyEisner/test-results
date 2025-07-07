package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/project/domain/models"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Project, error)
	GetAll(ctx context.Context) ([]*models.Project, error)
	GetByName(ctx context.Context, name string) (*models.Project, error)
	Create(ctx context.Context, p *models.Project) error
	Update(ctx context.Context, id int64, name string) (*models.Project, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}

// ProjectService defines the interface for project business logic
type ProjectService interface {
	GetProject(ctx context.Context, id int64) (*models.Project, error)
	GetAllProjects(ctx context.Context) ([]*models.Project, error)
	CreateProject(ctx context.Context, name string) (*models.Project, error)
	UpdateProject(ctx context.Context, id int64, name string) (*models.Project, error)
	DeleteProject(ctx context.Context, id int64) error
}
