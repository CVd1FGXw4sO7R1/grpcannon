package report

import (
	"fmt"
	"html/template"
	"io"
)

const htmlTmpl = `<!DOCTYPE html>
<html>
<head><title>grpcannon Report</title>
<style>
body { font-family: sans-serif; margin: 2rem; }
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #ccc; padding: 0.5rem 1rem; text-align: left; }
th { background: #f4f4f4; }
.pass { color: green; } .fail { color: red; }
</style>
</head>
<body>
<h1>grpcannon Load Test Report</h1>
<h2>Summary</h2>
<table>
<tr><th>Metric</th><th>Value</th></tr>
<tr><td>Total Requests</td><td>{{.Total}}</td></tr>
<tr><td>Successful</td><td class="pass">{{.Success}}</td></tr>
<tr><td>Failed</td><td class="fail">{{.Failure}}</td></tr>
<tr><td>Success Rate</td><td>{{printf "%.2f" .SuccessRate}}%</td></tr>
<tr><td>Avg Latency</td><td>{{printf "%.3f" .AvgMs}} ms</td></tr>
<tr><td>P50 Latency</td><td>{{printf "%.3f" .P50}} ms</td></tr>
<tr><td>P95 Latency</td><td>{{printf "%.3f" .P95}} ms</td></tr>
<tr><td>P99 Latency</td><td>{{printf "%.3f" .P99}} ms</td></tr>
<tr><td>Min Latency</td><td>{{printf "%.3f" .Min}} ms</td></tr>
<tr><td>Max Latency</td><td>{{printf "%.3f" .Max}} ms</td></tr>
</table>
</body>
</html>
`

// WriteHTML writes an HTML report to w.
func WriteHTML(w io.Writer, r *Report) error {
	if r == nil {
		return fmt.Errorf("report is nil")
	}
	tmpl, err := template.New("report").Parse(htmlTmpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	return tmpl.Execute(w, r)
}
