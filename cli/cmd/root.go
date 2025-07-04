package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "test-results",
	Short: "A CLI tool to post test results to the test-results API",
	Long: `Test Results CLI searches for and posts tests results to 
a RESTful API to be consumed by a custom dashboard for quality of life
improvements for QA and stakeholders.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `results --help` to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
