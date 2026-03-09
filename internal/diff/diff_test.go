package diff_test

import (
	"testing"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/reader"
)

func TestPositionalDiff_NoDifferences(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
	}

	result := diff.Compare(a, b, nil)
	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Modified) != 0 {
		t.Errorf("expected no diffs, got added=%d removed=%d modified=%d",
			len(result.Added), len(result.Removed), len(result.Modified))
	}
}

func TestPositionalDiff_AddedRows(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
	}

	result := diff.Compare(a, b, nil)
	if len(result.Added) != 1 {
		t.Errorf("added = %d, want 1", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Errorf("removed = %d, want 0", len(result.Removed))
	}
}

func TestPositionalDiff_RemovedRows(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}},
	}

	result := diff.Compare(a, b, nil)
	if len(result.Removed) != 1 {
		t.Errorf("removed = %d, want 1", len(result.Removed))
	}
}

func TestPositionalDiff_MixedChanges(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"3", "Charlie"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
	}

	result := diff.Compare(a, b, nil)
	if len(result.Removed) != 1 || len(result.Added) != 1 {
		t.Errorf("removed=%d added=%d, want 1 and 1", len(result.Removed), len(result.Added))
	}
}

func TestKeyDiff_NoDifferences(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name", "Score"},
		Rows:    [][]string{{"1", "Alice", "90"}, {"2", "Bob", "80"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name", "Score"},
		Rows:    [][]string{{"2", "Bob", "80"}, {"1", "Alice", "90"}},
	}

	result := diff.Compare(a, b, []string{"ID"})
	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Modified) != 0 {
		t.Errorf("expected no diffs, got added=%d removed=%d modified=%d",
			len(result.Added), len(result.Removed), len(result.Modified))
	}
}

func TestKeyDiff_Modified(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name", "Score"},
		Rows:    [][]string{{"1", "Alice", "90"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name", "Score"},
		Rows:    [][]string{{"1", "Alice", "95"}},
	}

	result := diff.Compare(a, b, []string{"ID"})
	if len(result.Modified) != 1 {
		t.Fatalf("modified = %d, want 1", len(result.Modified))
	}
	if result.Modified[0].Changes[0].OldValue != "90" || result.Modified[0].Changes[0].NewValue != "95" {
		t.Errorf("change = %+v, want 90 → 95", result.Modified[0].Changes[0])
	}
}

func TestKeyDiff_AddedAndRemoved(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"3", "Charlie"}},
	}

	result := diff.Compare(a, b, []string{"ID"})
	if len(result.Removed) != 1 || result.Removed[0].Values[1] != "Bob" {
		t.Errorf("removed = %v, want Bob", result.Removed)
	}
	if len(result.Added) != 1 || result.Added[0].Values[1] != "Charlie" {
		t.Errorf("added = %v, want Charlie", result.Added)
	}
}

func TestKeyDiff_CompositeKey(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"Date", "Ticker", "Price"},
		Rows:    [][]string{{"2024-01-01", "AAPL", "150"}, {"2024-01-01", "GOOG", "100"}},
	}
	b := &reader.Table{
		Headers: []string{"Date", "Ticker", "Price"},
		Rows:    [][]string{{"2024-01-01", "AAPL", "155"}, {"2024-01-01", "GOOG", "100"}},
	}

	result := diff.Compare(a, b, []string{"Date", "Ticker"})
	if len(result.Modified) != 1 {
		t.Fatalf("modified = %d, want 1", len(result.Modified))
	}
	if result.Modified[0].Changes[0].Column != "Price" {
		t.Errorf("changed column = %q, want Price", result.Modified[0].Changes[0].Column)
	}
}

func TestKeyDiff_DuplicateKeyWarning(t *testing.T) {
	a := &reader.Table{
		Headers:  []string{"ID", "Name"},
		Rows:     [][]string{{"1", "Alice"}, {"1", "Bob"}},
		FileName: "a.xlsx",
	}
	b := &reader.Table{
		Headers:  []string{"ID", "Name"},
		Rows:     [][]string{{"1", "Charlie"}},
		FileName: "b.xlsx",
	}

	_, warnings := diff.CompareWithWarnings(a, b, []string{"ID"})
	if len(warnings) == 0 {
		t.Error("expected duplicate key warning")
	}
}

func TestKeyDiff_NumericKeyIndex(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Bob"}},
	}

	result := diff.Compare(a, b, []string{"0"})
	if len(result.Modified) != 1 {
		t.Fatalf("modified = %d, want 1", len(result.Modified))
	}
}
