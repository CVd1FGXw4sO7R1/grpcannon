package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/your-org/grpcannon/report"
)

// RunPercentileHistory prints a multi-run percentile history table.
// Each entry in labeledResults maps a human-readable label to a *report.Report.
func RunPercentileHistory(labeledResults []struct {
	Label  string
	Report *report.Report
}) error {
	if len(labeledResults) == 0 {
		return fmt.Errorf("percentilehistory: no runs provided")
	}

	h := report.BuildPercentileHistory(labeledResults)
	if len(h.Entries) == 0 {
		fmt.Fprintln(os.Stdout, "no valid history entries to display")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Percentile History (%d runs)\n", len(h.Entries))
	fmt.Fprintf(os.Stdout, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	report.WritePercentileHistory(os.Stdout, h)
	return nil
}
