package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/output"
	"github.com/loulou/go-xldiff/internal/reader"
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
	file1, sheet1 := reader.ParseFileArg(args[0])
	file2, sheet2 := reader.ParseFileArg(args[1])

	skip1, skip2, err := reader.ParseSkipFlag(skipFlag)
	if err != nil {
		return fmt.Errorf("invalid --skip: %w", err)
	}

	table1, err := reader.ReadSheet(file1, sheet1, skip1, noHeaderFlag)
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	table2, err := reader.ReadSheet(file2, sheet2, skip2, noHeaderFlag)
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	var keys []string
	if keyFlag != "" {
		keys = strings.Split(keyFlag, ",")
	}

	result := diff.Compare(table1, table2, keys)

	switch formatFlag {
	case "json":
		if err := output.WriteJSON(os.Stdout, result); err != nil {
			return err
		}
	case "csv":
		if err := output.WriteCSV(os.Stdout, result); err != nil {
			return err
		}
	default:
		output.WriteText(os.Stdout, result, table1, table2, noColorFlag)
	}

	if result.HasDifferences() {
		os.Exit(1)
	}
	return nil
}
