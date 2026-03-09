package diff

import (
	"strings"

	"github.com/loulou/go-xldiff/internal/reader"
)

// Row holds the cell values for a single spreadsheet row.
type Row struct {
	Values []string
}

// CellChange records a single cell-level difference between two rows.
type CellChange struct {
	Column   string
	OldValue string
	NewValue string
}

// RowDiff records changes to a row identified by its key columns.
type RowDiff struct {
	Key     []string
	Changes []CellChange
}

// DiffResult holds the complete result of comparing two tables.
type DiffResult struct {
	Headers  []string
	Added    []Row
	Removed  []Row
	Modified []RowDiff
}

// HasDifferences returns true when any adds, removes, or modifications exist.
func (r *DiffResult) HasDifferences() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

// Compare diffs two tables. When keys is nil or empty, positional (all-column)
// mode is used. When keys are provided, rows are matched by key columns.
func Compare(a, b *reader.Table, keys []string) *DiffResult {
	result := &DiffResult{Headers: a.Headers}
	if len(keys) == 0 {
		comparePositional(a, b, result)
	} else {
		compareByKey(a, b, keys, result)
	}
	return result
}

func comparePositional(a, b *reader.Table, result *DiffResult) {
	countA := make(map[string]int)
	countB := make(map[string]int)

	for _, row := range a.Rows {
		countA[rowKey(row)]++
	}
	for _, row := range b.Rows {
		countB[rowKey(row)]++
	}

	// Find removed rows: rows in A whose count exceeds what B has.
	remaining := make(map[string]int)
	for k, v := range countB {
		remaining[k] = v
	}
	for _, row := range a.Rows {
		k := rowKey(row)
		if remaining[k] > 0 {
			remaining[k]--
		} else {
			result.Removed = append(result.Removed, Row{Values: row})
		}
	}

	// Find added rows: rows in B whose count exceeds what A has.
	remaining2 := make(map[string]int)
	for k, v := range countA {
		remaining2[k] = v
	}
	for _, row := range b.Rows {
		k := rowKey(row)
		if remaining2[k] > 0 {
			remaining2[k]--
		} else {
			result.Added = append(result.Added, Row{Values: row})
		}
	}
}

func compareByKey(a, b *reader.Table, keys []string, result *DiffResult) {
	// Placeholder — implemented in Task 5.
}

func rowKey(row []string) string {
	return strings.Join(row, "\x00")
}
