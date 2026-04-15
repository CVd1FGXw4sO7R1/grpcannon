package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"text", FormatText},
		{"TEXT", FormatText},
		{"", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"csv", FormatCSV},
		{"CSV", FormatCSV},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			f, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, f)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention the bad format, got: %v", err)
	}
}

func TestWrite_Formats(t *testing.T) {
	results := makeResults(5, 1)
	r := New(results, time.Second)

	formats := []Format{FormatText, FormatJSON, FormatCSV}
	for _, f := range formats {
		t.Run(fmt.Sprintf("format_%d", f), func(t *testing.T) {
			var buf bytes.Buffer
			if err := Write(r, f, &buf); err != nil {
				t.Fatalf("Write failed: %v", err)
			}
			if buf.Len() == 0 {
				t.Error("expected non-empty output")
			}
		})
	}
}
