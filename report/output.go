package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatTable Format = "table"
)

// ParseFormat parses a string into a Format, returning an error if unrecognised.
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
	default:
		return "", fmt.Errorf("unknown format %q: must be one of text, json, csv, table", s)
	}
}

// Write writes the report r to w using the specified format.
func Write(r *Report, format Format, w io.Writer) error {
	switch format {
	case FormatText:
		return WriteText(r, w)
	case FormatJSON:
		return WriteJSON(r, w)
	case FormatCSV:
		return WriteCSV(r, w)
	case FormatTable:
		return WriteTable(r, w)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
