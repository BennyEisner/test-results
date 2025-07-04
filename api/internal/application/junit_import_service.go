package application

import (
	"context"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/domain"
)

type JUnitImportService struct {
	buildService          domain.BuildService
	testSuiteService      domain.TestSuiteService
	testCaseService       domain.TestCaseService
	buildExecutionService domain.BuildExecutionService
}

func NewJUnitImportService(
	buildService domain.BuildService,
	testSuiteService domain.TestSuiteService,
	testCaseService domain.TestCaseService,
	buildExecutionService domain.BuildExecutionService,
) domain.JUnitImportService {
	return &JUnitImportService{
		buildService:          buildService,
		testSuiteService:      testSuiteService,
		testCaseService:       testCaseService,
		buildExecutionService: buildExecutionService,
	}
}

func (s *JUnitImportService) ProcessJUnitData(ctx context.Context, projectID int64, suiteID int64, junitData *domain.JUnitTestSuites) (*domain.Build, error) {
	// Validate suite exists
	if err := s.validateTestSuite(ctx, suiteID); err != nil {
		return nil, err
	}

	// Create Build
	build, err := s.createBuild(ctx, projectID, suiteID, junitData)
	if err != nil {
		return nil, err
	}

	// Process test suites and test cases
	if err := s.processTestCases(ctx, build.ID, suiteID, junitData); err != nil {
		return nil, err
	}

	return build, nil
}

func (s *JUnitImportService) validateTestSuite(ctx context.Context, suiteID int64) error {
	_, err := s.testSuiteService.GetTestSuiteByID(ctx, suiteID)
	if err != nil {
		return fmt.Errorf("test suite not found: %w", err)
	}
	return nil
}

func (s *JUnitImportService) createBuild(ctx context.Context, projectID int64, suiteID int64, junitData *domain.JUnitTestSuites) (*domain.Build, error) {
	buildName := s.determineBuildName(junitData)
	totalTestCaseCount := s.calculateTotalTestCaseCount(junitData)

	buildToCreate := &domain.Build{
		TestSuiteID:   suiteID,
		ProjectID:     projectID,
		BuildNumber:   buildName,
		CIProvider:    "JUnit Import",
		CIURL:         nil,
		CreatedAt:     time.Now(),
		TestCaseCount: totalTestCaseCount,
	}

	build, err := s.buildService.CreateBuild(ctx, buildToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create build: %w", err)
	}

	return build, nil
}

func (s *JUnitImportService) determineBuildName(junitData *domain.JUnitTestSuites) string {
	if junitData.Name != "" {
		return junitData.Name
	}
	if len(junitData.TestSuites) == 1 && junitData.TestSuites[0].Name != "" {
		return junitData.TestSuites[0].Name
	}
	return "JUnit Import"
}

func (s *JUnitImportService) calculateTotalTestCaseCount(junitData *domain.JUnitTestSuites) int64 {
	var totalTestCaseCount int64
	for _, junitSuite := range junitData.TestSuites {
		totalTestCaseCount += int64(len(junitSuite.TestCases))
	}
	return totalTestCaseCount
}

func (s *JUnitImportService) processTestCases(ctx context.Context, buildID int64, suiteID int64, junitData *domain.JUnitTestSuites) error {
	for _, junitSuite := range junitData.TestSuites {
		for _, junitCase := range junitSuite.TestCases {
			if err := s.processTestCase(ctx, buildID, suiteID, junitCase); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *JUnitImportService) processTestCase(ctx context.Context, buildID int64, suiteID int64, junitCase domain.JUnitTestCase) error {
	// Find or create test case
	testCase, err := s.getOrCreateTestCase(ctx, suiteID, junitCase)
	if err != nil {
		return err
	}

	// Create build execution
	execution := s.createBuildExecution(buildID, testCase.ID, junitCase)
	err = s.buildExecutionService.CreateBuildExecutions(ctx, buildID, []*domain.BuildExecution{execution})
	if err != nil {
		return fmt.Errorf("failed to create build execution: %w", err)
	}

	return nil
}

func (s *JUnitImportService) getOrCreateTestCase(ctx context.Context, suiteID int64, junitCase domain.JUnitTestCase) (*domain.TestCase, error) {
	testCase, err := s.testCaseService.GetTestCaseByName(ctx, suiteID, junitCase.Name)
	if err == domain.ErrTestCaseNotFound {
		testCase, err = s.testCaseService.CreateTestCase(ctx, suiteID, junitCase.Name, junitCase.Classname)
		if err != nil {
			return nil, fmt.Errorf("failed to create test case: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get test case: %w", err)
	}
	return testCase, nil
}

func (s *JUnitImportService) createBuildExecution(buildID int64, testCaseID int64, junitCase domain.JUnitTestCase) *domain.BuildExecution {
	status := s.determineTestCaseStatus(junitCase)

	return &domain.BuildExecution{
		BuildID:       buildID,
		TestCaseID:    testCaseID,
		Status:        status,
		ExecutionTime: junitCase.Time,
		CreatedAt:     time.Now(),
	}
}

func (s *JUnitImportService) determineTestCaseStatus(junitCase domain.JUnitTestCase) string {
	if junitCase.Failure != nil {
		return domain.StatusFailed
	}
	if junitCase.Error != nil {
		return domain.StatusError
	}
	if junitCase.Skipped != nil {
		return domain.StatusSkipped
	}
	return domain.StatusPassed
}
