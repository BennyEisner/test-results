package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type TestCaseService struct {
	repo domain.TestCaseRepository
}

func NewTestCaseService(repo domain.TestCaseRepository) domain.TestCaseService {
	return &TestCaseService{repo: repo}
}

func (s *TestCaseService) GetTestCaseByID(ctx context.Context, id int64) (*domain.TestCase, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}
	testCase, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test case by ID %d: %w", id, err)
	}
	if testCase == nil {
		return nil, domain.ErrTestCaseNotFound
	}
	return testCase, nil
}

func (s *TestCaseService) GetTestCasesBySuiteID(ctx context.Context, suiteID int64) ([]*domain.TestCase, error) {
	if suiteID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	testCases, err := s.repo.GetAllBySuiteID(ctx, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases by suite ID %d: %w", suiteID, err)
	}
	return testCases, nil
}

func (s *TestCaseService) GetTestCaseByName(ctx context.Context, suiteID int64, name string) (*domain.TestCase, error) {
	if suiteID <= 0 || name == "" {
		return nil, domain.ErrInvalidInput
	}
	testCase, err := s.repo.GetByName(ctx, suiteID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get test case by name %s in suite %d: %w", name, suiteID, err)
	}
	if testCase == nil {
		return nil, domain.ErrTestCaseNotFound
	}
	return testCase, nil
}

func (s *TestCaseService) CreateTestCase(ctx context.Context, suiteID int64, name, classname string) (*domain.TestCase, error) {
	if suiteID <= 0 || name == "" || classname == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if test case already exists
	existingTestCase, err := s.repo.GetByName(ctx, suiteID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing test case: %w", err)
	}
	if existingTestCase != nil {
		return nil, domain.ErrDuplicateTestCase
	}

	testCase := &domain.TestCase{
		SuiteID:   suiteID,
		Name:      name,
		Classname: classname,
	}

	if err := s.repo.Create(ctx, testCase); err != nil {
		return nil, fmt.Errorf("failed to create test case: %w", err)
	}
	return testCase, nil
}

func (s *TestCaseService) UpdateTestCase(ctx context.Context, id int64, name, classname string) (*domain.TestCase, error) {
	if id <= 0 || name == "" || classname == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if test case exists
	existingTestCase, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing test case: %w", err)
	}
	if existingTestCase == nil {
		return nil, domain.ErrTestCaseNotFound
	}

	updatedTestCase, err := s.repo.Update(ctx, id, name, classname)
	if err != nil {
		return nil, fmt.Errorf("failed to update test case: %w", err)
	}
	return updatedTestCase, nil
}

func (s *TestCaseService) DeleteTestCase(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	// Check if test case exists
	existingTestCase, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing test case: %w", err)
	}
	if existingTestCase == nil {
		return domain.ErrTestCaseNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete test case: %w", err)
	}
	return nil
}
