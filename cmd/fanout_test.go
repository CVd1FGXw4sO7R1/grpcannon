package cmd

import (
	"testing"
)

func TestParseLevels_Valid(t *testing.T) {
	levels, err := parseLevels("1,2,4,8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(levels) != 4 {
		t.Errorf("expected 4 levels, got %d", len(levels))
	}
	expected := []int{1, 2, 4, 8}
	for i, v := range expected {
		if levels[i] != v {
			t.Errorf("levels[%d]: expected %d, got %d", i, v, levels[i])
		}
	}
}

func TestParseLevels_WithSpaces(t *testing.T) {
	levels, err := parseLevels(" 1 , 4 , 16 ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(levels) != 3 {
		t.Errorf("expected 3 levels, got %d", len(levels))
	}
}

func TestParseLevels_Invalid(t *testing.T) {
	_, err := parseLevels("1,abc,4")
	if err == nil {
		t.Error("expected error for non-numeric level")
	}
}

func TestParseLevels_Empty(t *testing.T) {
	levels, err := parseLevels("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(levels) != 0 {
		t.Errorf("expected 0 levels, got %d", len(levels))
	}
}

func TestDefaultFanoutLevels_Count(t *testing.T) {
	if len(DefaultFanoutLevels) == 0 {
		t.Error("DefaultFanoutLevels should not be empty")
	}
}

func TestDefaultFanoutLevels_Ascending(t *testing.T) {
	for i := 1; i < len(DefaultFanoutLevels); i++ {
		if DefaultFanoutLevels[i] <= DefaultFanoutLevels[i-1] {
			t.Errorf("DefaultFanoutLevels not ascending at index %d", i)
		}
	}
}
