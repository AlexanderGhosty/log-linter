package rules

import (
	"go/ast"
	"go/constant"
	"go/token"
	"unicode"

	"github.com/AlexanderGhosty/log-linter/pkg/logsupport"
	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
)

// English checks that log messages consist of ASCII characters.
type English struct {
	registry *logsupport.Registry
}

// NewEnglish creates a new English rule.
func NewEnglish(registry *logsupport.Registry) Rule {
	if registry == nil {
		registry = logsupport.NewRegistry(nil)
	}

	return &English{
		registry: registry,
	}
}

// Name returns the name of the rule.
func (r *English) Name() string {
	return "english"
}

// Check validates a single log message string.
func (r *English) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	const maxASCII = 127
	for _, ch := range msg {
		if unicode.IsLetter(ch) && ch > maxASCII {
			return []analysis.Diagnostic{{
				Pos:     pos,
				End:     end,
				Message: "log message should be in English",
			}}
		}
	}
	return nil
}

// CheckCall analyzes a full log call expression.
func (r *English) CheckCall(call *ast.CallExpr, pass *analysis.Pass) []analysis.Diagnostic {
	var diags []analysis.Diagnostic

	// Determine message index to skip it (it is checked by Check method)
	pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
	msgIndex := -1
	if ok {
		msgIndex = r.registry.MessageIndex(pkgPath, funcName)
	}

	r.registry.InspectLogArgs(pass, call, msgIndex, func(arg ast.Expr, isKey bool) {
		if isKey {
			// Check if arg is a constant string
			tv, ok := pass.TypesInfo.Types[arg]
			if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
				val := constant.StringVal(tv.Value)
				if d := r.Check(val, arg.Pos(), arg.End()); len(d) > 0 {
					diags = append(diags, d...)
				}
			}
		}
	})

	return diags
}
