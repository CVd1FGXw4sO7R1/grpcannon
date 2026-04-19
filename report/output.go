package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents a supported output format.
type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatCSV
	FormatTable
	FormatMarkdown
	FormatPrometheus
	FormatHTML
	FormatXML
	FormatInflux
)

var formatNames = map[string]Format{
	"text":       FormatText,
	"json":       FormatJSON,
	"csv":        FormatCSV,
	"table":      FormatTable,
	"markdown":   FormatMarkdown,
	"prometheus": FormatPrometheus,
	"html":       FormatHTML,
	"xml":        FormatXML,
	"influx":     FormatInflux,
}

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	f, ok := formatNames[strings.ToLower(s)]
	if !ok {
		return FormatText, fmt.Errorf("unknown format %q", s)
	}
	return f, nil
}

// Write writes the report in the given format to w.
func Write(w io.Writer, r *Report, f Format) error {
	if r == nil {
		return fmt.Errorf("nil report")
	}
	switch f {
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
	default:
		return WriteText(w, r)
	}
}
