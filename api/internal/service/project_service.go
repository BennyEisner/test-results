package service

import (
	"database/sql"
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
}

// ProjectService provides services related to projects.
type ProjectService struct {
	DB *sql.DB
}

// NewProjectService creates a new ProjectService.
func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{DB: db}
}

// GetDBTestProjectCount returns the count of projects for the DB test endpoint.
func (s *ProjectService) GetDBTestProjectCount() (int, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("database error counting projects: %w", err)
	}
	return count, nil
}

// GetProjectByID fetches a single project by its ID.
func (s *ProjectService) GetProjectByID(id int64) (*models.Project, error) {
	var p models.Project
	err := s.DB.QueryRow("SELECT id, name FROM projects WHERE id = $1", id).Scan(&p.ID, &p.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Let handler decide on 404
		}
		return nil, fmt.Errorf("database error fetching project by ID %d: %w", id, err)
	}
	return &p, nil
}

// GetAllProjects fetches all projects from the database.
func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	rows, err := s.DB.Query("SELECT id, name FROM projects ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("database error fetching all projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, fmt.Errorf("error scanning project: %w", err)
		}
		projects = append(projects, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %w", err)
	}
	return projects, nil
}

// CreateProject creates a new project in the database.
func (s *ProjectService) CreateProject(name string) (*models.Project, error) {
	var p models.Project
	p.Name = name
	err := s.DB.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", name).Scan(&p.ID)
	if err != nil {
		return nil, fmt.Errorf("database error creating project: %w", err)
	}
	return &p, nil
}

// DeleteProject deletes a project by its ID and returns the number of rows affected.
func (s *ProjectService) DeleteProject(id int64) (int64, error) {
	result, err := s.DB.Exec("DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("database error deleting project ID %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error checking delete result for project ID %d: %w", id, err)
	}
	return rowsAffected, nil
}

// UpdateProject updates an existing project's name.
func (s *ProjectService) UpdateProject(id int64, name string) (*models.Project, error) {
	// Check if project exists first
	_, err := s.GetProjectByID(id)
	if err != nil {
		return nil, err // Propagate error (e.g., sql.ErrNoRows for handler to make 404)
	}

	var updatedProject models.Project
	err = s.DB.QueryRow("UPDATE projects SET name = $1 WHERE id = $2 RETURNING id, name", name, id).Scan(&updatedProject.ID, &updatedProject.Name)
	if err != nil {
		return nil, fmt.Errorf("database error updating project ID %d: %w", id, err)
	}
	return &updatedProject, nil
}
