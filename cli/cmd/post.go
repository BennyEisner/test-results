package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	project string
	tags    []string
	file    string
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

		fmt.Printf("Posting results:\n")
		fmt.Printf("- Project: %s\n", project)
		fmt.Printf("- File: %s\n", file)
		fmt.Printf("- Type: %s\n", testType)
		fmt.Printf("- Tags: %s\n", strings.Join(tags, ", "))

		// TODO: Load and transform results file
		// TODO: POST to your API

		return nil
	},
}
