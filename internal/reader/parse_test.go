package reader_test

import (
	"testing"

	"github.com/louloulibs/go-xldiff/internal/reader"
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
