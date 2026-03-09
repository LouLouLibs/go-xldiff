# go-xldiff Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go CLI tool that diffs two Excel sheets and reports added, removed, and modified rows.

**Architecture:** Layered internal packages — reader (Excel → Table), diff (Table × Table → DiffResult), output (DiffResult → stdout). CLI parsing with cobra. TDD throughout.

**Tech Stack:** Go 1.26, excelize/v2, cobra, standard library (encoding/json, encoding/csv)

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `cmd/root.go`

**Step 1: Initialize Go module**

Run: `cd /Users/loulou/Dropbox/projects_claude/go-xldiff && go mod init github.com/loulou/go-xldiff`
Expected: `go.mod` created

**Step 2: Install dependencies**

Run: `go get github.com/xuri/excelize/v2 github.com/spf13/cobra`
Expected: Dependencies added to `go.mod` and `go.sum`

**Step 3: Create main.go**

```go
package main

import "github.com/loulou/go-xldiff/cmd"

func main() {
	cmd.Execute()
}
```

**Step 4: Create cmd/root.go with flag definitions**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	keyFlag      string
	skipFlag     string
	formatFlag   string
	noHeaderFlag bool
	noColorFlag  bool
)

var rootCmd = &cobra.Command{
	Use:   "go-xldiff <file1>[:<sheet>] <file2>[:<sheet>]",
	Short: "Diff two Excel sheets",
	Long:  "Compare two Excel sheets and report added, removed, and modified rows.",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.Flags().StringVar(&keyFlag, "key", "", "Row identity columns (header name or 0-based index), comma-separated")
	rootCmd.Flags().StringVar(&skipFlag, "skip", "0", "Rows to skip before header. Single value or comma-separated pair (e.g. 3,5)")
	rootCmd.Flags().StringVar(&formatFlag, "format", "text", "Output format: text, json, csv")
	rootCmd.Flags().BoolVar(&noHeaderFlag, "no-header", false, "Treat first row as data, not headers")
	rootCmd.Flags().BoolVar(&noColorFlag, "no-color", false, "Disable colored terminal output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDiff(cmd *cobra.Command, args []string) error {
	fmt.Println("Not implemented yet")
	return nil
}
```

**Step 5: Verify it compiles and runs**

Run: `go build -o go-xldiff . && ./go-xldiff --help`
Expected: Help text with flag descriptions

**Step 6: Commit**

```bash
git add main.go cmd/root.go go.mod go.sum
git commit -m "feat: scaffold project with cobra CLI and flag definitions"
```

---

### Task 2: Argument Parsing — Parse file:sheet and --skip

**Files:**
- Create: `internal/reader/parse.go`
- Create: `internal/reader/parse_test.go`

**Step 1: Write failing tests for argument parsing**

```go
package reader_test

import (
	"testing"

	"github.com/loulou/go-xldiff/internal/reader"
)

func TestParseFileArg(t *testing.T) {
	tests := []struct {
		input     string
		wantFile  string
		wantSheet string
	}{
		{"data.xlsx", "data.xlsx", ""},
		{"data.xlsx:Sheet2", "data.xlsx", "Sheet2"},
		{"data.xlsx:0", "data.xlsx", "0"},
		{"/path/to/file.xlsx:MySheet", "/path/to/file.xlsx", "MySheet"},
		{"file.xlsx:", "file.xlsx", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			file, sheet := reader.ParseFileArg(tt.input)
			if file != tt.wantFile {
				t.Errorf("file = %q, want %q", file, tt.wantFile)
			}
			if sheet != tt.wantSheet {
				t.Errorf("sheet = %q, want %q", sheet, tt.wantSheet)
			}
		})
	}
}

