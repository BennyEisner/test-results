package service

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/models"
)

// ProjectServiceInterface defines the interface for project service operations.
type ProjectServiceInterface interface {
	GetProjectByID(id int64) (*models.Project, error)
	GetAllProjects() ([]models.Project, error)
	CreateProject(name string) (*models.Project, error)
	UpdateProject(id int64, name string) (*models.Project, error)
	DeleteProject(id int64) (int64, error) // Returns rows affected
	GetDBTestProjectCount() (int, error)
	GetProjectByName(name string) (*models.Project, error)
}

// ProjectService provides services related to projects.
type ProjectService struct {
	Repo models.ProjectRepository
}

// NewProjectService creates a new ProjectService.
func NewProjectService(repo models.ProjectRepository) ProjectServiceInterface {
	return &ProjectService{Repo: repo}
}

// GetDBTestProjectCount returns the count of projects for the DB test endpoint.
func (s *ProjectService) GetDBTestProjectCount() (int, error) {
	return s.Repo.Count(context.Background())
}

// GetProjectByID fetches a single project by its ID.
func (s *ProjectService) GetProjectByID(id int64) (*models.Project, error) {
	return s.Repo.GetByID(context.Background(), id)
}

// GetAllProjects fetches all projects from the database.
func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	projects, err := s.Repo.GetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("repository error fetching all projects: %w", err)
	}
	// Convert []*Project to []Project
	var result []models.Project
	for _, p := range projects {
		result = append(result, *p)
	}
	return result, nil
}

// CreateProject creates a new project in the database.
func (s *ProjectService) CreateProject(name string) (*models.Project, error) {
	p := &models.Project{Name: name}
	if err := s.Repo.Create(context.Background(), p); err != nil {
		return nil, fmt.Errorf("repository error creating project: %w", err)
	}
	return p, nil
}

// DeleteProject deletes a project by its ID and returns the number of rows affected.
func (s *ProjectService) DeleteProject(id int64) (int64, error) {
	if err := s.Repo.Delete(context.Background(), id); err != nil {
		return 0, fmt.Errorf("repository error deleting project ID %d: %w", id, err)
	}
	// We can't get rows affected from the repo, so just return 1 for success
	return 1, nil
}

// UpdateProject updates an existing project's name.
func (s *ProjectService) UpdateProject(id int64, name string) (*models.Project, error) {
	return s.Repo.Update(context.Background(), id, name)
}

// GetProjectByName fetches a single project by its name.
func (s *ProjectService) GetProjectByName(name string) (*models.Project, error) {
	return s.Repo.GetByName(context.Background(), name)
}
