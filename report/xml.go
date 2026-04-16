package report

import (
	"encoding/xml"
	"io"
	"time"
)

type xmlReport struct {
	XMLName    xml.Name     `xml:"report"`
	Total      int          `xml:"total"`
	Successes  int          `xml:"successes"`
	Failures   int          `xml:"failures"`
	SuccessRate float64     `xml:"success_rate"`
	Duration   xmlDuration  `xml:"duration"`
	Latency    xmlLatency   `xml:"latency"`
}

type xmlDuration struct {
	Total string `xml:"total"`
}

type xmlLatency struct {
	P50  string `xml:"p50"`
	P90  string `xml:"p90"`
	P95  string `xml:"p95"`
	P99  string `xml:"p99"`
	Mean string `xml:"mean"`
}

func WriteXML(w io.Writer, r *Report) error {
	if r == nil {
		_, err := w.Write([]byte("<report/>\n"))
		return err
	}

	meanDur := time.Duration(0)
	if r.Total > 0 {
		meanDur = r.TotalDuration / time.Duration(r.Total)
	}

	xr := xmlReport{
		Total:       r.Total,
		Successes:   r.Successes,
		Failures:    r.Failures,
		SuccessRate: r.SuccessRate,
		Duration:    xmlDuration{Total: roundDuration(r.TotalDuration).String()},
		Latency: xmlLatency{
			P50:  roundDuration(r.P50).String(),
			P90:  roundDuration(r.P90).String(),
			P95:  roundDuration(r.P95).String(),
			P99:  roundDuration(r.P99).String(),
			Mean: roundDuration(meanDur).String(),
		},
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(xr); err != nil {
		return err
	}
	return nil
}
