package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents an output format for reports.
type Format string

const (
	FormatText       Format = "text"
	FormatJSON       Format = "json"
	FormatCSV        Format = "csv"
	FormatTable      Format = "table"
	FormatMarkdown   Format = "markdown"
	FormatPrometheus Format = "prometheus"
	FormatHTML       Format = "html"
)

// ParseFormat parses a string into a Format.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown, FormatPrometheus, FormatHTML:
		return Format(strings.ToLower(s)), nil
	default:
		return "", fmt.Errorf("unknown format %q: choose one of text, json, csv, table, markdown, prometheus, html", s)
	}
}

// Write writes the report r to w in the given format.
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
	default:
		return fmt.Errorf("unsupported format: %s", f)
	}
}
