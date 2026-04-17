package report

import (
	"bytes"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []string{"text", "json", "csv", "table", "markdown", "rps"}
	for _, c := range cases {
		if _, err := ParseFormat(c); err != nil {
			t.Errorf("expected %q to be valid, got %v", c, err)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	if _, err := ParseFormat("bogus"); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_Formats(t *testing.T) {
	results := makeResults()
	r := New(results)
	formats := []Format{
		FormatText, FormatJSON, FormatCSV, FormatTable,
		FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML,
		FormatInflux, FormatDotPlot, FormatSparkline, FormatHeatmap,
		FormatTimeline, FormatFlamegraph, FormatRPS,
	}
	for _, f := range formats {
		var buf bytes.Buffer
		if err := Write(&buf, r, results, f); err != nil {
			t.Errorf("Write(%q) returned error: %v", f, err)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	results := makeResults()
	r := New(results)
	var buf bytes.Buffer
	// unknown falls back to text — no error expected
	if err := Write(&buf, r, results, Format("unknown")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
