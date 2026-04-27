package cmd

import (
	"os"

	"github.com/yourorg/grpcannon/report"
)

// DefaultPercentileTrendBuckets is the default number of time buckets.
const DefaultPercentileTrendBuckets = 10

// RunPercentileTrend builds and prints a percentile trend report from results.
func RunPercentileTrend(results []report.Result, buckets int) {
	if buckets <= 0 {
		buckets = DefaultPercentileTrendBuckets
	}
	points := report.BuildPercentileTrend(results, buckets)
	report.WritePercentileTrend(os.Stdout, points)
}
