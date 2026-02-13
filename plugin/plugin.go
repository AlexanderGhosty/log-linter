package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/tools/go/analysis"

	"github.com/AlexanderGhosty/log-linter/pkg/analyzer"
	"github.com/AlexanderGhosty/log-linter/pkg/config"
)

func init() {
	register.Plugin("loglinter", New)
}

func New(conf any) (register.LinterPlugin, error) {
	var cfg config.Config
	if err := mapstructure.Decode(conf, &cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Plugin{cfg: cfg}, nil
}

type Plugin struct {
	cfg config.Config
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.New(&p.cfg)}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
