package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case string(FormatText), "":
		return FormatText, nil
	case string(FormatJSON):
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of [text, json]", s)
	}
}

// Write writes the report to w using the specified format.
func Write(r *Report, format Format, w io.Writer) error {
	switch format {
	case FormatJSON:
		return WriteJSON(r, w)
	case FormatText:
		return WriteText(r, w)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
