package models

import "context"

type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (*Project, error)
	GetAll(ctx context.Context) ([]*Project, error)
	Create(ctx context.Context, p *Project) error
	Delete(ctx context.Context, id int64) error
}

