package report

import (
	"bytes"
	"testing"
	"time"
)

func TestParseFormat_Valid(t *testing.T) {
	formats := []string{"text", "json", "csv", "table", "markdown", "prometheus", "html", "xml", "influx"}
	for _, f := range formats {
		_, err := ParseFormat(f)
		if err != nil {
			t.Errorf("expected format %q to be valid, got error: %v", f, err)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("nope")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_Formats(t *testing.T) {
	results := makeResults(4, 1, 20*time.Millisecond)
	r := New(results)
	formats := []Format{FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML, FormatInflux}
	for _, f := range formats {
		var buf bytes.Buffer
		err := Write(&buf, r, f)
		if err != nil {
			t.Errorf("Write(%q) unexpected error: %v", f, err)
		}
		if buf.Len() == 0 {
			t.Errorf("Write(%q) produced no output", f)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, &Report{}, Format("unknown"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
