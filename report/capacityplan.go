package report

import (
	"fmt"
	"io"
	"sort"
)

// CapacityPoint represents a projected operating point at a given concurrency level.
type CapacityPoint struct {
	Concurrency int
	ProjectedRPS float64
	ProjectedP99Ms float64
	SaturationPct float64 // 0–100
}

// CapacityPlan holds the full capacity projection derived from a step-load run.
type CapacityPlan struct {
	Points          []CapacityPoint
	RecommendedMax  int     // concurrency at which SLO is still met
	PeakRPS         float64
	SLOThresholdMs  float64
}

// BuildCapacityPlan derives a capacity plan from step-load results.
// sloP99Ms is the maximum acceptable P99 latency in milliseconds.
func BuildCapacityPlan(steps []StepPoint, sloP99Ms float64) *CapacityPlan {
	if len(steps) == 0 {
		return &CapacityPlan{SLOThresholdMs: sloP99Ms}
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Concurrency < steps[j].Concurrency
	})

	plan := &CapacityPlan{
		SLOThresholdMs: sloP99Ms,
		Points:         make([]CapacityPoint, 0, len(steps)),
	}

	var peakRPS float64
	recommended := steps[0].Concurrency

	for _, s := range steps {
		var saturation float64
		if s.P99Ms > 0 && sloP99Ms > 0 {
			saturation = (s.P99Ms / sloP99Ms) * 100
			if saturation > 100 {
				saturation = 100
			}
		}
		pt := CapacityPoint{
			Concurrency:    s.Concurrency,
			ProjectedRPS:   s.RPS,
			ProjectedP99Ms: s.P99Ms,
			SaturationPct:  saturation,
		}
		plan.Points = append(plan.Points, pt)
		if s.RPS > peakRPS {
			peakRPS = s.RPS
		}
		if sloP99Ms <= 0 || s.P99Ms <= sloP99Ms {
			recommended = s.Concurrency
		}
	}

	plan.PeakRPS = peakRPS
	plan.RecommendedMax = recommended
	return plan
}

// WriteCapacityPlan writes a human-readable capacity plan to w.
func WriteCapacityPlan(w io.Writer, plan *CapacityPlan) {
	if plan == nil || len(plan.Points) == 0 {
		fmt.Fprintln(w, "capacity plan: no data")
		return
	}
	fmt.Fprintf(w, "Capacity Plan (SLO P99 ≤ %.1f ms)\n", plan.SLOThresholdMs)
	fmt.Fprintf(w, "  Recommended max concurrency : %d\n", plan.RecommendedMax)
	fmt.Fprintf(w, "  Peak RPS observed            : %.2f\n", plan.PeakRPS)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "  %-12s  %-10s  %-12s  %s\n", "Concurrency", "RPS", "P99 (ms)", "Saturation")
	fmt.Fprintf(w, "  %-12s  %-10s  %-12s  %s\n", "-----------", "---", "--------", "----------")
	for _, p := range plan.Points {
		marker := ""
		if p.Concurrency == plan.RecommendedMax {
			marker = " ◀ recommended"
		}
		fmt.Fprintf(w, "  %-12d  %-10.2f  %-12.2f  %.1f%%%s\n",
			p.Concurrency, p.ProjectedRPS, p.ProjectedP99Ms, p.SaturationPct, marker)
	}
}
