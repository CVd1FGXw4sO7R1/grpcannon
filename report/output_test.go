package report

import (
	"bytes"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []string{"text", "json", "csv", "table", "markdown", "prometheus", "html", "TEXT", "JSON", "HTML"}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			_, err := ParseFormat(c)
			if err != nil {
				t.Errorf("expected no error for %q, got %v", c, err)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_Formats(t *testing.T) {
	results := makeResults(4, 1)
	r := New(results)

	formats := []Format{
		FormatText, FormatJSON, FormatCSV, FormatTable,
		FormatMarkdown, FormatPrometheus, FormatHTML,
	}

	for _, f := range formats {
		t.Run(string(f), func(t *testing.T) {
			var buf bytes.Buffer
			if err := Write(&buf, r, f); err != nil {
				t.Errorf("Write(%s) error: %v", f, err)
			}
			if buf.Len() == 0 {
				t.Errorf("Write(%s) produced empty output", f)
			}
		})
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	results := makeResults(2, 0)
	r := New(results)
	var buf bytes.Buffer
	err := Write(&buf, r, Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
