package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// WriteCSV writes the report in CSV format to the given writer.
func WriteCSV(r *Report, w io.Writer) error {
	cw := csv.NewWriter(w)

	header := []string{
		"total", "successes", "failures", "success_rate",
		"rps", "min_ms", "mean_ms", "p50_ms", "p95_ms", "p99_ms", "max_ms",
	}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("csv: write header: %w", err)
	}

	var rps float64
	if r.Duration > 0 {
		rps = float64(r.Total) / r.Duration.Seconds()
	}

	row := []string{
		strconv.Itoa(r.Total),
		strconv.Itoa(r.Successes),
		strconv.Itoa(r.Failures),
		strconv.FormatFloat(r.SuccessRate()*100, 'f', 2, 64),
		strconv.FormatFloat(rps, 'f', 2, 64),
		strconv.FormatFloat(msFloat(r.Min), 'f', 3, 64),
		strconv.FormatFloat(msFloat(r.Mean), 'f', 3, 64),
		strconv.FormatFloat(msFloat(r.Percentile(50)), 'f', 3, 64),
		strconv.FormatFloat(msFloat(r.Percentile(95)), 'f', 3, 64),
		strconv.FormatFloat(msFloat(r.Percentile(99)), 'f', 3, 64),
		strconv.FormatFloat(msFloat(r.Max), 'f', 3, 64),
	}
	if err := cw.Write(row); err != nil {
		return fmt.Errorf("csv: write row: %w", err)
	}

	cw.Flush()
	return cw.Error()
}
