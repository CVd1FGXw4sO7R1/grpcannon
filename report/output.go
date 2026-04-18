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
	FormatCurve      Format = "curve"
)

var validFormats = []Format{
	FormatText, FormatJSON, FormatCSV, FormatTable, FormatMarkdown,
	FormatPrometheus, FormatHTML, FormatXML, FormatInflux, FormatDotPlot,
	FormatSparkline, FormatHeatmap, FormatTimeline, FormatFlamegraph, FormatCurve,
}

// ParseFormat parses a format string, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	f := Format(strings.ToLower(s))
	for _, v := range validFormats {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("unknown format %q", s)
}

// Write writes the report in the requested format.
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
		WriteFlamegraph(w, r)
	case FormatCurve:
		WritePercentileCurve(w, results, 20)
	default:
		return fmt.Errorf("unknown format %q", format)
	}
	return nil
}
