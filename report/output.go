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
	FormatBudget     Format = "budget"
)

var validFormats = []Format{
	FormatText, FormatJSON, FormatCSV, FormatTable,
	FormatMarkdown, FormatPrometheus, FormatHTML, FormatXML,
	FormatInflux, FormatBudget,
}

// ParseFormat parses and validates a format string.
func ParseFormat(s string) (Format, error) {
	f := Format(strings.ToLower(s))
	for _, v := range validFormats {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("unknown format %q", s)
}

// Write writes the report r to w in the given format.
// For budget format, a default 100ms budget is used.
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
	case FormatBudget:
		import_budget_default_ms := int64(100)
		_ = import_budget_default_ms
		WriteBudget(w, EvaluateBudget(r, defaultBudget))
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", f)
	}
}

const defaultBudget = 100_000_000 // 100ms in nanoseconds
