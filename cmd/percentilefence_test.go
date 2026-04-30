package cmd

import "testing"

func TestParseFenceThresholds_Valid(t *testing.T) {
	tm, err := parseFenceThresholds("50=20,90=50,99=200")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tm) != 3 {
		t.Errorf("expected 3 entries, got %d", len(tm))
	}
	if tm[99] != 200 {
		t.Errorf("expected tm[99]=200, got %f", tm[99])
	}
}

func TestParseFenceThresholds_Empty(t *testing.T) {
	tm, err := parseFenceThresholds("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tm) != 0 {
		t.Errorf("expected empty map, got %d entries", len(tm))
	}
}

func TestParseFenceThresholds_Invalid(t *testing.T) {
	_, err := parseFenceThresholds("badvalue")
	if err == nil {
		t.Error("expected error for bad input")
	}
}

func TestParseFenceThresholds_BadPercentile(t *testing.T) {
	_, err := parseFenceThresholds("abc=100")
	if err == nil {
		t.Error("expected error for non-numeric percentile")
	}
}

func TestParseFenceThresholds_BadMaxMs(t *testing.T) {
	_, err := parseFenceThresholds("99=abc")
	if err == nil {
		t.Error("expected error for non-numeric max_ms")
	}
}

func TestParseFenceThresholds_WithSpaces(t *testing.T) {
	tm, err := parseFenceThresholds(" 50 = 20 , 99 = 100 ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tm) != 2 {
		t.Errorf("expected 2 entries, got %d", len(tm))
	}
}
