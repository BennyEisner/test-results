package db

import (
	"context"
	"database/sql"
	"github.com/BennyEisner/test-results/api/models"
)

type SQLProjectRepository struct {
	DB *sql.DB
}

func NewSQLProjectRepository(db *sql.DB) *SQLProjectRepository {
	return &SQLProjectRepository{DB: db}
}

func (r *SQLProjectRepository) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, name FROM projects WHERE id = $1`, id)
	p := &models.Project{}
	if err := p.ScanFromRow(row); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *SQLProjectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		p := &models.Project{}
		if err := p.ScanFromRows(rows); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *SQLProjectRepository) Create(ctx context.Context, p *models.Project) error {
	return p.Insert(r.DB)
}

func (r *SQLProjectRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}