func TestParseSkipFlag(t *testing.T) {
	tests := []struct {
		input string
		want1 int
		want2 int
		err   bool
	}{
		{"0", 0, 0, false},
		{"3", 3, 3, false},
		{"3,5", 3, 5, false},
		{"abc", 0, 0, true},
		{"3,", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			s1, s2, err := reader.ParseSkipFlag(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("err = %v, wantErr %v", err, tt.err)
				return
			}
			if err == nil && (s1 != tt.want1 || s2 != tt.want2) {
				t.Errorf("got (%d, %d), want (%d, %d)", s1, s2, tt.want1, tt.want2)
			}
		})
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/reader/ -v`
Expected: FAIL — package/functions don't exist

**Step 3: Implement parse.go**

```go
package reader

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFileArg splits "file.xlsx:SheetName" into file path and sheet identifier.
// If no colon, sheet is empty (meaning first sheet).
func ParseFileArg(arg string) (file, sheet string) {
	idx := strings.LastIndex(arg, ":")
	if idx == -1 {
		return arg, ""
	}
	return arg[:idx], arg[idx+1:]
}

// ParseSkipFlag parses "--skip N" or "--skip N,M" into two skip values.
func ParseSkipFlag(flag string) (skip1, skip2 int, err error) {
	parts := strings.Split(flag, ",")
	switch len(parts) {
	case 1:
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid skip value %q: %w", parts[0], err)
		}
		return n, n, nil
	case 2:
		n1, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid skip value %q: %w", parts[0], err)
		}
		n2, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid skip value %q: %w", parts[1], err)
		}
		return n1, n2, nil
	default:
		return 0, 0, fmt.Errorf("invalid skip flag %q: expected N or N,M", flag)
	}
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/reader/ -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/reader/parse.go internal/reader/parse_test.go
git commit -m "feat: add file:sheet argument and --skip flag parsing"
```

---

### Task 3: Excel Reader — Load Sheet into Table

**Files:**
- Create: `internal/reader/reader.go`
- Modify: `internal/reader/parse_test.go` → rename to `internal/reader/reader_test.go` (merge all tests)

**Step 1: Define the Table type and write failing tests**

Add to `internal/reader/reader_test.go`:

```go
import "github.com/xuri/excelize/v2"

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
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/reader/ -v`
Expected: FAIL — ReadSheet not defined

**Step 3: Implement reader.go**

```go
package reader

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Table struct {
	Headers  []string
	Rows     [][]string
	FileName string
	Sheet    string
}

func ReadSheet(filePath, sheet string, skip int, noHeader bool) (*Table, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", filePath, err)
	}
	defer f.Close()

	sheetName, err := resolveSheet(f, sheet)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("reading rows from %s: %w", sheetName, err)
	}

	if skip > 0 {
		if skip >= len(rows) {
			return &Table{FileName: filePath, Sheet: sheetName}, nil
		}
		rows = rows[skip:]
	}

	if len(rows) == 0 {
		return &Table{FileName: filePath, Sheet: sheetName}, nil
	}

	var headers []string
	var dataRows [][]string

	if noHeader {
		maxCols := 0
		for _, row := range rows {
			if len(row) > maxCols {
				maxCols = len(row)
			}
		}
		headers = make([]string, maxCols)
		for i := range headers {
			headers[i] = fmt.Sprintf("Col%d", i)
		}
		dataRows = rows
	} else {
		headers = rows[0]
		dataRows = rows[1:]
	}

	normalized := make([][]string, len(dataRows))
	for i, row := range dataRows {
		n := make([]string, len(headers))
		copy(n, row)
		normalized[i] = n
	}

	return &Table{
		Headers:  headers,
		Rows:     normalized,
		FileName: filePath,
		Sheet:    sheetName,
	}, nil
}

func resolveSheet(f *excelize.File, sheet string) (string, error) {
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("workbook has no sheets")
	}

	if sheet == "" {
		return sheets[0], nil
	}

	for _, s := range sheets {
		if s == sheet {
			return s, nil
		}
	}

	idx, err := strconv.Atoi(sheet)
	if err == nil && idx >= 0 && idx < len(sheets) {
		return sheets[idx], nil
	}

	return "", fmt.Errorf("sheet %q not found", sheet)
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/reader/ -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/reader/reader.go internal/reader/reader_test.go
git commit -m "feat: add Excel reader with sheet selection, skip, and no-header support"
```

---

### Task 4: Diff Engine — Positional Mode

**Files:**
- Create: `internal/diff/diff.go`
- Create: `internal/diff/diff_test.go`

**Step 1: Define types and write failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/diff/ -v`
Expected: FAIL — package doesn't exist

**Step 3: Implement diff.go — types and positional mode**

```go
package diff

import (
	"strings"

	"github.com/loulou/go-xldiff/internal/reader"
)

type Row struct {
	Values []string
}

type CellChange struct {
	Column   string
	OldValue string
	NewValue string
}

type RowDiff struct {
	Key     []string
	Changes []CellChange
}

type DiffResult struct {
	Headers  []string
	Added    []Row
	Removed  []Row
	Modified []RowDiff
}

func (r *DiffResult) HasDifferences() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

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

	// Removed: in A but not enough in B
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

	// Added: in B but not enough in A
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
	// Placeholder — implemented in Task 5
}

func rowKey(row []string) string {
	return strings.Join(row, "\x00")
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/diff/ -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/diff/diff.go internal/diff/diff_test.go
git commit -m "feat: add diff engine with positional (all-column) mode"
```

