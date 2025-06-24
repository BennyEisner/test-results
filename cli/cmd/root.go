package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "test-results",
	Short: "A CLI tool to replay and compare CDN results (Akamai vs Fastly)",
	Long: `Test Results CLI replays requests and compares the outputs 
between Akamai and Fastly, collecting statistics and summarizing differences.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `test-results --help` to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
