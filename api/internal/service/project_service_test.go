package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/BennyEisner/test-results/internal/models"
)

type mockProjectRepo struct {
	GetByIDFunc   func(ctx context.Context, id int64) (*models.Project, error)
	GetAllFunc    func(ctx context.Context) ([]*models.Project, error)
	GetByNameFunc func(ctx context.Context, name string) (*models.Project, error)
	CreateFunc    func(ctx context.Context, p *models.Project) error
	UpdateFunc    func(ctx context.Context, id int64, name string) (*models.Project, error)
	DeleteFunc    func(ctx context.Context, id int64) error
	CountFunc     func(ctx context.Context) (int, error)
}

func (m *mockProjectRepo) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *mockProjectRepo) GetAll(ctx context.Context) ([]*models.Project, error) {
	return m.GetAllFunc(ctx)
}
func (m *mockProjectRepo) GetByName(ctx context.Context, name string) (*models.Project, error) {
	return m.GetByNameFunc(ctx, name)
}
func (m *mockProjectRepo) Create(ctx context.Context, p *models.Project) error {
	return m.CreateFunc(ctx, p)
}
func (m *mockProjectRepo) Update(ctx context.Context, id int64, name string) (*models.Project, error) {
	return m.UpdateFunc(ctx, id, name)
}
func (m *mockProjectRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
func (m *mockProjectRepo) Count(ctx context.Context) (int, error) {
	return m.CountFunc(ctx)
}

func TestProjectService_GetProjectByID(t *testing.T) {
	repo := &mockProjectRepo{
		GetByIDFunc: func(ctx context.Context, id int64) (*models.Project, error) {
			if id == 1 {
				return &models.Project{ID: 1, Name: "Test"}, nil
			}
			return nil, sql.ErrNoRows
		},
	}
	svc := NewProjectService(repo)

	t.Run("found", func(t *testing.T) {
		p, err := svc.GetProjectByID(1)
		if err != nil || p.ID != 1 {
			t.Errorf("expected project, got %v, err %v", p, err)
		}
	})
	t.Run("not found", func(t *testing.T) {
		p, err := svc.GetProjectByID(2)
		if err == nil || p != nil {
			t.Errorf("expected error, got %v, err %v", p, err)
		}
	})
}

func TestProjectService_GetAllProjects(t *testing.T) {
	repo := &mockProjectRepo{
		GetAllFunc: func(ctx context.Context) ([]*models.Project, error) {
			return []*models.Project{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, nil
		},
	}
	svc := NewProjectService(repo)

	ps, err := svc.GetAllProjects()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ps) != 2 || ps[0].Name != "A" || ps[1].Name != "B" {
		t.Errorf("unexpected projects: %+v", ps)
	}
}

func TestProjectService_CreateProject(t *testing.T) {
	repo := &mockProjectRepo{
		CreateFunc: func(ctx context.Context, p *models.Project) error {
			if p.Name == "fail" {
				return errors.New("fail")
			}
			p.ID = 42
			return nil
		},
	}
	svc := NewProjectService(repo)

	t.Run("success", func(t *testing.T) {
		p, err := svc.CreateProject("ok")
		if err != nil || p.ID != 42 || p.Name != "ok" {
			t.Errorf("unexpected result: %v, err %v", p, err)
		}
	})
	t.Run("fail", func(t *testing.T) {
		p, err := svc.CreateProject("fail")
		if err == nil || p != nil {
			t.Errorf("expected error, got %v, err %v", p, err)
		}
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	repo := &mockProjectRepo{
		UpdateFunc: func(ctx context.Context, id int64, name string) (*models.Project, error) {
			if id == 1 {
				return &models.Project{ID: 1, Name: name}, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewProjectService(repo)

	p, err := svc.UpdateProject(1, "new")
	if err != nil || p.Name != "new" {
		t.Errorf("unexpected result: %v, err %v", p, err)
	}
	p, err = svc.UpdateProject(2, "fail")
	if err == nil || p != nil {
		t.Errorf("expected error, got %v, err %v", p, err)
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	repo := &mockProjectRepo{
		DeleteFunc: func(ctx context.Context, id int64) error {
			if id == 1 {
				return nil
			}
			return errors.New("not found")
		},
	}
	svc := NewProjectService(repo)

	t.Run("success", func(t *testing.T) {
		n, err := svc.DeleteProject(1)
		if err != nil || n != 1 {
			t.Errorf("unexpected result: %d, err %v", n, err)
		}
	})
	t.Run("fail", func(t *testing.T) {
		n, err := svc.DeleteProject(2)
		if err == nil || n != 0 {
			t.Errorf("expected error, got %d, err %v", n, err)
		}
	})
}

func TestProjectService_GetDBTestProjectCount(t *testing.T) {
	repo := &mockProjectRepo{
		CountFunc: func(ctx context.Context) (int, error) {
			return 7, nil
		},
	}
	svc := NewProjectService(repo)
	count, err := svc.GetDBTestProjectCount()
	if err != nil || count != 7 {
		t.Errorf("unexpected result: %d, err %v", count, err)
	}
}

func TestProjectService_GetProjectByName(t *testing.T) {
	repo := &mockProjectRepo{
		GetByNameFunc: func(ctx context.Context, name string) (*models.Project, error) {
			if name == "foo" {
				return &models.Project{ID: 1, Name: "foo"}, nil
			}
			return nil, sql.ErrNoRows
		},
	}
	svc := NewProjectService(repo)

	p, err := svc.GetProjectByName("foo")
	if err != nil || p.Name != "foo" {
		t.Errorf("unexpected result: %v, err %v", p, err)
	}
	p, err = svc.GetProjectByName("bar")
	if err == nil || p != nil {
		t.Errorf("expected error, got %v, err %v", p, err)
	}
}
