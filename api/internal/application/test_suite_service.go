package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type TestSuiteService struct {
	repo domain.TestSuiteRepository
}

func NewTestSuiteService(repo domain.TestSuiteRepository) domain.TestSuiteService {
	return &TestSuiteService{repo: repo}
}

func (s *TestSuiteService) GetTestSuiteByID(ctx context.Context, id int64) (*domain.TestSuite, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}
	ts, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test suite by ID %d: %w", id, err)
	}
	if ts == nil {
		return nil, domain.ErrTestSuiteNotFound
	}
	return ts, nil
}

func (s *TestSuiteService) GetTestSuitesByProjectID(ctx context.Context, projectID int64) ([]*domain.TestSuite, error) {
	if projectID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	ts, err := s.repo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test suites by project ID %d: %w", projectID, err)
	}
	return ts, nil
}

func (s *TestSuiteService) GetTestSuiteByName(ctx context.Context, projectID int64, name string) (*domain.TestSuite, error) {
	if projectID <= 0 || name == "" {
		return nil, domain.ErrInvalidInput
	}
	ts, err := s.repo.GetByName(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get test suite by name: %w", err)
	}
	if ts == nil {
		return nil, domain.ErrTestSuiteNotFound
	}
	return ts, nil
}

func (s *TestSuiteService) CreateTestSuite(ctx context.Context, projectID int64, name string, parentID *int64) (*domain.TestSuite, error) {
	if projectID <= 0 || name == "" {
		return nil, domain.ErrInvalidInput
	}
	// Check for duplicate
	existing, err := s.repo.GetByName(ctx, projectID, name)
	if err == nil && existing != nil {
		return nil, domain.ErrDuplicateTestSuite
	}
	ts := &domain.TestSuite{
		ProjectID: projectID,
		Name:      name,
		ParentID:  parentID,
	}
	if err := s.repo.Create(ctx, ts); err != nil {
		return nil, fmt.Errorf("failed to create test suite: %w", err)
	}
	return ts, nil
}

func (s *TestSuiteService) UpdateTestSuite(ctx context.Context, id int64, name string) (*domain.TestSuite, error) {
	if id <= 0 || name == "" {
		return nil, domain.ErrInvalidInput
	}
	ts, err := s.repo.Update(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to update test suite: %w", err)
	}
	return ts, nil
}

func (s *TestSuiteService) DeleteTestSuite(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete test suite: %w", err)
	}
	return nil
}
