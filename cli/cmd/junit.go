package cmd

import (
	"fmt"
	"log"

	"github.com/BennyEisner/test-results/cli/internal/client"
	"github.com/BennyEisner/test-results/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	projectID int64
	suiteID   int64
	filePath  string
)

var junitCmd = &cobra.Command{
	Use:   "junit",
	Short: "Upload a JUnit XML file to test-results API",
	Long:  `Upload a JUnit XML file for a specific project and test-suite. A new build will be created in the database with the test-case results from the JUnit file`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg := config.LoadConfig()

		// Create API client
		apiClient := client.NewAPIClient(cfg)

		// Call the client to upload the file
		response, err := apiClient.PostJUnitFile(projectID, suiteID, filePath)
		if err != nil {
			log.Fatalf("Error uploading JUnit file: %v", err)
		}

		fmt.Println("Successfully uploaded JUnit file.")
		fmt.Println("API Response:", response)
	},
}

func init() {
	rootCmd.AddCommand(junitCmd)

	// Define flags for the junit command
	junitCmd.Flags().Int64Var(&projectID, "project-id", 0, "The ID of the project")
	junitCmd.Flags().Int64Var(&suiteID, "suite-id", 0, "The ID of the test suite")
	junitCmd.Flags().StringVar(&filePath, "file", "", "The path to the JUnit XML file")

	// Mark flags as required
	junitCmd.MarkFlagRequired("project-id")
	junitCmd.MarkFlagRequired("suite-id")
	junitCmd.MarkFlagRequired("file")
}
