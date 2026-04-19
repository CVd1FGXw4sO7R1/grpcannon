package report

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeBaselineReport() *Report {
	return New([]Result{
		{Duration: ms(10), Error: nil},
		{Duration: ms(20), Error: nil},
		{Duration: ms(30), Error: nil},
		{Duration: ms(40), Error: nil},
		{Duration: ms(200), Error: nil},
	})
}

func TestCaptureBaseline_Nil(t *testing.T) {
	_, err := CaptureBaseline(nil)
	if err == nil {
		t.Fatal("expected error for nil report")
	}
}

func TestCaptureBaseline_Fields(t *testing.T) {
	r := makeBaselineReport()
	snap, err := CaptureBaseline(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.P50Ms <= 0 {
		t.Error("expected positive P50")
	}
	if snap.P99Ms < snap.P50Ms {
		t.Error("P99 should be >= P50")
	}
	if snap.SuccessRate != 100.0 {
		t.Errorf("expected 100%% success, got %.1f", snap.SuccessRate)
	}
	if snap.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	snap := &BaselineSnapshot{
		Timestamp:   time.Now().UTC(),
		P50Ms:       15.0,
		P95Ms:       40.0,
		P99Ms:       80.0,
		AvgMs:       20.0,
		SuccessRate: 99.5,
		RPS:         120.0,
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := SaveBaseline(snap, path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.P99Ms != snap.P99Ms {
		t.Errorf("P99 mismatch: got %.2f want %.2f", loaded.P99Ms, snap.P99Ms)
	}
	if loaded.RPS != snap.RPS {
		t.Errorf("RPS mismatch")
	}
}

func TestLoadBaseline_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadBaseline(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSaveBaseline_IsValidJSON(t *testing.T) {
	snap := &BaselineSnapshot{P50Ms: 5, P99Ms: 50, SuccessRate: 98}
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	SaveBaseline(snap, path)
	data, _ := os.ReadFile(path)
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Errorf("not valid JSON: %v", err)
	}
}

func TestWriteBaseline_Output(t *testing.T) {
	snap := &BaselineSnapshot{P50Ms: 10, P95Ms: 20, P99Ms: 30, AvgMs: 15, SuccessRate: 100, RPS: 50}
	var buf bytes.Buffer
	WriteBaseline(&buf, snap)
	out := buf.String()
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestWriteBaseline_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteBaseline(&buf, nil)
	if buf.String() != "no baseline\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
