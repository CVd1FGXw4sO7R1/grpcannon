package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// WriteCSV writes the report results as CSV to the provided writer.
func WriteCSV(w io.Writer, r *Report) error {
	cw := csv.NewWriter(w)

	header := []string{"status", "duration_ms", "error"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("writing csv header: %w", err)
	}

	for _, res := range r.Results {
		errStr := ""
		if res.Error != nil {
			errStr = res.Error.Error()
		}
		row := []string{
			res.Status,
			strconv.FormatFloat(msFloat(res.Duration), 'f', 3, 64),
			errStr,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("writing csv row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
