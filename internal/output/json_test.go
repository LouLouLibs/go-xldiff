package output_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/output"
)

func TestJSONOutput(t *testing.T) {
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
	if err := output.WriteJSON(&buf, result); err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	added := parsed["added"].([]interface{})
	if len(added) != 1 {
		t.Errorf("added = %d, want 1", len(added))
	}
}

func TestJSONOutput_Empty(t *testing.T) {
	result := &diff.DiffResult{Headers: []string{"ID"}}
	var buf bytes.Buffer
	if err := output.WriteJSON(&buf, result); err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
