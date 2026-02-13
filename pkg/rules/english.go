package rules

import (
	"go/ast"
	"go/constant"
	"go/token"
	"strings"
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
	for _, ch := range msg {
		if unicode.IsLetter(ch) && ch > 127 {
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

	for i, arg := range call.Args {
		if i == msgIndex {
			continue
		}

		// 1. Check for attribute constructors (slog.String("key", ...), zap.String("key", ...))
		if callExpr, ok := arg.(*ast.CallExpr); ok {
			constructorPkg, constructorFunc, resolved := utils.ResolveCallPackagePath(pass, callExpr)
			if resolved && utils.IsFieldConstructor(constructorPkg, constructorFunc) {
				// The first argument is the key.
				if len(callExpr.Args) > 0 {
					firstArg := callExpr.Args[0]
					// Check if first arg is constant string
					tv, ok := pass.TypesInfo.Types[firstArg]
					if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
						val := constant.StringVal(tv.Value)
						if d := r.Check(val, firstArg.Pos(), firstArg.End()); len(d) > 0 {
							diags = append(diags, d...)
						}
					}
				}
				continue
			}
		}

		// 2. Check for key-value pairs in SugaredLogger (zap.Infow, etc.) AND slog functions
		// slog.Info("msg", "key", "value") -> "key" is at index 1 (msg at 0)
		// zap.Infow("msg", "key", "value") -> "key" is at index 1
		if utils.IsSupportedLogger(pkgPath, funcName) {
			isSlog := strings.Contains(pkgPath, "log/slog")
			isZapSugared := strings.HasSuffix(funcName, "w")

			if isSlog || isZapSugared {
				// Args after msgIndex are keys and values (mostly).
				// For slog, it can be mixed with Attrs, but if it's a string literal that is NOT an Attr constructor, it's likely a key.
				// A simple heuristic: if we are at an odd position relative to msgIndex (1st arg after msg, 3rd arg...), it's a key.
				// (i - msgIndex) starts at 1.
				if (i-msgIndex)%2 == 1 {
					tv, ok := pass.TypesInfo.Types[arg]
					// Check if it's a string constant
					if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
						// Additional check for slog: ensure it's not an Attr (like slog.String(...))
						// We already checked for CallExpr above. If it's a BasicLit string, it's a key.
						val := constant.StringVal(tv.Value)
						if d := r.Check(val, arg.Pos(), arg.End()); len(d) > 0 {
							diags = append(diags, d...)
						}
					}
				}
			}
		}
	}
	return diags
}
