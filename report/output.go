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
	FormatXML        Format = "xml"
	FormatInflux     Format = "influx"
	FormatDotPlot    Format = "dotplot"
	FormatSparkline  Format = "sparkline"
	FormatHeatmap    Format = "heatmap"
	FormatTimeline   Format = "timeline"
	FormatFlamegraph Format = "flamegraph"
)

// ParseFormat parses a string into a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown,
		FormatPrometheus, FormatHTML, FormatXML, FormatInflux,
		FormatDotPlot, FormatSparkline, FormatHeatmap, FormatTimeline, FormatFlamegraph:
		return Format(strings.ToLower(s)), nil
	default:
		return "", fmt.Errorf("unknown format: %q", s)
	}
}

// Write writes the report in the given format to w.
func Write(w io.Writer, r *Report, results []Result, f Format) error {
	switch f {
	case FormatText:
		return WriteText(w, r)
	case FormatJSON:
		return WriteJSON(w, r)
	case FormatCSV:
		return WriteCSV(w, results)
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
	case FormatDotPlot:
		return WriteDotPlot(w, results)
	case FormatSparkline:
		return WriteSparkline(w, results)
	case FormatHeatmap:
		return WriteHeatmap(w, results, 0)
	case FormatTimeline:
		return WriteTimeline(w, results, 0)
	case FormatFlamegraph:
		return WriteFlamegraph(w, results)
	default:
		return fmt.Errorf("unsupported format: %q", f)
	}
}
