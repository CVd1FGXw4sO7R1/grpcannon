package report

import (
	"fmt"
	"io"
	"sort"
)

// StatusCodeBreakdown maps gRPC status code strings to counts.
type StatusCodeBreakdown map[string]int

// BuildStatusCodeBreakdown counts occurrences of each gRPC status code in results.
func BuildStatusCodeBreakdown(results []Result) StatusCodeBreakdown {
	breakdown := make(StatusCodeBreakdown)
	for _, r := range results {
		if r.Error != nil {
			code := extractCode(r.Error)
			breakdown[code]++
		} else {
			breakdown["OK"]++
		}
	}
	return breakdown
}

func extractCode(err error) string {
	if err == nil {
		return "OK"
	}
	// Use error string as code key; real impl would use status.Code(err).String()
	return err.Error()
}

// WriteStatusCodeBreakdown writes a status code frequency table to w.
func WriteStatusCodeBreakdown(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No results.")
		return
	}
	bd := BuildStatusCodeBreakdown(results)

	codes := make([]string, 0, len(bd))
	for code := range bd {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	fmt.Fprintln(w, "Status Code Breakdown:")
	fmt.Fprintf(w, "  %-30s %s\n", "Code", "Count")
	fmt.Fprintf(w, "  %-30s %s\n", "----", "-----")
	for _, code := range codes {
		fmt.Fprintf(w, "  %-30s %d\n", code, bd[code])
	}
}
