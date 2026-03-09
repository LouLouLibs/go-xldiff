package output_test

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/output"
)

func TestCSVOutput(t *testing.T) {
	result := &diff.DiffResult{
		Headers: []string{"ID", "Name"},
		Added:   []diff.Row{{Values: []string{"2", "Bob"}}},
		Removed: []diff.Row{{Values: []string{"3", "Charlie"}}},
		Modified: []diff.RowDiff{{
			Key:     []string{"1"},
			Changes: []diff.CellChange{{Column: "Name", OldValue: "Alice", NewValue: "Alicia"}},
		}},
	}
	var buf bytes.Buffer
	if err := output.WriteCSV(&buf, result); err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 4 {
		t.Errorf("rows = %d, want 4 (1 header + 3 data)", len(records))
	}
	if records[0][0] != "_status" {
		t.Errorf("first header = %q, want _status", records[0][0])
	}
}

func TestCSVOutput_Empty(t *testing.T) {
	result := &diff.DiffResult{Headers: []string{"ID"}}
	var buf bytes.Buffer
	if err := output.WriteCSV(&buf, result); err != nil {
		t.Fatal(err)
	}
	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Errorf("rows = %d, want 1 (header only)", len(records))
	}
}
