package rules

import (
	"go/ast"
	"go/constant"
	"go/token"
	"unicode"

	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
)

type English struct{}

func NewEnglish() Rule {
	return &English{}
}

func (r *English) Name() string {
	return "english"
}

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

func (r *English) CheckCall(call *ast.CallExpr, pass *analysis.Pass) []analysis.Diagnostic {
	var diags []analysis.Diagnostic

	// Determine message index to skip it (it is checked by Check method)
	pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
	msgIndex := -1
	if ok {
		msgIndex = utils.GetMessageIndex(pkgPath, funcName)
	}

	utils.InspectLogArgs(pass, call, msgIndex, func(arg ast.Expr, isKey bool) {
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
