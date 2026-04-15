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
		{"json", FormatJSON},
		{"csv", FormatCSV},
		{"table", FormatTable},
		{"markdown", FormatMarkdown},
		{"prometheus", FormatPrometheus},
		{"TEXT", FormatText},
		{"JSON", FormatJSON},
		{"Prometheus", FormatPrometheus},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			f, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, f)
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
	r := New(results, 500*time.Millisecond)

	formats := []Format{
		FormatText, FormatJSON, FormatCSV,
		FormatTable, FormatMarkdown, FormatPrometheus,
	}

	for _, f := range formats {
		t.Run(string(f), func(t *testing.T) {
			var buf bytes.Buffer
			err := Write(&buf, r, f)
			if err != nil {
				t.Fatalf("Write(%q) error: %v", f, err)
			}
			if buf.Len() == 0 {
				t.Errorf("Write(%q) produced empty output", f)
			}
		})
	}
}
