package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"grpcannon/report"
)

var deadLetterN int

var deadLetterCmd = &cobra.Command{
	Use:   "dead-letter",
	Short: "Show the top-N slowest failed requests from the last run",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunDeadLetter(deadLetterN)
	},
}

func init() {
	deadLetterCmd.Flags().IntVarP(&deadLetterN, "top", "n", 10, "Number of failed requests to display")
	rootCmd.AddCommand(deadLetterCmd)
}

// RunDeadLetter builds and prints a dead-letter queue from synthetic results.
func RunDeadLetter(n int) error {
	results := syntheticDeadLetterResults()
	dlq := report.BuildDeadLetterQueue(results, n)
	report.WriteDeadLetterQueue(os.Stdout, dlq)
	return nil
}

// syntheticDeadLetterResults generates a small set of results for demo purposes.
func syntheticDeadLetterResults() []report.Result {
	return []report.Result{
		{Duration: 12 * time.Millisecond, Err: nil},
		{Duration: 340 * time.Millisecond, Err: fmt.Errorf("rpc error: code = Unavailable")},
		{Duration: 5 * time.Millisecond, Err: nil},
		{Duration: 820 * time.Millisecond, Err: fmt.Errorf("rpc error: code = DeadlineExceeded")},
		{Duration: 15 * time.Millisecond, Err: nil},
		{Duration: 90 * time.Millisecond, Err: fmt.Errorf("rpc error: code = ResourceExhausted")},
	}
}
