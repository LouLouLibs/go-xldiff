# go-xldiff

Compare two Excel sheets and report added, removed, and modified rows.

`go-xldiff` reads two `.xlsx` sheets, diffs them row by row, and prints the differences to stdout. It supports positional comparison (all columns must match) and key-based comparison (match rows by designated columns, then report cell-level changes).

## Installation

```bash
go install github.com/louloulibs/go-xldiff@latest
```

Or build from source:

```bash
git clone https://github.com/louloulibs/go-xldiff.git
cd go-xldiff
go build -o go-xldiff .
```

## Quick Start

```bash
# Compare two files (first sheet, all columns)
go-xldiff old.xlsx new.xlsx

# Compare by a key column
go-xldiff old.xlsx new.xlsx --key ID

# Compare specific sheets
go-xldiff data.xlsx:January data.xlsx:February --key Date,Ticker

# Output as JSON
go-xldiff old.xlsx new.xlsx --key ID --format json
```

## Usage

```
go-xldiff <file1>[:<sheet>] <file2>[:<sheet>] [flags]
```

### Arguments

Each argument is a path to an `.xlsx` file, optionally followed by `:<sheet>` to select a specific sheet.

| Syntax | Meaning |
|--------|---------|
| `file.xlsx` | Use the first sheet |
| `file.xlsx:Sales` | Use the sheet named "Sales" |
| `file.xlsx:2` | Use the third sheet (0-indexed) |

Sheet selection tries an exact name match first. If no sheet matches and the value is a number, it uses the value as a 0-based index.

To compare two sheets within the same file:

```bash
go-xldiff report.xlsx:Q1 report.xlsx:Q2
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--key <cols>` | *(all columns)* | Columns that identify a row. Accepts header names or 0-based indices, comma-separated. |
| `--skip <n>[,<m>]` | `0` | Rows to skip before the header. A single value applies to both sheets; a comma-separated pair sets each independently. |
| `--format <fmt>` | `text` | Output format: `text`, `json`, or `csv`. |
| `--no-header` | `false` | Treat the first data row as data rather than a header. Columns are labeled `Col0`, `Col1`, etc. |
| `--no-color` | `false` | Suppress ANSI color codes in text output. |

## Diff Modes

### Positional Mode (default)

When no `--key` is specified, every column defines row identity. A row either matches exactly or differs. The tool reports rows unique to the first sheet as **removed** and rows unique to the second sheet as **added**. The tool counts duplicate rows correctly (multiset comparison).

This mode omits **modified** rows because any column change creates a new identity.

### Key-Based Mode

When `--key` is specified, the tool matches rows across sheets by their key column values, then compares the remaining columns cell by cell.

```bash
# Single key
go-xldiff old.xlsx new.xlsx --key ID

# Composite key (neither column is unique alone)
go-xldiff old.xlsx new.xlsx --key Date,Ticker

# Key by column index
go-xldiff old.xlsx new.xlsx --key 0
```

Key-based mode reports three categories:

- **Added** -- rows whose key exists only in the second sheet
- **Removed** -- rows whose key exists only in the first sheet
- **Modified** -- rows sharing a key across both sheets but differing in non-key columns, showing per-cell old/new values

If a sheet contains duplicate keys, `go-xldiff` prints a warning to stderr and continues.

## Output Formats

### Text (default)

Colored terminal output resembling `diff`. Removed rows print in red, added rows in green, modified rows in yellow. Use `--no-color` to strip ANSI codes when piping output.

```
--- Sheet1 (old.xlsx)
+++ Sheet1 (new.xlsx)

Added rows: 1
Removed rows: 1
Modified rows: 1

- ID: "3"  Name: "Charlie"  Score: "70"
+ ID: "4"  Name: "Dana"  Score: "85"

~ [1]  Score: "90" -> "95"
```

### JSON

Structured output suitable for scripting and further processing.

```json
{
  "added": [
    {"ID": "4", "Name": "Dana", "Score": "85"}
  ],
  "removed": [
    {"ID": "3", "Name": "Charlie", "Score": "70"}
  ],
  "modified": [
    {
      "key": {"ID": "1"},
      "changes": [
        {"column": "Score", "old": "90", "new": "95"}
      ]
    }
  ]
}
```

Empty categories appear as `[]`, never `null`.

### CSV

Flat format with a `_status` column and `_old_<column>` columns for previous values.

```csv
_status,ID,Name,Score,_old_ID,_old_Name,_old_Score
added,4,Dana,85,,,
removed,3,Charlie,70,,,
modified,,,,,,90
```

For modified rows, only changed columns populate the current and `_old_` positions; unchanged columns remain empty.

## Skipping Rows

Many spreadsheets have title rows, metadata, or blank lines before the data table. Use `--skip` to ignore them.

```bash
# Skip 2 rows in both sheets before the header
go-xldiff file1.xlsx file2.xlsx --skip 2

# Skip 3 rows in the first sheet, 5 in the second
go-xldiff file1.xlsx file2.xlsx --skip 3,5
```

The first row after the skipped rows becomes the header row (or the first data row when using `--no-header`).

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | No differences found |
| `1` | Differences found |
| `2` | Error (bad arguments, file not found, etc.) |

This follows the convention of `diff(1)`, making `go-xldiff` easy to use in shell scripts:

```bash
if go-xldiff old.xlsx new.xlsx --key ID --no-color > /dev/null 2>&1; then
    echo "Files are identical"
else
    echo "Files differ"
fi
```

## Architecture

```
go-xldiff/
  main.go              Entry point
  cmd/root.go           CLI parsing (cobra)
  internal/
    reader/             Excel -> Table (excelize)
    diff/               Table x Table -> DiffResult
    output/             DiffResult -> stdout (text/json/csv)
```

The three internal packages are fully decoupled: the reader ignores diffs, the diff engine ignores Excel, and the formatters ignore both. Each accepts and returns plain data structures.

## Development

```bash
# Run all tests
go test ./... -v

# Run a single package
go test ./internal/diff/ -v

# Build
go build -o go-xldiff .

# Lint
go vet ./...
```

## License

MIT
