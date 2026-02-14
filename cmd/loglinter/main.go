// Package main is the entry point for the loglinter executable.
package main

import (
	"github.com/AlexanderGhosty/log-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.New(nil))
}
