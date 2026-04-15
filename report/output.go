package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents an output format for the report.
type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatCSV
)

// ParseFormat parses a format string into a Format constant.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	default:
		return FormatText, fmt.Errorf("unknown format %q: must be one of text, json, csv", s)
	}
}

// Write writes the report r to w using the given format.
func Write(r *Report, format Format, w io.Writer) error {
	switch format {
	case FormatJSON:
		return WriteJSON(r, w)
	case FormatCSV:
		return WriteCSV(r, w)
	default:
		return WriteText(r, w)
	}
}
