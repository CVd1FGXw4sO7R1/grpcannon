package report

import (
	"encoding/json"
	"io"
	"time"
)

// JSONSummary is a JSON-serialisable representation of Summary.
type JSONSummary struct {
	Total      int     `json:"total"`
	Successes  int     `json:"successes"`
	Failures   int     `json:"failures"`
	TotalTimeMs float64 `json:"total_time_ms"`
	P50Ms      float64 `json:"p50_ms"`
	P90Ms      float64 `json:"p90_ms"`
	P99Ms      float64 `json:"p99_ms"`
}

func msFloat(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// WriteJSON encodes the summary as JSON into w.
func (s *Summary) WriteJSON(w io.Writer) error {
	js := JSONSummary{
		Total:       s.Total,
		Successes:   s.Successes,
		Failures:    s.Failures,
		TotalTimeMs: msFloat(s.TotalTime),
		P50Ms:       msFloat(s.Percentile(50)),
		P90Ms:       msFloat(s.Percentile(90)),
		P99Ms:       msFloat(s.Percentile(99)),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(js)
}
