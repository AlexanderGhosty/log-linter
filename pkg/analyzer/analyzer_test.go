package analyzer_test

import (
	"path/filepath"
	"testing"

	"github.com/AlexanderGhosty/log-linter/pkg/analyzer"
	"github.com/AlexanderGhosty/log-linter/pkg/config"
	"golang.org/x/tools/go/analysis/analysistest"

	_ "go.uber.org/zap" // Forced dependency for testdata
)

func TestAnalyzer(t *testing.T) {
	// Points to <project_root>/testdata
	testdata, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatal(err)
	}

	// Running tests on package "example" inside testdata/src/example
	analysistest.Run(t, testdata, analyzer.New(nil), "example")
}

func TestAnalyzer_Config(t *testing.T) {
	testdata, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sensitive: config.SensitiveConfig{
			Keywords: []string{"magic_word"},
		},
		Symbols: config.SymbolsConfig{
			Allowed: "@",
		},
	}

	analysistest.Run(t, testdata, analyzer.New(cfg), "configcheck")
}

func TestAnalyzer_CustomLogger(t *testing.T) {
	testdata, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Loggers: []config.LoggerConfig{{
			Package:      "custom",
			UserType:     "generic",
			MessageIndex: 0,
		}},
	}

	analysistest.Run(t, testdata, analyzer.New(cfg), "custom")
}
