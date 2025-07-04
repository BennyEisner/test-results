package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// TestCaseService implements the TestCaseService interface
type TestCaseService struct {
	repo ports.TestCaseRepository
}

func NewTestCaseService(repo ports.TestCaseRepository) ports.TestCaseService {
	return &TestCaseService{repo: repo}
}

func (s *TestCaseService) GetTestCaseByID(ctx context.Context, id int64) (*models.TestCase, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidTestCaseName
	}

	tc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test case by ID %d: %w", id, err)
	}
	if tc == nil {
		return nil, domain.ErrTestCaseNotFound
	}
	return tc, nil
}

func (s *TestCaseService) GetTestCasesBySuiteID(ctx context.Context, suiteID int64) ([]*models.TestCase, error) {
	if suiteID <= 0 {
		return nil, domain.ErrInvalidTestSuiteName
	}
	tcs, err := s.repo.GetAllBySuiteID(ctx, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases by suite ID %d: %w", suiteID, err)
	}
	return tcs, nil
}

func (s *TestCaseService) GetTestCaseByName(ctx context.Context, suiteID int64, name string) (*models.TestCase, error) {
	if suiteID <= 0 || name == "" {
		return nil, domain.ErrInvalidTestCaseName
	}

	tc, err := s.repo.GetByName(ctx, suiteID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get test case by name: %w", err)
	}
	if tc == nil {
		return nil, domain.ErrTestCaseNotFound
	}
	return tc, nil
}

func (s *TestCaseService) CreateTestCase(ctx context.Context, suiteID int64, name, classname string) (*models.TestCase, error) {
	if suiteID <= 0 || name == "" || classname == "" {
		return nil, domain.ErrInvalidTestCaseName
	}

	tc := &models.TestCase{
		SuiteID:   suiteID,
		Name:      name,
		Classname: classname,
	}

	if err := s.repo.Create(ctx, tc); err != nil {
		return nil, fmt.Errorf("failed to create test case: %w", err)
	}
	return tc, nil
}

func (s *TestCaseService) UpdateTestCase(ctx context.Context, id int64, name, classname string) (*models.TestCase, error) {
	if id <= 0 || name == "" || classname == "" {
		return nil, domain.ErrInvalidTestCaseName
	}

	tc, err := s.repo.Update(ctx, id, name, classname)
	if err != nil {
		return nil, fmt.Errorf("failed to update test case: %w", err)
	}
	return tc, nil
}

func (s *TestCaseService) DeleteTestCase(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidTestCaseName
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete test case: %w", err)
	}
	return nil
}
