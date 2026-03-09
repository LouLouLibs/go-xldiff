package reader_test

import (
	"testing"

	"github.com/loulou/go-xldiff/internal/reader"
	"github.com/xuri/excelize/v2"
)

func TestReadSheet(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "A2", "1")
	f.SetCellValue("Sheet1", "B2", "Alice")
	f.SetCellValue("Sheet1", "A3", "2")
	f.SetCellValue("Sheet1", "B3", "Bob")

	tmpFile := t.TempDir() + "/test.xlsx"
	if err := f.SaveAs(tmpFile); err != nil {
		t.Fatal(err)
	}

	table, err := reader.ReadSheet(tmpFile, "", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Headers) != 2 || table.Headers[0] != "ID" || table.Headers[1] != "Name" {
		t.Errorf("headers = %v, want [ID Name]", table.Headers)
	}
	if len(table.Rows) != 2 {
		t.Errorf("rows = %d, want 2", len(table.Rows))
	}
	if table.Rows[0][0] != "1" || table.Rows[0][1] != "Alice" {
		t.Errorf("row 0 = %v, want [1 Alice]", table.Rows[0])
	}
}

func TestReadSheetWithSkip(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	f.SetCellValue("Sheet1", "A1", "metadata row")
	f.SetCellValue("Sheet1", "A2", "another metadata")
	f.SetCellValue("Sheet1", "A3", "ID")
	f.SetCellValue("Sheet1", "B3", "Name")
	f.SetCellValue("Sheet1", "A4", "1")
	f.SetCellValue("Sheet1", "B4", "Alice")

	tmpFile := t.TempDir() + "/skip.xlsx"
	if err := f.SaveAs(tmpFile); err != nil {
		t.Fatal(err)
	}

	table, err := reader.ReadSheet(tmpFile, "", 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if table.Headers[0] != "ID" {
		t.Errorf("headers[0] = %q, want ID", table.Headers[0])
	}
	if len(table.Rows) != 1 {
		t.Errorf("rows = %d, want 1", len(table.Rows))
	}
}

func TestReadSheetByName(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	f.NewSheet("Data")
	f.SetCellValue("Data", "A1", "Col1")
	f.SetCellValue("Data", "A2", "val1")

	tmpFile := t.TempDir() + "/named.xlsx"
	if err := f.SaveAs(tmpFile); err != nil {
		t.Fatal(err)
	}

	table, err := reader.ReadSheet(tmpFile, "Data", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if table.Headers[0] != "Col1" {
		t.Errorf("headers[0] = %q, want Col1", table.Headers[0])
	}
}

func TestReadSheetByIndex(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	f.NewSheet("Second")
	f.SetCellValue("Second", "A1", "X")
	f.SetCellValue("Second", "A2", "1")

	tmpFile := t.TempDir() + "/indexed.xlsx"
	if err := f.SaveAs(tmpFile); err != nil {
		t.Fatal(err)
	}

	table, err := reader.ReadSheet(tmpFile, "1", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if table.Headers[0] != "X" {
		t.Errorf("headers[0] = %q, want X", table.Headers[0])
	}
}

func TestReadSheetNoHeader(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	f.SetCellValue("Sheet1", "A1", "val1")
	f.SetCellValue("Sheet1", "B1", "val2")

	tmpFile := t.TempDir() + "/noheader.xlsx"
	if err := f.SaveAs(tmpFile); err != nil {
		t.Fatal(err)
	}

	table, err := reader.ReadSheet(tmpFile, "", 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if table.Headers[0] != "Col0" || table.Headers[1] != "Col1" {
		t.Errorf("headers = %v, want [Col0 Col1]", table.Headers)
	}
	if len(table.Rows) != 1 {
		t.Errorf("rows = %d, want 1", len(table.Rows))
	}
}
