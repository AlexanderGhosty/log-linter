package analyzer_test

import (
	"path/filepath"
	"testing"

	"github.com/AlexanderGhosty/log-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"

	_ "go.uber.org/zap" // Forced dependency for testdata
)

func TestAnalyzer(t *testing.T) {
	// Points to <project_root>/testdata
	testdata, _ := filepath.Abs("../../testdata")

	// Running tests on package "example" inside testdata/src/example
	analysistest.Run(t, testdata, analyzer.New(nil), "example")
}
