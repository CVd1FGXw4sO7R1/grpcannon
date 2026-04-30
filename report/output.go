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
)

// ParseFormat converts a string to a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatText, FormatJSON, FormatCSV, FormatTable,
		FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML, FormatInflux:
		return Format(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unknown format: %q", s)
}

// Write dispatches the report to w in the requested format.
func Write(w io.Writer, rep *Report, f Format) error {
	switch f {
	case FormatText:
		return WriteText(w, rep)
	case FormatJSON:
		return WriteJSON(w, rep)
	case FormatCSV:
		return WriteCSV(w, rep)
	case FormatTable:
		WriteTable(w, rep)
		return nil
	case FormatMarkdown:
		WriteMarkdown(w, rep)
		return nil
	case FormatPrometheus:
		WritePrometheus(w, rep)
		return nil
	case FormatHTML:
		WriteHTML(w, rep)
		return nil
	case FormatXML:
		WriteXML(w, rep)
		return nil
	case FormatInflux:
		WriteInflux(w, rep)
		return nil
	default:
		return fmt.Errorf("unsupported format: %q", f)
	}
}
