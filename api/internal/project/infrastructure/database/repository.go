package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/project/domain/models"
	"github.com/BennyEisner/test-results/internal/project/domain/ports"
)

// SQLProjectRepository implements the ProjectRepository interface
type SQLProjectRepository struct {
	db *sql.DB
}

// NewSQLProjectRepository creates a new SQL project repository
func NewSQLProjectRepository(db *sql.DB) ports.ProjectRepository {
	return &SQLProjectRepository{db: db}
}

// GetByID retrieves a project by its ID
func (r *SQLProjectRepository) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	query := `SELECT id, name FROM projects WHERE id = $1`

	var project models.Project
	err := r.db.QueryRowContext(ctx, query, id).Scan(&project.ID, &project.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get project by ID: %w", err)
	}

	return &project, nil
}

// GetAll retrieves all projects
func (r *SQLProjectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	query := `SELECT id, name FROM projects ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// GetByName retrieves a project by its name
func (r *SQLProjectRepository) GetByName(ctx context.Context, name string) (*models.Project, error) {
	query := `SELECT id, name FROM projects WHERE name = $1`

	var project models.Project
	err := r.db.QueryRowContext(ctx, query, name).Scan(&project.ID, &project.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get project by name: %w", err)
	}

	return &project, nil
}

// Create creates a new project
func (r *SQLProjectRepository) Create(ctx context.Context, p *models.Project) error {
	query := `INSERT INTO projects (name) VALUES ($1) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, p.Name).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// Update updates an existing project
func (r *SQLProjectRepository) Update(ctx context.Context, id int64, name string) (*models.Project, error) {
	query := `UPDATE projects SET name = $1 WHERE id = $2 RETURNING id, name`

	var project models.Project
	err := r.db.QueryRowContext(ctx, query, name, id).Scan(&project.ID, &project.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &project, nil
}

// Delete deletes a project by its ID
func (r *SQLProjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

// Count returns the total number of projects
func (r *SQLProjectRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM projects`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count projects: %w", err)
	}

	return count, nil
}
