package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	keyFlag      string
	skipFlag     string
	formatFlag   string
	noHeaderFlag bool
	noColorFlag  bool
)

var rootCmd = &cobra.Command{
	Use:   "go-xldiff <file1>[:<sheet>] <file2>[:<sheet>]",
	Short: "Diff two Excel sheets",
	Long:  "Compare two Excel sheets and report added, removed, and modified rows.",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.Flags().StringVar(&keyFlag, "key", "", "Row identity columns (header name or 0-based index), comma-separated")
	rootCmd.Flags().StringVar(&skipFlag, "skip", "0", "Rows to skip before header. Single value or comma-separated pair (e.g. 3,5)")
	rootCmd.Flags().StringVar(&formatFlag, "format", "text", "Output format: text, json, csv")
	rootCmd.Flags().BoolVar(&noHeaderFlag, "no-header", false, "Treat first row as data, not headers")
	rootCmd.Flags().BoolVar(&noColorFlag, "no-color", false, "Disable colored terminal output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDiff(cmd *cobra.Command, args []string) error {
	fmt.Println("Not implemented yet")
	return nil
}
