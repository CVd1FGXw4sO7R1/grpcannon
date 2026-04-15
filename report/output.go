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
)

// ParseFormat parses a string into a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatText:
		return FormatText, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatCSV:
		return FormatCSV, nil
	case FormatTable:
		return FormatTable, nil
	case FormatMarkdown:
		return FormatMarkdown, nil
	case FormatPrometheus:
		return FormatPrometheus, nil
	default:
		return "", fmt.Errorf("unknown format %q: supported formats are text, json, csv, table, markdown, prometheus", s)
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
	default:
		return fmt.Errorf("unsupported format: %s", f)
	}
}
