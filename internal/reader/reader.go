package reader

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Table holds the parsed contents of a single Excel sheet.
type Table struct {
	Headers  []string
	Rows     [][]string
	FileName string
	Sheet    string
}

// ReadSheet opens an xlsx file and reads the specified sheet into a Table.
// If sheet is empty, the first sheet is used. If sheet is a numeric string,
// it is treated as a zero-based sheet index. skip controls how many leading
// rows to ignore before reading headers. If noHeader is true, synthetic
// column names (Col0, Col1, ...) are generated and all rows are treated as data.
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
