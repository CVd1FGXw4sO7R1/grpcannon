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
		{"html", FormatHTML},
		{"xml", FormatXML},
		{"influx", FormatInflux},
		{"hdr", FormatHDR},
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
	_, err := ParseFormat("bogus")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestParseFormat_RoundTrip(t *testing.T) {
	f, err := ParseFormat("hdr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != FormatHDR {
		t.Errorf("expected FormatHDR")
	}
}

func TestWrite_Formats(t *testing.T) {
	results := []Result{
		{Duration: 10 * time.Millisecond},
		{Duration: 20 * time.Millisecond},
	}
	r := New(results)

	formats := []Format{
		FormatText, FormatCSV, FormatTable, FormatMarkdown,
		FormatPrometheus, FormatHTML, FormatInflux,
		FormatDotPlot, FormatSparkline, FormatHeatmap,
		FormatTimeline, FormatFlamegraph, FormatHDR,
	}
	for _, f := range formats {
		var buf bytes.Buffer
		if err := Write(&buf, r, results, f); err != nil {
			t.Errorf("Write format %d returned error: %v", f, err)
		}
		if buf.Len() == 0 {
			t.Errorf("Write format %d produced no output", f)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, nil, nil, Format(9999))
	if err == nil {
		t.Error("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected 'unsupported' in error, got: %v", err)
	}
}

func TestWrite_HDRFormat(t *testing.T) {
	results := []Result{
		{Duration: 5 * time.Millisecond},
		{Duration: 15 * time.Millisecond},
		{Duration: 50 * time.Millisecond},
	}
	r := New(results)
	var buf bytes.Buffer
	if err := Write(&buf, r, results, FormatHDR); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "HDR Histogram") {
		t.Errorf("expected HDR Histogram header in output")
	}
}
