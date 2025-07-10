package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/project/domain"
	"github.com/BennyEisner/test-results/internal/project/domain/models"
	"github.com/BennyEisner/test-results/internal/project/domain/ports"
)

// ProjectService implements the domain ProjectService interface
type ProjectService struct {
	projectRepo ports.ProjectRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService(projectRepo ports.ProjectRepository) ports.ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// GetProject retrieves a project by its ID
func (s *ProjectService) GetProject(ctx context.Context, id int64) (*models.Project, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidProjectName
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by ID %d: %w", id, err)
	}

	if project == nil {
		return nil, domain.ErrProjectNotFound
	}

	return project, nil
}

// GetAllProjects retrieves all projects
func (s *ProjectService) GetAllProjects(ctx context.Context) ([]*models.Project, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all projects: %w", err)
	}

	return projects, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, name string) (*models.Project, error) {
	if name == "" {
		return nil, domain.ErrInvalidProjectName
	}

	// Check if project with same name already exists
	existingProject, err := s.projectRepo.GetByName(ctx, name)
	if err == nil && existingProject != nil {
		return nil, domain.ErrProjectAlreadyExists
	}

	project := &models.Project{
		Name: name,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, id int64, name string) (*models.Project, error) {
	if id <= 0 || name == "" {
		return nil, domain.ErrInvalidProjectName
	}

	// Check if project exists
	existingProject, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}
	if existingProject == nil {
		return nil, domain.ErrProjectNotFound
	}

	// Check if new name conflicts with existing project
	conflictingProject, err := s.projectRepo.GetByName(ctx, name)
	if err == nil && conflictingProject != nil && conflictingProject.ID != id {
		return nil, domain.ErrProjectAlreadyExists
	}

	updatedProject, err := s.projectRepo.Update(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return updatedProject, nil
}

// DeleteProject deletes a project by its ID
func (s *ProjectService) DeleteProject(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidProjectName
	}

	// Check if project exists
	existingProject, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %w", err)
	}
	if existingProject == nil {
		return domain.ErrProjectNotFound
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// GetProjectByName retrieves a project by its name
func (s *ProjectService) GetProjectByName(ctx context.Context, name string) (*models.Project, error) {
	if name == "" {
		return nil, domain.ErrInvalidProjectName
	}

	project, err := s.projectRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by name: %w", err)
	}

	if project == nil {
		return nil, domain.ErrProjectNotFound
	}

	return project, nil
}
