package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/louloulibs/go-xldiff/internal/diff"
	"github.com/louloulibs/go-xldiff/internal/reader"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func WriteText(w io.Writer, result *diff.DiffResult, a, b *reader.Table, noColor bool) {
	red, green, yellow, reset := colorRed, colorGreen, colorYellow, colorReset
	if noColor {
		red, green, yellow, reset = "", "", "", ""
	}

	fmt.Fprintf(w, "--- %s (%s)\n", a.Sheet, a.FileName)
	fmt.Fprintf(w, "+++ %s (%s)\n", b.Sheet, b.FileName)
	fmt.Fprintln(w)

	if !result.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	fmt.Fprintf(w, "Added rows: %d\n", len(result.Added))
	fmt.Fprintf(w, "Removed rows: %d\n", len(result.Removed))
	fmt.Fprintf(w, "Modified rows: %d\n", len(result.Modified))
	fmt.Fprintln(w)

	for _, row := range result.Removed {
		fmt.Fprintf(w, "%s- %s%s\n", red, formatRow(result.Headers, row.Values), reset)
	}
	for _, row := range result.Added {
		fmt.Fprintf(w, "%s+ %s%s\n", green, formatRow(result.Headers, row.Values), reset)
	}
	if len(result.Modified) > 0 && (len(result.Added) > 0 || len(result.Removed) > 0) {
		fmt.Fprintln(w)
	}
	for _, mod := range result.Modified {
		parts := make([]string, len(mod.Changes))
		for i, c := range mod.Changes {
			parts[i] = fmt.Sprintf("%s: %q → %q", c.Column, c.OldValue, c.NewValue)
		}
		keyStr := strings.Join(mod.Key, ", ")
		fmt.Fprintf(w, "%s~ [%s]  %s%s\n", yellow, keyStr, strings.Join(parts, ", "), reset)
	}
}

func formatRow(headers, values []string) string {
	parts := make([]string, len(headers))
	for i, h := range headers {
		v := ""
		if i < len(values) {
			v = values[i]
		}
		parts[i] = fmt.Sprintf("%s: %q", h, v)
	}
	return strings.Join(parts, "  ")
}
