package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/grpcannon/report"
)

// RunScore loads a baseline and prints the composite performance score.
func RunScore(baselinePath string, jsonOut bool) error {
	bl, err := report.LoadBaseline(baselinePath)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}

	// Reconstruct a minimal Report from baseline snapshot.
	r := &report.Report{
		Results: bl.Results,
	}

	sc := report.CalcScore(r)

	if jsonOut {
		return json.NewEncoder(os.Stdout).Encode(sc)
	}

	report.WriteScore(os.Stdout, r)
	return nil
}
