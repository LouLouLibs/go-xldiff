package output

import (
	"encoding/csv"
	"io"

	"github.com/loulou/go-xldiff/internal/diff"
)

func WriteCSV(w io.Writer, result *diff.DiffResult) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"_status"}
	header = append(header, result.Headers...)
	for _, h := range result.Headers {
		header = append(header, "_old_"+h)
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, row := range result.Added {
		record := []string{"added"}
		record = append(record, row.Values...)
		for range result.Headers {
			record = append(record, "")
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	for _, row := range result.Removed {
		record := []string{"removed"}
		record = append(record, row.Values...)
		for range result.Headers {
			record = append(record, "")
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	for _, mod := range result.Modified {
		record := []string{"modified"}
		changeMap := make(map[string]diff.CellChange)
		for _, c := range mod.Changes {
			changeMap[c.Column] = c
		}
		for _, h := range result.Headers {
			if c, ok := changeMap[h]; ok {
				record = append(record, c.NewValue)
			} else {
				record = append(record, "")
			}
		}
		for _, h := range result.Headers {
			if c, ok := changeMap[h]; ok {
				record = append(record, c.OldValue)
			} else {
				record = append(record, "")
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
