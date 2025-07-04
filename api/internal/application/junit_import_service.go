package application

import (
	"context"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// JUnitImportService implements the JUnitImportService interface
type JUnitImportService struct {
	buildService          ports.BuildService
	testSuiteService      ports.TestSuiteService
	testCaseService       ports.TestCaseService
	buildExecutionService ports.BuildExecutionService
}

func NewJUnitImportService(
	buildService ports.BuildService,
	testSuiteService ports.TestSuiteService,
	testCaseService ports.TestCaseService,
	buildExecutionService ports.BuildExecutionService,
) ports.JUnitImportService {
	return &JUnitImportService{
		buildService:          buildService,
		testSuiteService:      testSuiteService,
		testCaseService:       testCaseService,
		buildExecutionService: buildExecutionService,
	}
}

func (s *JUnitImportService) ProcessJUnitData(ctx context.Context, projectID int64, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, error) {
	if junitData == nil {
		return nil, fmt.Errorf("JUnit data cannot be nil")
	}

	// Create a new build for this import
	build := &models.Build{
		TestSuiteID: suiteID,
		ProjectID:   projectID,
		BuildNumber: fmt.Sprintf("import-%d", time.Now().Unix()),
		CIProvider:  "junit-import",
	}

	createdBuild, err := s.buildService.CreateBuild(ctx, build)
	if err != nil {
		return nil, fmt.Errorf("failed to create build: %w", err)
	}

	// Process each test suite in the JUnit data
	for _, testSuite := range junitData.TestSuites {
		if err := s.processTestSuite(ctx, createdBuild.ID, testSuite); err != nil {
			return nil, fmt.Errorf("failed to process test suite %s: %w", testSuite.Name, err)
		}
	}

	return createdBuild, nil
}

func (s *JUnitImportService) processTestSuite(ctx context.Context, buildID int64, testSuite models.JUnitTestSuite) error {
	// Create test cases and executions for this suite
	for _, testCase := range testSuite.TestCases {
		if err := s.processTestCase(ctx, buildID, testCase); err != nil {
			return fmt.Errorf("failed to process test case %s: %w", testCase.Name, err)
		}
	}
	return nil
}

func (s *JUnitImportService) processTestCase(ctx context.Context, buildID int64, testCase models.JUnitTestCase) error {
	// Create or find the test case
	createdTestCase, err := s.testCaseService.CreateTestCase(ctx, buildID, testCase.Name, testCase.Classname)
	if err != nil {
		return fmt.Errorf("failed to create test case: %w", err)
	}

	// Determine status based on JUnit data
	status := "passed"
	if testCase.Failure != nil {
		status = "failed"
	} else if testCase.Error != nil {
		status = "error"
	} else if testCase.Skipped != nil {
		status = "skipped"
	}

	// Create build execution
	execution := &models.BuildExecution{
		BuildID:       buildID,
		TestCaseID:    createdTestCase.ID,
		Status:        status,
		ExecutionTime: testCase.Time,
	}

	if err := s.buildExecutionService.CreateBuildExecutions(ctx, buildID, []*models.BuildExecution{execution}); err != nil {
		return fmt.Errorf("failed to create build execution: %w", err)
	}

	return nil
}
