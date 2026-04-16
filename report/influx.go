package report

import (
	"fmt"
	"io"
	"time"
)

// WriteInflux writes results in InfluxDB line protocol format.
func WriteInflux(w io.Writer, r *Report) error {
	if r == nil {
		return fmt.Errorf("report is nil")
	}

	now := time.Now().UnixNano()

	fmt.Fprintf(w, "grpcannon_summary total=%di,success=%di,failures=%di %d\n",
		r.Total, r.Success, r.Failures, now)

	if r.Total > 0 {
		fmt.Fprintf(w, "grpcannon_latency p50=%.4f,p90=%.4f,p95=%.4f,p99=%.4f,mean=%.4f %d\n",
			r.P50, r.P90, r.P95, r.P99, r.Mean, now)
	}

	return nil
}
