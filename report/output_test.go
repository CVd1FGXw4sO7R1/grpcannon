package report

import (
	"bytes"
	"testing"
	"time"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []string{"text", "json", "csv", "table", "markdown", "prometheus", "html", "xml", "influx"}
	for _, c := range cases {
		if _, err := ParseFormat(c); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", c, err)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	if _, err := ParseFormat("nope"); err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestParseFormat_RoundTrip(t *testing.T) {
	f, err := ParseFormat("json")
	if err != nil {
		t.Fatal(err)
	}
	if f != FormatJSON {
		t.Fatalf("expected FormatJSON")
	}
}

func TestWrite_Formats(t *testing.T) {
	r := New([]Result{
		{Duration: 10 * time.Millisecond},
		{Duration: 20 * time.Millisecond},
	})
	formats := []Format{FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML, FormatInflux}
	for _, f := range formats {
		var buf bytes.Buffer
		if err := Write(&buf, r, f); err != nil {
			t.Errorf("Write failed for format %v: %v", f, err)
		}
		if buf.Len() == 0 {
			t.Errorf("expected non-empty output for format %v", f)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	r := New([]Result{{Duration: 5 * time.Millisecond}})
	var buf bytes.Buffer
	if err := Write(&buf, r, Format(999)); err != nil {
		t.Fatalf("unexpected error for unknown format (falls back to text): %v", err)
	}
}

func TestWrite_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, nil, FormatText); err == nil {
		t.Fatal("expected error for nil report")
	}
}
