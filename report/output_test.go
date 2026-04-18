package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []string{"text", "json", "csv", "table", "markdown", "curve"}
	for _, c := range cases {
		if _, err := ParseFormat(c); err != nil {
			t.Errorf("expected %q to be valid: %v", c, err)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	if _, err := ParseFormat("nope"); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_Formats(t *testing.T) {
	r := &Report{Total: 2, Success: 2}
	results := makeCurveResults()

	formats := []Format{
		FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown,
		FormatPrometheus, FormatHTML, FormatXML, FormatInflux, FormatDotPlot,
		FormatSparkline, FormatHeatmap, FormatTimeline, FormatFlamegraph, FormatCurve,
	}
	for _, f := range formats {
		var buf bytes.Buffer
		if err := Write(&buf, r, results, f); err != nil {
			t.Errorf("Write(%q) returned error: %v", f, err)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, &Report{}, nil, Format("bogus"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("error should mention format name: %v", err)
	}
}
