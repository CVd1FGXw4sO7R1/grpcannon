package report

import (
	"fmt"
	"io"
)

// Grade represents a letter grade for a load test run.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// GradeResult holds the computed grade and the reasons behind it.
type GradeResult struct {
	Grade   Grade
	Score   int // 0-100
	Reasons []string
}

// AssignGrade scores a report and returns a GradeResult.
// Scoring:
//   - Success rate >= 99.9% → 40 pts
//   - Success rate >= 99%   → 30 pts
//   - Success rate >= 95%   → 20 pts
//   - P99 < 100ms           → 30 pts
//   - P99 < 500ms           → 20 pts
//   - P99 < 1000ms          → 10 pts
//   - P50 < 50ms            → 30 pts
//   - P50 < 200ms           → 20 pts
//   - P50 < 500ms           → 10 pts
func AssignGrade(r *Report) GradeResult {
	if r == nil || len(r.Results) == 0 {
		return GradeResult{Grade: GradeF, Score: 0, Reasons: []string{"no results"}}
	}

	var score int
	var reasons []string

	succRate := float64(r.SuccessCount) / float64(r.Total) * 100
	switch {
	case succRate >= 99.9:
		score += 40
		reasons = append(reasons, fmt.Sprintf("success rate %.2f%% (+40)", succRate))
	case succRate >= 99:
		score += 30
		reasons = append(reasons, fmt.Sprintf("success rate %.2f%% (+30)", succRate))
	case succRate >= 95:
		score += 20
		reasons = append(reasons, fmt.Sprintf("success rate %.2f%% (+20)", succRate))
	default:
		reasons = append(reasons, fmt.Sprintf("success rate %.2f%% (+0)", succRate))
	}

	p99 := Percentile(r.Results, 99)
	switch {
	case p99 < 100:
		score += 30
		reasons = append(reasons, fmt.Sprintf("p99 %.1fms (+30)", p99))
	case p99 < 500:
		score += 20
		reasons = append(reasons, fmt.Sprintf("p99 %.1fms (+20)", p99))
	case p99 < 1000:
		score += 10
		reasons = append(reasons, fmt.Sprintf("p99 %.1fms (+10)", p99))
	default:
		reasons = append(reasons, fmt.Sprintf("p99 %.1fms (+0)", p99))
	}

	p50 := Percentile(r.Results, 50)
	switch {
	case p50 < 50:
		score += 30
		reasons = append(reasons, fmt.Sprintf("p50 %.1fms (+30)", p50))
	case p50 < 200:
		score += 20
		reasons = append(reasons, fmt.Sprintf("p50 %.1fms (+20)", p50))
	case p50 < 500:
		score += 10
		reasons = append(reasons, fmt.Sprintf("p50 %.1fms (+10)", p50))
	default:
		reasons = append(reasons, fmt.Sprintf("p50 %.1fms (+0)", p50))
	}

	var grade Grade
	switch {
	case score >= 90:
		grade = GradeA
	case score >= 75:
		grade = GradeB
	case score >= 60:
		grade = GradeC
	case score >= 40:
		grade = GradeD
	default:
		grade = GradeF
	}

	return GradeResult{Grade: grade, Score: score, Reasons: reasons}
}

// WriteGrade writes a human-readable grade report to w.
func WriteGrade(w io.Writer, r *Report) {
	gr := AssignGrade(r)
	fmt.Fprintf(w, "Grade: %s (score: %d/100)\n", gr.Grade, gr.Score)
	for _, reason := range gr.Reasons {
		fmt.Fprintf(w, "  • %s\n", reason)
	}
}
