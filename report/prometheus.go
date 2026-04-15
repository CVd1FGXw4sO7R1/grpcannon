package report

import (
	"fmt"
	"io"
	"time"
)

// WritePrometheus writes the report in Prometheus exposition format.
func WritePrometheus(w io.Writer, r *Report) error {
	if r == nil {
		return fmt.Errorf("report is nil")
	}

	totalDuration := r.TotalDuration.Seconds()

	fmt.Fprintf(w, "# HELP grpcannon_requests_total Total number of gRPC requests made.\n")
	fmt.Fprintf(w, "# TYPE grpcannon_requests_total counter\n")
	fmt.Fprintf(w, "grpcannon_requests_total %d\n", r.Total)

	fmt.Fprintf(w, "# HELP grpcannon_requests_success Total number of successful gRPC requests.\n")
	fmt.Fprintf(w, "# TYPE grpcannon_requests_success counter\n")
	fmt.Fprintf(w, "grpcannon_requests_success %d\n", r.Success)

	fmt.Fprintf(w, "# HELP grpcannon_requests_failed Total number of failed gRPC requests.\n")
	fmt.Fprintf(w, "# TYPE grpcannon_requests_failed counter\n")
	fmt.Fprintf(w, "grpcannon_requests_failed %d\n", r.Failures)

	fmt.Fprintf(w, "# HELP grpcannon_duration_seconds Total wall-clock duration of the load test.\n")
	fmt.Fprintf(w, "# TYPE grpcannon_duration_seconds gauge\n")
	fmt.Fprintf(w, "grpcannon_duration_seconds %g\n", totalDuration)

	fmt.Fprintf(w, "# HELP grpcannon_rps Requests per second throughput.\n")
	fmt.Fprintf(w, "# TYPE grpcannon_rps gauge\n")
	fmt.Fprintf(w, "grpcannon_rps %g\n", r.RPS)

	percentiles := []struct {
		q     float64
		value time.Duration
	}{
		{0.50, r.Fastest},
		{0.50, r.P50},
		{0.90, r.P90},
		{0.95, r.P95},
		{0.99, r.P99},
	}

	fmt.Fprintf(w, "# HELP grpcannon_latency_seconds gRPC request latency percentiles in seconds.\n")
	fmt.Fprintf(w, "# TYPE grpcannon_latency_seconds summary\n")
	for _, p := range percentiles[1:] {
		fmt.Fprintf(w, "grpcannon_latency_seconds{quantile=\"%g\"} %g\n", p.q, p.value.Seconds())
	}
	fmt.Fprintf(w, "grpcannon_latency_seconds_fastest %g\n", r.Fastest.Seconds())
	fmt.Fprintf(w, "grpcannon_latency_seconds_slowest %g\n", r.Slowest.Seconds())
	fmt.Fprintf(w, "grpcannon_latency_seconds_mean %g\n", r.Mean.Seconds())

	_ = percentiles[0]
	return nil
}
