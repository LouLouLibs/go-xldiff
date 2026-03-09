package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/loulou/go-xldiff/internal/diff"
)

type jsonOutput struct {
	Added    []map[string]string `json:"added"`
	Removed  []map[string]string `json:"removed"`
	Modified []jsonModified      `json:"modified"`
}

type jsonModified struct {
	Key     map[string]string `json:"key"`
	Changes []jsonChange      `json:"changes"`
}

type jsonChange struct {
	Column string `json:"column"`
	Old    string `json:"old"`
	New    string `json:"new"`
}

func WriteJSON(w io.Writer, result *diff.DiffResult) error {
	out := jsonOutput{
		Added:    rowsToMaps(result.Headers, result.Added),
		Removed:  rowsToMaps(result.Headers, result.Removed),
		Modified: modifiedToJSON(result.Modified, result.KeyColumns),
	}
	if out.Added == nil {
		out.Added = []map[string]string{}
	}
	if out.Removed == nil {
		out.Removed = []map[string]string{}
	}
	if out.Modified == nil {
		out.Modified = []jsonModified{}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func rowsToMaps(headers []string, rows []diff.Row) []map[string]string {
	var result []map[string]string
	for _, row := range rows {
		m := make(map[string]string)
		for i, h := range headers {
			if i < len(row.Values) {
				m[h] = row.Values[i]
			}
		}
		result = append(result, m)
	}
	return result
}

func modifiedToJSON(mods []diff.RowDiff, keyColumns []string) []jsonModified {
	var result []jsonModified
	for _, mod := range mods {
		keyMap := make(map[string]string)
		for i, v := range mod.Key {
			name := fmt.Sprintf("key%d", i)
			if i < len(keyColumns) && keyColumns[i] != "" {
				name = keyColumns[i]
			}
			keyMap[name] = v
		}
		var changes []jsonChange
		for _, c := range mod.Changes {
			changes = append(changes, jsonChange{Column: c.Column, Old: c.OldValue, New: c.NewValue})
		}
		result = append(result, jsonModified{Key: keyMap, Changes: changes})
	}
	return result
}
