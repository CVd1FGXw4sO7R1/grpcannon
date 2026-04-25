package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents a supported output format for reports.
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
	FormatPercentileStats
)

// ParseFormat converts a string format name to a Format constant.
// Returns an error for unknown format names.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	case "table":
		return FormatTable, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	case "prometheus", "prom":
		return FormatPrometheus, nil
	case "html":
		return FormatHTML, nil
	case "xml":
		return FormatXML, nil
	case "influx":
		return FormatInflux, nil
	case "percentilesats", "pstats":
		return FormatPercentileStats, nil
	}
	return 0, fmt.Errorf("unknown format: %q", s)
}

// Write renders the report r in the given format to w.
// Returns an error for unsupported formats.
func Write(w io.Writer, r *Report, results []Result, f Format) error {
	switch f {
	case FormatText:
		WriteText(w, r)
	case FormatJSON:
		WriteJSON(w, r)
	case FormatCSV:
		WriteCSV(w, r)
	case FormatTable:
		WriteTable(w, r)
	case FormatMarkdown:
		WriteMarkdown(w, r)
	case FormatPrometheus:
		WritePrometheus(w, r)
	case FormatHTML:
		WriteHTML(w, r)
	case FormatXML:
		WriteXML(w, r)
	case FormatInflux:
		WriteInflux(w, r)
	case FormatPercentileStats:
		WritePercentileStats(w, results)
	default:
		return fmt.Errorf("unsupported format: %d", f)
	}
	return nil
}
