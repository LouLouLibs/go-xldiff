# go-xldiff Design

A CLI tool that diffs two Excel sheets and reports row-level and cell-level differences.

## CLI Interface

```
go-xldiff <file1>[:<sheet>] <file2>[:<sheet>]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--key <col1>[,col2,...]` | Row identity columns (header name or 0-based index) | All columns |
| `--skip <n>[,<n>]` | Rows to skip before header. Single value for both sheets, comma-separated for each. | 0 |
| `--format <text\|json\|csv>` | Output format | `text` |
| `--no-header` | Treat first row (after skip) as data, not headers | false |
| `--no-color` | Disable colored terminal output | false |

### Sheet Selection

- `file.xlsx` — uses first sheet
- `file.xlsx:Sheet2` — by name
- `file.xlsx:0` — by index (try name first, fall back to index if it's a number)
- Same-file comparison: `go-xldiff file.xlsx:Sheet1 file.xlsx:Sheet2`

## Architecture

Layered internal packages, each with a single responsibility.

### Data Flow

```
Input files → Reader (excelize, skip rows, extract headers+rows)
    → Normalized Table {Headers []string, Rows [][]string}
    → Diff Engine (positional or key-based matching)
    → DiffResult {Added, Removed, Modified rows with per-cell changes}
    → Formatter (text/json/csv) → stdout
```

### Package Structure

```
go-xldiff/
├── main.go                  # Entry point, wires everything together
├── go.mod
├── cmd/
│   └── root.go              # CLI parsing with cobra, flag definitions
├── internal/
│   ├── reader/
│   │   └── reader.go        # Excel loading via excelize, skip logic, returns Table
│   ├── diff/
│   │   └── diff.go          # Positional & key-based diff engine, returns DiffResult
│   └── output/
│       ├── text.go           # Colored terminal formatter
│       ├── json.go           # JSON formatter
│       └── csv.go            # CSV formatter
```

### Dependencies

- `github.com/xuri/excelize/v2` — Excel reading
- `github.com/spf13/cobra` — CLI parsing
- Standard library for everything else

## Diff Engine

### Positional Mode (default, no `--key`)

All columns define row identity. A row is either an exact match or a difference. Rows present in one sheet but not the other are reported as added/removed. No "modified" state — just rows unique to sheet 1 (removed) and rows unique to sheet 2 (added).

### Key-Based Mode (`--key` specified)

- Build a map of key → row for each sheet
- Keys in both sheets: compare remaining columns, report modified cells
- Keys only in sheet 1: removed
- Keys only in sheet 2: added
- Duplicate keys within a sheet: warn the user

### DiffResult Structure

```go
type DiffResult struct {
    Headers     []string
    Added       []Row
    Removed     []Row
    Modified    []RowDiff
}

type RowDiff struct {
    Key         []string
    Changes     []CellChange
}

type CellChange struct {
    Column      string
    OldValue    string
    NewValue    string
}
```

## Output Formats

### Text (default)

```
--- Sheet1 (file1.xlsx)
+++ Sheet2 (file2.xlsx)

Added rows: 2
Removed rows: 1
Modified rows: 3

- [ID=42]  Name: "Alice"  Score: 95
+ [ID=55]  Name: "Charlie"  Score: 88
+ [ID=60]  Name: "Dana"  Score: 72

~ [ID=10]  Score: 80 → 85
~ [ID=23]  Name: "Bob" → "Robert", Score: 90 → 92
~ [ID=31]  Name: "Eve" → "Eva"
```

Colors: green for `+` (added), red for `-` (removed), yellow for `~` (modified). `--no-color` strips ANSI codes.

### JSON

```json
{
  "added": [{"ID": "55", "Name": "Charlie", "Score": "88"}],
  "removed": [{"ID": "42", "Name": "Alice", "Score": "95"}],
  "modified": [
    {"key": {"ID": "10"}, "changes": [{"column": "Score", "old": "80", "new": "85"}]}
  ]
}
```

### CSV

Flat format with a `_status` column (`added`, `removed`, `modified`) and `_old_<col>` columns for changed values.
