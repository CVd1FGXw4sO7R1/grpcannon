package cmd

import (
	"fmt"
	"os"

	"github.com/yourusername/grpcannon/report"
)

// BaselineFlags holds CLI flags for the baseline subcommand.
type BaselineFlags struct {
	SavePath string
	LoadPath string
}

// RunSaveBaseline captures and saves a baseline from the provided results.
func RunSaveBaseline(results []report.Result, path string) error {
	r := report.New(results)
	snap, err := report.CaptureBaseline(r)
	if err != nil {
		return fmt.Errorf("capture baseline: %w", err)
	}
	if err := report.SaveBaseline(snap, path); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	report.WriteBaseline(os.Stdout, snap)
	fmt.Fprintf(os.Stdout, "Baseline saved to %s\n", path)
	return nil
}

// RunShowBaseline loads and prints a previously saved baseline.
func RunShowBaseline(path string) error {
	snap, err := report.LoadBaseline(path)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	report.WriteBaseline(os.Stdout, snap)
	return nil
}
