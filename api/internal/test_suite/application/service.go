package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/project/domain"
	"github.com/BennyEisner/test-results/internal/test_suite/domain/models"
	"github.com/BennyEisner/test-results/internal/test_suite/domain/ports"
)

// TestSuiteService implements the TestSuiteService interface
type TestSuiteService struct {
	repo ports.TestSuiteRepository
}

func NewTestSuiteService(repo ports.TestSuiteRepository) ports.TestSuiteService {
	return &TestSuiteService{repo: repo}
}

func (s *TestSuiteService) GetTestSuite(ctx context.Context, id int64) (*models.TestSuite, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidTestSuiteName
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

func (s *TestSuiteService) GetTestSuitesByProject(ctx context.Context, projectID int64) ([]*models.TestSuite, error) {
	if projectID <= 0 {
		return nil, domain.ErrInvalidProjectName
	}
	ts, err := s.repo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test suites by project ID %d: %w", projectID, err)
	}
	return ts, nil
}

func (s *TestSuiteService) GetTestSuiteByName(ctx context.Context, projectID int64, name string) (*models.TestSuite, error) {
	if projectID <= 0 || name == "" {
		return nil, domain.ErrInvalidTestSuiteName
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

func (s *TestSuiteService) CreateTestSuite(ctx context.Context, projectID int64, name string, parentID *int64, time float64) (*models.TestSuite, error) {
	if projectID <= 0 || name == "" {
		return nil, domain.ErrInvalidTestSuiteName
	}
	// Check for duplicate
	existing, err := s.repo.GetByName(ctx, projectID, name)
	if err == nil && existing != nil {
		return nil, domain.ErrDuplicateTestSuite
	}
	ts := &models.TestSuite{
		ProjectID: projectID,
		Name:      name,
		ParentID:  parentID,
		Time:      time,
	}
	if err := s.repo.Create(ctx, ts); err != nil {
		return nil, fmt.Errorf("failed to create test suite: %w", err)
	}
	return ts, nil
}

func (s *TestSuiteService) UpdateTestSuite(ctx context.Context, id int64, name string) (*models.TestSuite, error) {
	if id <= 0 || name == "" {
		return nil, domain.ErrInvalidTestSuiteName
	}
	ts, err := s.repo.Update(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to update test suite: %w", err)
	}
	return ts, nil
}

func (s *TestSuiteService) DeleteTestSuite(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidTestSuiteName
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete test suite: %w", err)
	}
	return nil
}
