package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents a supported output format.
type Format string

const (
	FormatText       Format = "text"
	FormatJSON       Format = "json"
	FormatCSV        Format = "csv"
	FormatTable      Format = "table"
	FormatMarkdown   Format = "markdown"
	FormatPrometheus Format = "prometheus"
	FormatHTML       Format = "html"
	FormatXML        Format = "xml"
	FormatInflux     Format = "influx"
	FormatBaseline   Format = "baseline"
)

// ParseFormat converts a string to a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatText, FormatJSON, FormatCSV, FormatTable,
		FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML,
		FormatInflux, FormatBaseline:
		return Format(strings.ToLower(s)), nil
	default:
		return "", fmt.Errorf("unknown format: %q", s)
	}
}

// Write renders the report in the requested format to w.
func Write(w io.Writer, r *Report, f Format) error {
	switch f {
	case FormatText:
		return WriteText(w, r)
	case FormatJSON:
		return WriteJSON(w, r)
	case FormatCSV:
		return WriteCSV(w, r)
	case FormatTable:
		return WriteTable(w, r)
	case FormatMarkdown:
		return WriteMarkdown(w, r)
	case FormatPrometheus:
		return WritePrometheus(w, r)
	case FormatHTML:
		return WriteHTML(w, r)
	case FormatXML:
		return WriteXML(w, r)
	case FormatInflux:
		return WriteInflux(w, r)
	case FormatBaseline:
		snap, err := CaptureBaseline(r)
		if err != nil {
			return err
		}
		WriteBaseline(w, snap)
		return nil
	default:
		return fmt.Errorf("unsupported format: %q", f)
	}
}
