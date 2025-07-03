package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/cli/internal/client"
	"github.com/BennyEisner/test-results/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	project  string
	tags     []string
	file     string
	testType string
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post test results to the REST API",
	Long: `Upload test results (JUnit or ReadyAPI format) to a centralized results API.

Example:
  test-results post --project myproj --file results.xml --type junit --tags smoke,api`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if project == "" {
			return fmt.Errorf("required flag --project not set")
		}
		if file == "" {
			return fmt.Errorf("required flag --file not set")
		}
		if testType != "junit" && testType != "readyapi" {
			return fmt.Errorf("unsupported --type: %s (must be 'junit' or 'readyapi')", testType)
		}

		// For now, only implement JUnit posting
		if testType != "junit" {
			return fmt.Errorf("only 'junit' type is currently supported")
		}

		// Parse project and suite IDs from the project flag
		// Expected format: "projectID:suiteID"
		parts := strings.Split(project, ":")
		if len(parts) != 2 {
			return fmt.Errorf("project flag must be in format 'projectID:suiteID' (e.g., '1:2')")
		}

		projectID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid project ID: %s", parts[0])
		}

		suiteID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid suite ID: %s", parts[1])
		}

		fmt.Printf("Posting JUnit results:\n")
		fmt.Printf("- Project ID: %d\n", projectID)
		fmt.Printf("- Suite ID: %d\n", suiteID)
		fmt.Printf("- File: %s\n", file)
		fmt.Printf("- Type: %s\n", testType)
		if len(tags) > 0 {
			fmt.Printf("- Tags: %s\n", strings.Join(tags, ", "))
		}

		// Load configuration
		cfg := config.LoadConfig()

		// Create API client
		apiClient := client.NewAPIClient(cfg)

		// Call the client to upload the file
		response, err := apiClient.PostJUnitFile(projectID, suiteID, file)
		if err != nil {
			log.Fatalf("Error uploading JUnit file: %v", err)
		}

		fmt.Println("Successfully uploaded JUnit file.")
		fmt.Println("API Response:", response)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.Flags().StringVar(&project, "project", "", "Project ID (required)")
	postCmd.Flags().StringVar(&file, "file", "junit.xml", "Path to JUnit XML file (optional)")
	postCmd.Flags().StringVar(&testType, "type", "junit", "Test type: junit or readyapi (optional)")
	postCmd.Flags().StringSliceVar(&tags, "tags", nil, "Comma-separated tags (optional)")
	postCmd.MarkFlagRequired("project")
}
