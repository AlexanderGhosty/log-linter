package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/AlexanderGhosty/log-linter/pkg/analyzer"
)

func init() {
	register.Plugin("loglinter", New)
}

func New(conf any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

type Plugin struct{}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
