package rules

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
)

type Symbols struct{}

func NewSymbols() Rule {
	return &Symbols{}
}

func (r *Symbols) Name() string {
	return "symbols"
}

func (r *Symbols) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	var cleanMsg []rune
	hasBadChars := false

	for _, ch := range msg {
		if isAllowed(ch) {
			cleanMsg = append(cleanMsg, ch)
		} else {
			hasBadChars = true
		}
	}

	if !hasBadChars {
		return nil
	}

	// Suggest replacement
	newMsg := string(cleanMsg)

	return []analysis.Diagnostic{{
		Pos:     pos,
		End:     end,
		Message: "log message should not contain special characters or emoji",
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: fmt.Sprintf("remove special characters -> %q", newMsg),
			TextEdits: []analysis.TextEdit{{
				Pos:     pos,
				End:     end,
				NewText: []byte(strconv.Quote(newMsg)),
			}},
		}},
	}}
}

func (r *Symbols) CheckCall(call *ast.CallExpr, pass *analysis.Pass) []analysis.Diagnostic {
	var diags []analysis.Diagnostic

	// Determine message index
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

		// 2. Check for key-value pairs in SugaredLogger (Infow, Warnw...) AND slog functions
		if utils.IsSupportedLogger(pkgPath, funcName) {
			isSlog := strings.Contains(pkgPath, "log/slog")
			isZapSugared := strings.HasSuffix(funcName, "w")

			if isSlog || isZapSugared {
				// Args after msgIndex are (key, value) pairs.
				// (i - msgIndex) starts at 1.
				// 1 -> key, 2 -> value, 3 -> key, 4 -> value...
				if (i-msgIndex)%2 == 1 {
					tv, ok := pass.TypesInfo.Types[arg]
					if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
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

func isAllowed(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
		return true
	}
	switch r {
	case '.', ',', '-', '_', ':', '/', '=', '%', '(', ')', '\'':
		return true
	}
	return false
}
