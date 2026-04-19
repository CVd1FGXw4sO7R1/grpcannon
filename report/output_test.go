package report

import (
	"bytes"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	formats := []string{"text", "json", "csv", "table", "markdown", "prometheus", "html", "xml", "influx", "baseline"}
	for _, s := range formats {
		if _, err := ParseFormat(s); err != nil {
			t.Errorf("expected %q to be valid: %v", s, err)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("nope")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestParseFormat_RoundTrip(t *testing.T) {
	formats := []Format{FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML, FormatInflux, FormatBaseline}
	for _, f := range formats {
		parsed, err := ParseFormat(string(f))
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", f, err)
		}
		if parsed != f {
			t.Errorf("ParseFormat(%q) = %q, want %q", f, parsed, f)
		}
	}
}

func TestWrite_Formats(t *testing.T) {
	r := New([]Result{
		{Duration: ms(10), Error: nil},
		{Duration: ms(50), Error: nil},
	})
	formats := []Format{FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML, FormatInflux, FormatBaseline}
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
	err := Write(&buf, r, Format("unknown"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_BaselineNilReport(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, nil, FormatBaseline)
	if err == nil {
		t.Error("expected error writing baseline from nil report")
	}
}
