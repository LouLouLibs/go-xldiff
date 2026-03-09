package diff

import (
	"fmt"
	"strconv"
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
	Headers    []string
	KeyColumns []string
	Added      []Row
	Removed    []Row
	Modified   []RowDiff
}

// HasDifferences returns true when any adds, removes, or modifications exist.
func (r *DiffResult) HasDifferences() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

// CompareWithWarnings diffs two tables and returns any warnings (e.g. duplicate keys).
func CompareWithWarnings(a, b *reader.Table, keys []string) (*DiffResult, []string) {
	var warnings []string
	if len(keys) > 0 {
		keyIndices := resolveKeyIndices(a.Headers, keys)
		warnings = append(warnings, checkDuplicateKeys(a.Rows, keyIndices, a.FileName)...)
		warnings = append(warnings, checkDuplicateKeys(b.Rows, keyIndices, b.FileName)...)
	}
	result := Compare(a, b, keys)
	return result, warnings
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
	keyIndices := resolveKeyIndices(a.Headers, keys)

	// Store key column names for output formatters
	result.KeyColumns = make([]string, len(keyIndices))
	for i, idx := range keyIndices {
		if idx < len(a.Headers) {
			result.KeyColumns[i] = a.Headers[idx]
		}
	}

	mapA := buildKeyMap(a.Rows, keyIndices)
	mapB := buildKeyMap(b.Rows, keyIndices)

	for k, rowA := range mapA {
		rowB, exists := mapB[k]
		if !exists {
			result.Removed = append(result.Removed, Row{Values: rowA})
			continue
		}
		var changes []CellChange
		for i, header := range a.Headers {
			if isKeyIndex(i, keyIndices) {
				continue
			}
			if rowA[i] != rowB[i] {
				changes = append(changes, CellChange{
					Column:   header,
					OldValue: rowA[i],
					NewValue: rowB[i],
				})
			}
		}
		if len(changes) > 0 {
			keyVals := extractKey(rowA, keyIndices)
			result.Modified = append(result.Modified, RowDiff{Key: keyVals, Changes: changes})
		}
	}

	for k, rowB := range mapB {
		if _, exists := mapA[k]; !exists {
			result.Added = append(result.Added, Row{Values: rowB})
		}
	}
}

func resolveKeyIndices(headers []string, keys []string) []int {
	indices := make([]int, len(keys))
	for i, key := range keys {
		if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < len(headers) {
			indices[i] = idx
			continue
		}
		for j, h := range headers {
			if h == key {
				indices[i] = j
				break
			}
		}
	}
	return indices
}

func buildKeyMap(rows [][]string, keyIndices []int) map[string][]string {
	m := make(map[string][]string)
	for _, row := range rows {
		k := extractKeyString(row, keyIndices)
		m[k] = row
	}
	return m
}

func extractKey(row []string, keyIndices []int) []string {
	vals := make([]string, len(keyIndices))
	for i, idx := range keyIndices {
		if idx < len(row) {
			vals[i] = row[idx]
		}
	}
	return vals
}

func extractKeyString(row []string, keyIndices []int) string {
	return strings.Join(extractKey(row, keyIndices), "\x00")
}

func isKeyIndex(i int, keyIndices []int) bool {
	for _, idx := range keyIndices {
		if i == idx {
			return true
		}
	}
	return false
}

func rowKey(row []string) string {
	return strings.Join(row, "\x00")
}

func checkDuplicateKeys(rows [][]string, keyIndices []int, label string) []string {
	var warnings []string
	seen := make(map[string]bool)
	for _, row := range rows {
		k := extractKeyString(row, keyIndices)
		if seen[k] {
			warnings = append(warnings, fmt.Sprintf("duplicate key %q in %s", k, label))
		}
		seen[k] = true
	}
	return warnings
}