---

### Task 5: Diff Engine — Key-Based Mode

**Files:**
- Modify: `internal/diff/diff.go`
- Modify: `internal/diff/diff_test.go`

**Step 1: Write failing tests for key-based mode**

Append to `internal/diff/diff_test.go`:

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/diff/ -v`
Expected: FAIL — key-based tests fail

**Step 3: Implement compareByKey in diff.go**

Replace the `compareByKey` placeholder and add helpers:

```go
import "strconv"

func compareByKey(a, b *reader.Table, keys []string, result *DiffResult) {
	keyIndices := resolveKeyIndices(a.Headers, keys)

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
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/diff/ -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/diff/diff.go internal/diff/diff_test.go
git commit -m "feat: add key-based diff mode with composite key support"
```

---

### Task 6: Text Output Formatter

**Files:**
- Create: `internal/output/text.go`
- Create: `internal/output/text_test.go`

**Step 1: Write failing tests**

```go
package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/output"
	"github.com/loulou/go-xldiff/internal/reader"
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
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/output/ -v`
Expected: FAIL

**Step 3: Implement text.go**

```go
package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/reader"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func WriteText(w io.Writer, result *diff.DiffResult, a, b *reader.Table, noColor bool) {
	red, green, yellow, reset := colorRed, colorGreen, colorYellow, colorReset
	if noColor {
		red, green, yellow, reset = "", "", "", ""
	}

	fmt.Fprintf(w, "--- %s (%s)\n", a.Sheet, a.FileName)
	fmt.Fprintf(w, "+++ %s (%s)\n", b.Sheet, b.FileName)
	fmt.Fprintln(w)

	if !result.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	fmt.Fprintf(w, "Added rows: %d\n", len(result.Added))
	fmt.Fprintf(w, "Removed rows: %d\n", len(result.Removed))
	fmt.Fprintf(w, "Modified rows: %d\n", len(result.Modified))
	fmt.Fprintln(w)

	for _, row := range result.Removed {
		fmt.Fprintf(w, "%s- %s%s\n", red, formatRow(result.Headers, row.Values), reset)
	}
	for _, row := range result.Added {
		fmt.Fprintf(w, "%s+ %s%s\n", green, formatRow(result.Headers, row.Values), reset)
	}
	if len(result.Modified) > 0 && (len(result.Added) > 0 || len(result.Removed) > 0) {
		fmt.Fprintln(w)
	}
	for _, mod := range result.Modified {
		parts := make([]string, len(mod.Changes))
		for i, c := range mod.Changes {
			parts[i] = fmt.Sprintf("%s: %q → %q", c.Column, c.OldValue, c.NewValue)
		}
		keyStr := strings.Join(mod.Key, ", ")
		fmt.Fprintf(w, "%s~ [%s]  %s%s\n", yellow, keyStr, strings.Join(parts, ", "), reset)
	}
}

