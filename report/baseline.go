package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// BaselineSnapshot captures key metrics for future regression comparison.
type BaselineSnapshot struct {
	Timestamp  time.Time `json:"timestamp"`
	P50Ms      float64   `json:"p50_ms"`
	P95Ms      float64   `json:"p95_ms"`
	P99Ms      float64   `json:"p99_ms"`
	AvgMs      float64   `json:"avg_ms"`
	SuccessRate float64  `json:"success_rate"`
	RPS        float64   `json:"rps"`
}

// CaptureBaseline builds a BaselineSnapshot from a Report.
func CaptureBaseline(r *Report) (*BaselineSnapshot, error) {
	if r == nil {
		return nil, fmt.Errorf("report is nil")
	}
	durations := SortedDurationsMs(r.Results)
	total := len(r.Results)
	succ := 0
	for _, res := range r.Results {
		if res.Error == nil {
			succ++
		}
	}
	var successRate float64
	if total > 0 {
		successRate = float64(succ) / float64(total) * 100
	}
	return &BaselineSnapshot{
		Timestamp:   time.Now().UTC(),
		P50Ms:       Percentile(durations, 50),
		P95Ms:       Percentile(durations, 95),
		P99Ms:       Percentile(durations, 99),
		AvgMs:       r.Avg.Seconds() * 1000,
		SuccessRate: successRate,
		RPS:         r.RPS,
	}, nil
}

// SaveBaseline writes a BaselineSnapshot as JSON to path.
func SaveBaseline(snap *BaselineSnapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(snap)
}

// LoadBaseline reads a BaselineSnapshot from a JSON file.
func LoadBaseline(path string) (*BaselineSnapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var snap BaselineSnapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// WriteBaseline prints a BaselineSnapshot summary to w.
func WriteBaseline(w io.Writer, snap *BaselineSnapshot) {
	if snap == nil {
		fmt.Fprintln(w, "no baseline")
		return
	}
	fmt.Fprintf(w, "Baseline captured at %s\n", snap.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  P50: %.2f ms  P95: %.2f ms  P99: %.2f ms\n", snap.P50Ms, snap.P95Ms, snap.P99Ms)
	fmt.Fprintf(w, "  Avg: %.2f ms  Success: %.1f%%  RPS: %.2f\n", snap.AvgMs, snap.SuccessRate, snap.RPS)
}
