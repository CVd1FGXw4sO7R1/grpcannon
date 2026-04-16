package report

import (
	"bytes"
	"testing"
	"time"
)

func TestParseFormat_Valid(t *testing.T) {
	formats := []string{"text", "json", "csv", "table", "markdown",
		"prometheus", "html", "xml", "influx", "dotplot", "sparkline"}
	for _, f := range formats {
		if _, err := ParseFormat(f); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", f, err)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	if _, err := ParseFormat("nope"); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_Formats(t *testing.T) {
	results := makeResults(5, 1, 2*time.Millisecond, 20*time.Millisecond)
	r := New(results)
	formats := []Format{
		FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown,
		FormatPrometheus, FormatHTML, FormatXML, FormatInflux, FormatDotPlot, FormatSparkline,
	}
	for _, f := range formats {
		var buf bytes.Buffer
		if err := Write(&buf, r, f); err != nil {
			t.Errorf("Write(%q) error: %v", f, err)
		}
		if buf.Len() == 0 {
			t.Errorf("Write(%q) produced no output", f)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	r := New(nil)
	var buf bytes.Buffer
	if err := Write(&buf, r, Format("unknown")); err == nil {
		t.Error("expected error for unknown format")
	}
}
