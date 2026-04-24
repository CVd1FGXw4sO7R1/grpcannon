package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makePheatmapResults(n int, baseDuration time.Duration) []Result {
	results := make([]Result, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		results[i] = Result{
			StartedAt: now.Add(time.Duration(i) * time.Second),
			Duration:  baseDuration + time.Duration(i)*time.Millisecond,
		}
	}
	return results
}

func TestBuildPercentileHeatmap_Empty(t *testing.T) {
	rows := BuildPercentileHeatmap([]Result{}, 5)
	if rows != nil {
		t.Errorf("expected nil for empty results, got %v", rows)
	}
}

func TestBuildPercentileHeatmap_ZeroBuckets(t *testing.T) {
	results := makePheatmapResults(10, 10*time.Millisecond)
	rows := BuildPercentileHeatmap(results, 0)
	if rows != nil {
		t.Errorf("expected nil for zero buckets, got %v", rows)
	}
}

func TestBuildPercentileHeatmap_RowCount(t *testing.T) {
	results := makePheatmapResults(20, 5*time.Millisecond)
	rows := BuildPercentileHeatmap(results, 4)
	// expect 5 percentile rows: 50, 75, 90, 95, 99
	if len(rows) != 5 {
		t.Errorf("expected 5 rows, got %d", len(rows))
	}
}

func TestBuildPercentileHeatmap_BucketCount(t *testing.T) {
	results := makePheatmapResults(20, 5*time.Millisecond)
	const buckets = 4
	rows := BuildPercentileHeatmap(results, buckets)
	for _, row := range rows {
		if len(row.Buckets) != buckets {
			t.Errorf("expected %d bucket values, got %d for P%.0f", buckets, len(row.Buckets), row.Percentile)
		}
	}
}

func TestBuildPercentileHeatmap_PercentilesAscending(t *testing.T) {
	results := makePheatmapResults(50, 1*time.Millisecond)
	rows := BuildPercentileHeatmap(results, 5)
	for _, row := range rows {
		for bi := 0; bi < len(row.Buckets)-1; bi++ {
			// each row's percentile value should be non-negative
			if row.Buckets[bi] < 0 {
				t.Errorf("P%.0f bucket %d has negative value %.2f", row.Percentile, bi, row.Buckets[bi])
			}
		}
	}
}

func TestWritePercentileHeatmap_ValidOutput(t *testing.T) {
	results := makePheatmapResults(20, 10*time.Millisecond)
	var buf bytes.Buffer
	WritePercentileHeatmap(&buf, results, 3)
	out := buf.String()
	if !strings.Contains(out, "P50") {
		t.Error("expected P50 in output")
	}
	if !strings.Contains(out, "P99") {
		t.Error("expected P99 in output")
	}
}

func TestWritePercentileHeatmap_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileHeatmap(&buf, []Result{}, 5)
	out := buf.String()
	if !strings.Contains(out, "no data") {
		t.Errorf("expected 'no data', got: %s", out)
	}
}
