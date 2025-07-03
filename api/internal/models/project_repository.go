package models

import "context"

type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (*Project, error)
	GetAll(ctx context.Context) ([]*Project, error)
	GetByName(ctx context.Context, name string) (*Project, error)
	Create(ctx context.Context, p *Project) error
	Update(ctx context.Context, id int64, name string) (*Project, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}
