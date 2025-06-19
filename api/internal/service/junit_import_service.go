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

// ProcessJUnitData will contain the core logic for parsing and saving JUnit data.
// It creates a new Build for this import, associated with the given projectID and suiteID.
func (s *JUnitImportService) ProcessJUnitData(projectID int64, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, []string, error) {
	var processingErrors []string
	var createdBuild *models.Build

	// 1. Start a database transaction.
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, processingErrors, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			tx.Rollback() // Rollback on error
		} else {
			err = tx.Commit() // Commit on success
			if err != nil {
				processingErrors = append(processingErrors, "Failed to commit transaction: "+err.Error())
			}
		}
	}()

	// 2. Validate that the provided suiteID exists for the projectID.
	_, err = s.TestSuiteService.GetProjectTestSuiteByIDWithTx(tx, projectID, suiteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, processingErrors, fmt.Errorf("test suite ID %d not found for project ID %d", suiteID, projectID)
		}
		return nil, processingErrors, fmt.Errorf("failed to validate suite ID %d for project %d: %w", suiteID, projectID, err)
	}

	// 3. Create a new models.Build using BuildService (associated with projectID and suiteID).
	buildName := "JUnit Import" // Default name
	if junitData.Name != "" {   // Use name from <testsuites name="X"> if present
		buildName = junitData.Name
	} else if len(junitData.TestSuites) == 1 && junitData.TestSuites[0].Name != "" { // Or from the single <testsuite name="Y">
		buildName = junitData.TestSuites[0].Name
	}

	// Timestamp for the build. Could also parse from junitData.TestSuites[0].Timestamp if available and reliable.
	buildTimestamp := time.Now()
	// Example: Parsing timestamp from JUnit XML if needed
	// if len(junitData.TestSuites) > 0 && junitData.TestSuites[0].Timestamp != "" {
	//	parsedTs, tsErr := time.Parse(time.RFC3339Nano, junitData.TestSuites[0].Timestamp) // Or other expected format
	//	if tsErr == nil {
	//		buildTimestamp = parsedTs
	//	} else {
	//		processingErrors = append(processingErrors, fmt.Sprintf("Warning: could not parse build timestamp '%s': %v", junitData.TestSuites[0].Timestamp, tsErr))
	//	}
	//}

	buildToCreate := &models.Build{
		TestSuiteID: suiteID,
		BuildNumber: buildName,      // Using the derived name as BuildNumber
		CIProvider:  "JUnit Import", // Or derive from junitData.TestSuites[0].Hostname
		CIURL:       nil,            // Can be set if available in XML
		CreatedAt:   buildTimestamp, // This will be overridden by DB's NOW() in CreateBuildWithTx, but good to have
	}

	createdBuild, err = s.BuildService.CreateBuildWithTx(tx, buildToCreate)
	if err != nil {
		return nil, processingErrors, fmt.Errorf("failed to create build: %w", err)
	}

	// 4. Iterate through junitData.TestSuites.
	for _, junitSuite := range junitData.TestSuites {
		// Optional: Validate junitSuite.Name against the name of the suite from suiteID if desired.
		// log.Printf("Processing TestSuite from XML: %s (DB Suite ID: %d, Build ID: %d)\n", junitSuite.Name, suiteID, createdBuild.ID)

		var executionInputs []models.BuildExecutionInput

		// 5. For each JUnitTestCase within the JUnitTestSuite(s):
		for _, junitCase := range junitSuite.TestCases {
			testCase, tcErr := s.TestCaseService.FindOrCreateTestCaseWithTx(tx, suiteID, junitCase.Name, junitCase.Classname)
			if tcErr != nil {
				errMsg := fmt.Sprintf("Error finding/creating test case '%s' (class: '%s'): %v", junitCase.Name, junitCase.Classname, tcErr)
				processingErrors = append(processingErrors, errMsg)
				continue // Skip this test case, proceed with others
			}

			status := "passed"
			var failureMessage, failureType, failureDetails *string

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
				if junitCase.Skipped.Message != "" { // Store skipped message if available
					// Assuming FailureMessage can be used for skipped messages, or add a specific field
					failureMessage = &junitCase.Skipped.Message
				}
			}

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

		// b. Create models.BuildTestCaseExecution records using BuildExecutionService.
		if len(executionInputs) > 0 {
			_, batchErrors, execErr := s.BuildExecutionService.CreateBuildExecutionsWithTx(tx, createdBuild.ID, executionInputs)
			if execErr != nil {
				// This is a more fatal error for this batch of executions
				errMsg := fmt.Sprintf("Fatal error during batch creation of executions for build %d, suite '%s': %v", createdBuild.ID, junitSuite.Name, execErr)
				processingErrors = append(processingErrors, errMsg)
				// Depending on severity, you might choose to return `createdBuild, processingErrors, execErr` here to force rollback.
				// For now, we collect the error and let the transaction attempt to commit with partial data if other parts succeeded.
			}
			if len(batchErrors) > 0 {
				for _, batchErr := range batchErrors {
					processingErrors = append(processingErrors, fmt.Sprintf("Error creating execution for build %d, suite '%s': %s", createdBuild.ID, junitSuite.Name, batchErr))
				}
			}
		}
	}
	if len(processingErrors) > 0 {
		// Even if there are processing errors, the transaction might still commit if `err` is nil here.
		// The `err` variable in the defer function's scope is key.
		// If we want to force a rollback on processingErrors, we should set `err` here.
		// For now, let's assume processingErrors are non-fatal for the transaction itself unless a DB operation failed.
	}

	return createdBuild, processingErrors, err // err will be nil if commit was successful, or the commit error
}
