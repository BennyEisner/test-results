package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

// ProjectService implements the domain ProjectService interface
type ProjectService struct {
	projectRepo domain.ProjectRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService(projectRepo domain.ProjectRepository) domain.ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// GetProjectByID retrieves a project by its ID
func (s *ProjectService) GetProjectByID(ctx context.Context, id int64) (*domain.Project, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
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
func (s *ProjectService) GetAllProjects(ctx context.Context) ([]*domain.Project, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all projects: %w", err)
	}

	return projects, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, name string) (*domain.Project, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if project with same name already exists
	existingProject, err := s.projectRepo.GetByName(ctx, name)
	if err == nil && existingProject != nil {
		return nil, domain.ErrDuplicateProject
	}

	project := &domain.Project{
		Name: name,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, id int64, name string) (*domain.Project, error) {
	if id <= 0 || name == "" {
		return nil, domain.ErrInvalidInput
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
		return nil, domain.ErrDuplicateProject
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
		return domain.ErrInvalidInput
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
func (s *ProjectService) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
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
