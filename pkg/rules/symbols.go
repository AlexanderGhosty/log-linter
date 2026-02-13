package rules

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
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

	utils.InspectLogArgs(pass, call, msgIndex, func(arg ast.Expr, isKey bool) {
		if isKey {
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
