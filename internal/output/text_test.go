package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/louloulibs/go-xldiff/internal/diff"
	"github.com/louloulibs/go-xldiff/internal/output"
	"github.com/louloulibs/go-xldiff/internal/reader"
)

func TestTextOutput_NoChanges(t *testing.T) {
	result := &diff.DiffResult{Headers: []string{"ID", "Name"}}
	var buf bytes.Buffer
	output.WriteText(&buf, result, &reader.Table{FileName: "a.xlsx", Sheet: "Sheet1"}, &reader.Table{FileName: "b.xlsx", Sheet: "Sheet1"}, true)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected 'No differences', got %q", buf.String())
	}
}

func TestTextOutput_WithChanges(t *testing.T) {
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
	output.WriteText(&buf, result, &reader.Table{FileName: "a.xlsx", Sheet: "Sheet1"}, &reader.Table{FileName: "b.xlsx", Sheet: "Sheet1"}, true)
	out := buf.String()

	if !strings.Contains(out, "Added rows: 1") {
		t.Errorf("missing added count")
	}
	if !strings.Contains(out, "Removed rows: 1") {
		t.Errorf("missing removed count")
	}
	if !strings.Contains(out, "Modified rows: 1") {
		t.Errorf("missing modified count")
	}
}

func TestTextOutput_NoColor(t *testing.T) {
	result := &diff.DiffResult{
		Headers: []string{"ID"},
		Added:   []diff.Row{{Values: []string{"1"}}},
	}
	var buf bytes.Buffer
	output.WriteText(&buf, result, &reader.Table{FileName: "a.xlsx", Sheet: "S1"}, &reader.Table{FileName: "b.xlsx", Sheet: "S1"}, true)
	out := buf.String()
	if strings.Contains(out, "\033[") {
		t.Errorf("no-color output contains ANSI escape codes")
	}
}
