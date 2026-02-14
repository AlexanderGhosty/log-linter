package rules

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
	"unicode"

	"github.com/AlexanderGhosty/log-linter/pkg/logsupport"
	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
)

// Symbols checks for disallowed symbols in log messages.
type Symbols struct {
	registry *logsupport.Registry
	allowed  string
}

// NewSymbols creates a new Symbols rule.
func NewSymbols(registry *logsupport.Registry, allowed string) Rule {
	if registry == nil {
		registry = logsupport.NewRegistry(nil)
	}

	return &Symbols{
		allowed:  allowed,
		registry: registry,
	}
}

// Name returns the name of the rule.
func (r *Symbols) Name() string {
	return "symbols"
}

// Check validates a single log message string.
func (r *Symbols) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	var cleanMsg []rune
	hasBadChars := false

	for _, ch := range msg {
		if r.isAllowed(ch) {
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

// CheckCall analyzes a full log call expression.
func (r *Symbols) CheckCall(call *ast.CallExpr, pass *analysis.Pass) []analysis.Diagnostic {
	var diags []analysis.Diagnostic

	// Determine message index
	pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
	msgIndex := -1
	if ok {
		msgIndex = r.registry.MessageIndex(pkgPath, funcName)
	}

	r.registry.InspectLogArgs(pass, call, msgIndex, func(arg ast.Expr, isKey bool) {
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

func (r *Symbols) isAllowed(ch rune) bool {
	if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == ' ' {
		return true
	}

	// Check configured allowed chars
	for _, allowedCh := range r.allowed {
		if ch == allowedCh {
			return true
		}
	}

	switch ch {
	case '.', ',', '-', '_', ':', '/', '=', '%', '(', ')', '\'':
		return true
	}
	return false
}