func formatRow(headers, values []string) string {
	parts := make([]string, len(headers))
	for i, h := range headers {
		v := ""
		if i < len(values) {
			v = values[i]
		}
		parts[i] = fmt.Sprintf("%s: %q", h, v)
	}
	return strings.Join(parts, "  ")
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/output/ -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/output/text.go internal/output/text_test.go
git commit -m "feat: add colored text output formatter"
```

---

### Task 7: JSON Output Formatter

**Files:**
- Create: `internal/output/json.go`
- Create: `internal/output/json_test.go`

**Step 1: Write failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/output/ -v -run TestJSON`
Expected: FAIL

**Step 3: Implement json.go**

```go
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
		Modified: modifiedToJSON(result.Modified),
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

func modifiedToJSON(mods []diff.RowDiff) []jsonModified {
	var result []jsonModified
	for _, mod := range mods {
		keyMap := make(map[string]string)
		for i, v := range mod.Key {
			keyMap[fmt.Sprintf("key%d", i)] = v
		}
		var changes []jsonChange
		for _, c := range mod.Changes {
			changes = append(changes, jsonChange{Column: c.Column, Old: c.OldValue, New: c.NewValue})
		}
		result = append(result, jsonModified{Key: keyMap, Changes: changes})
	}
	return result
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/output/ -v -run TestJSON`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/output/json.go internal/output/json_test.go
git commit -m "feat: add JSON output formatter"
```

---

### Task 8: CSV Output Formatter

**Files:**
- Create: `internal/output/csv.go`
- Create: `internal/output/csv_test.go`

**Step 1: Write failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/output/ -v -run TestCSV`
Expected: FAIL

**Step 3: Implement csv.go**

```go
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
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/output/ -v -run TestCSV`
Expected: All PASS

**Step 5: Commit**

```bash
git add internal/output/csv.go internal/output/csv_test.go
git commit -m "feat: add CSV output formatter"
```

---

### Task 9: Wire Everything Together in cmd/root.go

**Files:**
- Modify: `cmd/root.go`

**Step 1: Implement runDiff in cmd/root.go**

Replace the placeholder `runDiff` function:

```go
import (
	"os"
	"strings"

	"github.com/loulou/go-xldiff/internal/diff"
	"github.com/loulou/go-xldiff/internal/output"
	"github.com/loulou/go-xldiff/internal/reader"
	"github.com/spf13/cobra"
)

func runDiff(cmd *cobra.Command, args []string) error {
	file1, sheet1 := reader.ParseFileArg(args[0])
	file2, sheet2 := reader.ParseFileArg(args[1])

	skip1, skip2, err := reader.ParseSkipFlag(skipFlag)
	if err != nil {
		return fmt.Errorf("invalid --skip: %w", err)
	}

	table1, err := reader.ReadSheet(file1, sheet1, skip1, noHeaderFlag)
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	table2, err := reader.ReadSheet(file2, sheet2, skip2, noHeaderFlag)
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	var keys []string
	if keyFlag != "" {
		keys = strings.Split(keyFlag, ",")
	}

	result := diff.Compare(table1, table2, keys)

	switch formatFlag {
	case "json":
		if err := output.WriteJSON(os.Stdout, result); err != nil {
			return err
		}
	case "csv":
		if err := output.WriteCSV(os.Stdout, result); err != nil {
			return err
		}
	default:
		output.WriteText(os.Stdout, result, table1, table2, noColorFlag)
	}

	if result.HasDifferences() {
		os.Exit(1)
	}
	return nil
}
```

**Step 2: Verify the build**

Run: `go build -o go-xldiff .`
Expected: Compiles without errors

**Step 3: Manual smoke test**

Create two small test xlsx files and run the tool to verify end-to-end behavior.

**Step 4: Run all tests**

Run: `go test ./... -v`
Expected: All PASS

**Step 5: Commit**

```bash
git add cmd/root.go
git commit -m "feat: wire CLI together — reader, diff, and output formatters"
```

---

### Task 10: Duplicate Key Warning and Edge Cases

**Files:**
- Modify: `internal/diff/diff.go`
- Modify: `internal/diff/diff_test.go`

**Step 1: Write failing test for duplicate key warning**

```go
func TestKeyDiff_DuplicateKeyWarning(t *testing.T) {
	a := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Alice"}, {"1", "Bob"}},
	}
	b := &reader.Table{
		Headers: []string{"ID", "Name"},
		Rows:    [][]string{{"1", "Charlie"}},
	}

	warnings := diff.CompareWithWarnings(a, b, []string{"ID"})
	if len(warnings) == 0 {
		t.Error("expected duplicate key warning")
	}
}
```

**Step 2: Implement CompareWithWarnings**

Add a `CompareWithWarnings` function that returns `(DiffResult, []string)` — the second return is a list of warning messages. Detect duplicate keys in both tables before diffing.

```go
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
```

Add `"fmt"` and `"os"` to imports. Update `cmd/root.go` to use `CompareWithWarnings` and print warnings to stderr.

**Step 3: Run tests to verify they pass**

Run: `go test ./... -v`
Expected: All PASS

**Step 4: Commit**

```bash
git add internal/diff/diff.go internal/diff/diff_test.go cmd/root.go
git commit -m "feat: add duplicate key warnings to stderr"
```

---

### Task 11: Final Polish

**Files:**
- Create: `CLAUDE.md`
- Create: `.gitignore`

**Step 1: Create .gitignore**

```
go-xldiff
*.xlsx
!testdata/**/*.xlsx
```

**Step 2: Create CLAUDE.md**

```markdown
# go-xldiff

Excel sheet diff tool written in Go.

## Build & Test

- Build: `go build -o go-xldiff .`
- Test: `go test ./... -v`
- Single package: `go test ./internal/diff/ -v`

## Architecture

- `cmd/` — CLI (cobra)
- `internal/reader/` — Excel → Table
- `internal/diff/` — Table × Table → DiffResult
- `internal/output/` — DiffResult → stdout (text/json/csv)
```

**Step 3: Run final verification**

Run: `go build -o go-xldiff . && go test ./... -v && go vet ./...`
Expected: All pass, no warnings

**Step 4: Commit**

```bash
git add CLAUDE.md .gitignore
git commit -m "chore: add CLAUDE.md and .gitignore"
```
