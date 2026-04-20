package report

import (
	"fmt"
	"io"
	"strings"
)

// Format represents an output format for reports.
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
	FormatDotPlot
	FormatSparkline
	FormatHeatmap
	FormatTimeline
	FormatFlamegraph
	FormatHDR
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
	"dotplot":    FormatDotPlot,
	"sparkline":  FormatSparkline,
	"heatmap":    FormatHeatmap,
	"timeline":   FormatTimeline,
	"flamegraph": FormatFlamegraph,
	"hdr":        FormatHDR,
}

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	f, ok := formatNames[strings.ToLower(s)]
	if !ok {
		return FormatText, fmt.Errorf("unknown format: %q", s)
	}
	return f, nil
}

// Write writes the report in the given format to w.
func Write(w io.Writer, r *Report, results []Result, format Format) error {
	switch format {
	case FormatText:
		WriteText(w, r)
	case FormatJSON:
		return WriteJSON(w, r)
	case FormatCSV:
		WriteCSV(w, results)
	case FormatTable:
		WriteTable(w, r)
	case FormatMarkdown:
		WriteMarkdown(w, r)
	case FormatPrometheus:
		WritePrometheus(w, r)
	case FormatHTML:
		WriteHTML(w, r)
	case FormatXML:
		return WriteXML(w, r)
	case FormatInflux:
		WriteInflux(w, r)
	case FormatDotPlot:
		WriteDotPlot(w, results)
	case FormatSparkline:
		WriteSparkline(w, results)
	case FormatHeatmap:
		WriteHeatmap(w, results, 0)
	case FormatTimeline:
		WriteTimeline(w, results, 0)
	case FormatFlamegraph:
		WriteFlamegraph(w, results)
	case FormatHDR:
		WriteHDRHistogram(w, results, 10)
	default:
		return fmt.Errorf("unsupported format: %d", format)
	}
	return nil
}
