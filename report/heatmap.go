package report

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteHeatmap writes a time-bucket heatmap of latencies to w.
// Buckets are columns (time windows), rows are latency bands.
func WriteHeatmap(w io.Writer, results []Result, buckets int) error {
	if results == nil {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no results")
		return err
	}
	if buckets <= 0 {
		buckets = 10
	}

	// latency bands in ms
	bands := []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000}
	bandLabels := []string{"<1ms", "<5ms", "<10ms", "<25ms", "<50ms", "<100ms", "<250ms", "<500ms", "<1s", ">=1s"}

	// split results into time buckets
	start := results[0].Start
	end := results[len(results)-1].Start
	span := end.Sub(start)
	if span == 0 {
		span = time.Millisecond
	}
	bucketSize := span / time.Duration(buckets)

	grid := make([][]int, len(bandLabels))
	for i := range grid {
		grid[i] = make([]int, buckets)
	}

	for _, r := range results {
		bi := int(r.Start.Sub(start) / bucketSize)
		if bi >= buckets {
			bi = buckets - 1
		}
		ms := float64(r.Duration) / float64(time.Millisecond)
		row := len(bands)
		for j, b := range bands {
			if ms < b {
				row = j
				break
			}
		}
		grid[row][bi]++
	}

	shades := []string{" ", "░", "▒", "▓", "█"}
	max := 1
	for _, row := range grid {
		for _, v := range row {
			if v > max {
				max = v
			}
		}
	}

	fmt.Fprintln(w, "Latency Heatmap")
	fmt.Fprintln(w, strings.Repeat("-", buckets+14))
	for i := len(grid) - 1; i >= 0; i-- {
		fmt.Fprintf(w, "%-8s │", bandLabels[i])
		for _, v := range grid[i] {
			idx := v * (len(shades) - 1) / max
			fmt.Fprint(w, shades[idx])
		}
		fmt.Fprintln(w)
	}
	return nil
}
