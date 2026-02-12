package main

import (
	"github.com/AlexanderGhosty/log-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// main is a no-op required for build hygiene (go build ./...).
// This binary is not run directly — it is compiled into golangci-lint
// via `golangci-lint custom` using the New function below.
func main() {}

// New initializes the analyzer plugin.
// This function signature is required by golangci-lint module plugin system.
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}
