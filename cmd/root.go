// Package cmd provides the CLI entry point and subcommands for grpcannon.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags shared across subcommands.
	target      string
	call        string
	concurrency int
	total       int
	timeout     int
	outputFmt   string
	baselineFile string
)

// rootCmd is the base command for grpcannon.
var rootCmd = &cobra.Command{
	Use:   "grpcannon",
	Short: "A lightweight gRPC load-testing CLI",
	Long: `grpcannon is a configurable gRPC load-testing tool that supports
concurrency control, latency histograms, and multiple output formats.

Examples:
  grpcannon run --target localhost:50051 --call helloworld.Greeter/SayHello
  grpcannon run --target localhost:50051 --call myservice.Svc/Ping --concurrency 20 --total 1000
  grpcannon baseline save --target localhost:50051 --call myservice.Svc/Ping
  grpcannon score --target localhost:50051 --call myservice.Svc/Ping`,
	SilenceUsage: true,
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all subcommands.
	rootCmd.PersistentFlags().StringVar(&target, "target", "", "gRPC server address (host:port)")
	rootCmd.PersistentFlags().StringVar(&call, "call", "", "Fully-qualified gRPC method (package.Service/Method)")
	rootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 10, "Number of concurrent workers")
	rootCmd.PersistentFlags().IntVar(&total, "total", 200, "Total number of requests to send")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Per-request timeout in seconds")
	rootCmd.PersistentFlags().StringVar(&outputFmt, "output", "text", "Output format: text, json, csv, table, markdown, prometheus, html, xml, influx")
	rootCmd.PersistentFlags().StringVar(&baselineFile, "baseline", "grpcannon_baseline.json", "Path to baseline file for comparison commands")

	// Register subcommands.
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

// versionCmd prints the current build version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the grpcannon version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grpcannon v0.1.0")
	},
}
