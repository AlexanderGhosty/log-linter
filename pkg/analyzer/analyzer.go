// Package analyzer defines the loglinter Analyzer.
package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"

	"github.com/AlexanderGhosty/log-linter/pkg/config"
	"github.com/AlexanderGhosty/log-linter/pkg/logsupport"
	"github.com/AlexanderGhosty/log-linter/pkg/rules"
	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// New returns a new loglinter analyzer.
func New(cfg *config.Config) *analysis.Analyzer {
	if cfg == nil {
		cfg = &config.Config{}
	}

	registry := logsupport.NewRegistry(cfg.Loggers)

	registeredRules := []rules.Rule{
		rules.NewLowercase(),
		rules.NewEnglish(registry),
		rules.NewSymbols(registry, cfg.Symbols.Allowed),
		rules.NewSensitive(registry, cfg.Sensitive.Keywords, cfg.Sensitive.Patterns),
	}

	return &analysis.Analyzer{
		Name: "loglinter",
		Doc:  "checks log messages for common style issues",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass, registry, registeredRules)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass, registry *logsupport.Registry, registeredRules []rules.Rule) (interface{}, error) {
	inspectAnalyzer := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspectAnalyzer.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
		if !ok {
			return
		}

		if !registry.IsSupportedLogger(pkgPath, funcName) {
			return
		}

		msgArgIndex := registry.MessageIndex(pkgPath, funcName)
		msg, pos, end, found := extractLogMessage(pass, call, msgArgIndex)

		for _, rule := range registeredRules {
			// Basic string rules
			if found {
				if diags := rule.Check(msg, pos, end); len(diags) > 0 {
					for _, d := range diags {
						pass.Report(d)
					}
				}
			}

			// Advanced expression rules (check even if string literal wasn't found,
			// e.g. for analyzing context or other args)
			if exprRule, ok := rule.(rules.ExprRule); ok {
				if diags := exprRule.CheckCall(call, pass); len(diags) > 0 {
					for _, d := range diags {
						pass.Report(d)
					}
				}
			}
		}
	})

	return nil, nil
}

func extractLogMessage(pass *analysis.Pass, call *ast.CallExpr, argIndex int) (string, token.Pos, token.Pos, bool) {
	if len(call.Args) == 0 {
		return "", token.NoPos, token.NoPos, false
	}

	if len(call.Args) <= argIndex {
		return "", token.NoPos, token.NoPos, false
	}

	targetArg := call.Args[argIndex]

	tv, ok := pass.TypesInfo.Types[targetArg]
	if !ok || tv.Value == nil {
		return "", token.NoPos, token.NoPos, false
	}

	if tv.Value.Kind() != constant.String {
		return "", token.NoPos, token.NoPos, false
	}

	msg := constant.StringVal(tv.Value)
	return msg, targetArg.Pos(), targetArg.End(), true
}
