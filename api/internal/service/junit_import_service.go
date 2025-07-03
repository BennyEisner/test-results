package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/models"
)

// JUnitImportServiceInterface defines the operations for importing JUnit data.
type JUnitImportServiceInterface interface {
	// ProcessJUnitData will create a new Build for this import, associated with projectID and suiteID.
	// It returns the created Build, a list of any non-fatal processing errors, and a fatal error if one occurs.
	ProcessJUnitData(projectID int64, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, []string, error)
}

// JUnitImportService provides services for importing JUnit XML data.
type JUnitImportService struct {
	DB                    *sql.DB
	BuildService          BuildServiceInterface
	TestSuiteService      TestSuiteServiceInterface
	TestCaseService       TestCaseServiceInterface
	BuildExecutionService BuildExecutionServiceInterface
}

// NewJUnitImportService creates a new JUnitImportService.
func NewJUnitImportService(
	db *sql.DB,
	bs BuildServiceInterface,
	tss TestSuiteServiceInterface,
	tcs TestCaseServiceInterface,
	bes BuildExecutionServiceInterface,
) *JUnitImportService {
	return &JUnitImportService{
		DB:                    db,
		BuildService:          bs,
		TestSuiteService:      tss,
		TestCaseService:       tcs,
		BuildExecutionService: bes,
	}
}

// ProcessJUnitData will contain the core logic for parsing and saving JUnit data
// It creates a new Build for this import, associated with the given projectID and suiteID
func (s *JUnitImportService) ProcessJUnitData(projectID int64, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, []string, error) {
	return s.processJUnitDataWithTx(projectID, suiteID, junitData)
}

// processJUnitDataWithTx manages the transaction and delegates to helpers.
func (s *JUnitImportService) processJUnitDataWithTx(projectID int64, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, []string, error) {
	var processingErrors []string
	var createdBuild *models.Build

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, processingErrors, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but re-panic with original panic
				_ = rbErr // Suppress unused variable warning
			}
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but return the original error
				_ = rbErr // Suppress unused variable warning
			}
		} else {
			err = tx.Commit()
			if err != nil {
				processingErrors = append(processingErrors, "Failed to commit transaction: "+err.Error())
			}
		}
	}()

	if err := s.validateSuiteExists(tx, projectID, suiteID); err != nil {
		return nil, processingErrors, err
	}

	createdBuild, err = s.createBuildWithTx(tx, suiteID, junitData)
	if err != nil {
		return nil, processingErrors, err
	}

	processingErrors = s.processTestSuitesWithTx(tx, createdBuild, suiteID, junitData, processingErrors)

	return createdBuild, processingErrors, err
}

// validateSuiteExists checks that the suite exists for the project.
func (s *JUnitImportService) validateSuiteExists(tx *sql.Tx, projectID, suiteID int64) error {
	_, err := s.TestSuiteService.GetProjectTestSuiteByIDWithTx(tx, projectID, suiteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("test suite ID %d not found for project ID %d", suiteID, projectID)
		}
		return fmt.Errorf("failed to validate suite ID %d for project %d: %w", suiteID, projectID, err)
	}
	return nil
}

// createBuildWithTx creates a new build for the suite.
func (s *JUnitImportService) createBuildWithTx(tx *sql.Tx, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, error) {
	buildName := "JUnit Import"
	if junitData.Name != "" {
		buildName = junitData.Name
	} else if len(junitData.TestSuites) == 1 && junitData.TestSuites[0].Name != "" {
		buildName = junitData.TestSuites[0].Name
	}
	buildTimestamp := time.Now()
	var totalTestCaseCount int64
	for _, junitSuite := range junitData.TestSuites {
		totalTestCaseCount += int64(len(junitSuite.TestCases))
	}
	buildToCreate := &models.Build{
		TestSuiteID:   suiteID,
		BuildNumber:   buildName,
		CIProvider:    "JUnit Import",
		CIURL:         nil,
		CreatedAt:     buildTimestamp,
		TestCaseCount: totalTestCaseCount,
	}
	return s.BuildService.CreateBuildWithTx(tx, buildToCreate)
}

// processTestSuitesWithTx processes all test suites and their test cases.
func (s *JUnitImportService) processTestSuitesWithTx(tx *sql.Tx, createdBuild *models.Build, suiteID int64, junitData *models.JUnitTestSuites, processingErrors []string) []string {
	for _, junitSuite := range junitData.TestSuites {
		executionInputs, errs := s.createExecutionInputsForSuite(tx, suiteID, junitSuite)
		processingErrors = append(processingErrors, errs...)
		if len(executionInputs) > 0 {
			_, batchErrors, execErr := s.BuildExecutionService.CreateBuildExecutionsWithTx(tx, createdBuild.ID, executionInputs)
			if execErr != nil {
				errMsg := fmt.Sprintf("Fatal error during batch creation of executions for build %d, suite '%s': %v", createdBuild.ID, junitSuite.Name, execErr)
				processingErrors = append(processingErrors, errMsg)
			}
			for _, batchErr := range batchErrors {
				processingErrors = append(processingErrors, fmt.Sprintf("Error creating execution for build %d, suite '%s': %s", createdBuild.ID, junitSuite.Name, batchErr))
			}
		}
	}
	return processingErrors
}

// createExecutionInputsForSuite creates execution inputs for all test cases in a suite.
func (s *JUnitImportService) createExecutionInputsForSuite(tx *sql.Tx, suiteID int64, junitSuite models.JUnitTestSuite) ([]models.BuildExecutionInput, []string) {
	var executionInputs []models.BuildExecutionInput
	var processingErrors []string
	for _, junitCase := range junitSuite.TestCases {
		testCase, tcErr := s.TestCaseService.FindOrCreateTestCaseWithTx(tx, suiteID, junitCase.Name, junitCase.Classname)
		if tcErr != nil {
			errMsg := fmt.Sprintf("Error finding/creating test case '%s' (class: '%s'): %v", junitCase.Name, junitCase.Classname, tcErr)
			processingErrors = append(processingErrors, errMsg)
			continue
		}
		status, failureMessage, failureType, failureDetails := getJUnitCaseStatus(junitCase)
		currentExecutionInput := models.BuildExecutionInput{
			TestCaseID:     testCase.ID,
			Status:         status,
			ExecutionTime:  junitCase.Time,
			FailureMessage: failureMessage,
			FailureType:    failureType,
			FailureDetails: failureDetails,
		}
		executionInputs = append(executionInputs, currentExecutionInput)
	}
	return executionInputs, processingErrors
}

// getJUnitCaseStatus determines the status and failure details for a JUnit test case.
func getJUnitCaseStatus(junitCase models.JUnitTestCase) (status string, failureMessage, failureType, failureDetails *string) {
	status = "passed"
	if junitCase.Failure != nil {
		status = "failed"
		failureMessage = &junitCase.Failure.Message
		failureType = &junitCase.Failure.Type
		failureDetails = &junitCase.Failure.Value
	} else if junitCase.Error != nil {
		status = "error"
		failureMessage = &junitCase.Error.Message
		failureType = &junitCase.Error.Type
		failureDetails = &junitCase.Error.Value
	} else if junitCase.Skipped != nil {
		status = "skipped"
		if junitCase.Skipped.Message != "" {
			failureMessage = &junitCase.Skipped.Message
		}
	}
	return
}
