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

var validFormats = []Format{
	FormatText, FormatJSON, FormatCSV, FormatTable,
	FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML, FormatInflux,
}

// ParseFormat parses a string into a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	f := Format(strings.ToLower(strings.TrimSpace(s)))
	for _, v := range validFormats {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("unknown format %q", s)
}

// Write writes the report in the given format to w.
func Write(w io.Writer, r *Report, f Format) error {
	switch f {
	case FormatJSON:
		return WriteJSON(w, r)
	case FormatCSV:
		return WriteCSV(w, r)
	case FormatTable:
		WriteTable(w, r)
		return nil
	case FormatMarkdown:
		WriteMarkdown(w, r)
		return nil
	case FormatPrometheus:
		WritePrometheus(w, r)
		return nil
	case FormatHTML:
		WriteHTML(w, r)
		return nil
	case FormatXML:
		return WriteXML(w, r)
	case FormatInflux:
		WriteInflux(w, r)
		return nil
	case FormatText:
		WriteText(w, r)
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", f)
	}
}
