package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/yourorg/grpcannon/report"
)

// DefaultSaturationSteps is the default number of concurrency steps to test.
const DefaultSaturationSteps = 5

// RunSaturation executes a saturation sweep across concurrency levels and
// writes the analysis to stdout. groups maps concurrency level to the slice of
// results collected at that level; window is the measurement duration used to
// compute RPS.
func RunSaturation(groups map[int][]report.Result, window time.Duration) error {
	if len(groups) == 0 {
		return fmt.Errorf("saturation: no result groups provided")
	}
	sr := report.BuildSaturation(groups, window)
	report.WriteSaturation(os.Stdout, sr)
	return nil
}
