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
