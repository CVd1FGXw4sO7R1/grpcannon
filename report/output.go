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
	FormatDotPlot    Format = "dotplot"
	FormatSparkline  Format = "sparkline"
	FormatHeatmap    Format = "heatmap"
	FormatTimeline   Format = "timeline"
	FormatFlamegraph Format = "flamegraph"
	FormatLatency    Format = "latency"
)

var validFormats = []Format{
	FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown,
	FormatPrometheus, FormatHTML, FormatXML, FormatInflux, FormatDotPlot,
	FormatSparkline, FormatHeatmap, FormatTimeline, FormatFlamegraph, FormatLatency,
}

// ParseFormat parses a string into a Format, returning an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	f := Format(strings.ToLower(strings.TrimSpace(s)))
	for _, v := range validFormats {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("unknown format %q", s)
}

// Write dispatches to the appropriate writer based on format.
func Write(w io.Writer, r *Report, results []Result, format Format) error {
	switch format {
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
	case FormatLatency:
		return WriteLatencyBands(w, results, 0)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
