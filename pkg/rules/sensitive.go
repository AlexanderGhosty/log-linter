package rules

import (
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
	"strings"

	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
)

type Sensitive struct {
	keywords []string
}

func NewSensitive() Rule {
	return &Sensitive{
		keywords: []string{
			"password", "passwd", "secret", "token",
			"api_key", "apikey", "access_key", "auth_token",
			"credential", "private_key",
		},
	}
}

func (r *Sensitive) Name() string {
	return "sensitive"
}

func (r *Sensitive) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	if r.containsSensitiveInfo(msg) {
		return []analysis.Diagnostic{{
			Pos:     pos,
			End:     end,
			Message: "log message may contain sensitive data",
		}}
	}
	return nil
}

func (r *Sensitive) CheckCall(call *ast.CallExpr, pass *analysis.Pass) []analysis.Diagnostic {
	var diags []analysis.Diagnostic

	// Helper to report diagnostic
	report := func(pos, end token.Pos, msg string) {
		diags = append(diags, analysis.Diagnostic{
			Pos:     pos,
			End:     end,
			Message: msg,
		})
	}

	// Determine message index to skip it (it is checked by Check method)
	pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
	msgIndex := -1
	if ok {
		msgIndex = utils.GetMessageIndex(pkgPath, funcName)
	}

	for i, arg := range call.Args {
		if i == msgIndex {
			// If the message argument is a constant string, it's checked by Rule.Check.
			// Skip it here to avoid duplicates.
			// If it's not constant (e.g. concatenation), Rule.Check won't run,
			// so we MUST check it here.
			tv, ok := pass.TypesInfo.Types[arg]
			if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
				continue
			}
		}

		// 1. Check for attribute constructors (slog.String("key", ...), zap.String("key", ...))
		if callExpr, ok := arg.(*ast.CallExpr); ok {
			// Resolve the function being called and verify it's a known field constructor
			constructorPkg, constructorFunc, resolved := utils.ResolveCallPackagePath(pass, callExpr)
			if resolved && utils.IsFieldConstructor(constructorPkg, constructorFunc) {
				// The first argument is the key.
				if len(callExpr.Args) > 0 {
					firstArg := callExpr.Args[0]

					// Check if first arg is constant string
					tv, ok := pass.TypesInfo.Types[firstArg]
					if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
						val := constant.StringVal(tv.Value)
						if r.containsSensitiveInfo(val) {
							report(firstArg.Pos(), firstArg.End(), "log field key may contain sensitive data")
						}
					}

					// Check values (subsequent arguments)
					for _, arg := range callExpr.Args[1:] {
						checkOperand(arg, r, report, "log attribute contains sensitive data")
					}
				}
				continue
			}
		}

		// 2. Check all other argument types (BasicLit, Ident, SelectorExpr, BinaryExpr, etc.)
		msg := "log attribute contains sensitive data"
		if i == msgIndex {
			msg = "log message may contain sensitive data"
		}
		checkOperand(arg, r, report, msg)
	}

	return diags
}

func checkOperand(expr ast.Expr, r *Sensitive, report func(token.Pos, token.Pos, string), stringLiteralMsg string) {
	// Check variable names for sensitive keywords
	if ident, ok := expr.(*ast.Ident); ok {
		if r.containsSensitiveInfo(ident.Name) {
			report(ident.Pos(), ident.End(), "variable name suggests sensitive data")
		}
	}
	// Check struct fields/selectors (e.g. user.Password)
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if r.containsSensitiveInfo(sel.Sel.Name) {
			report(sel.Pos(), sel.End(), "field name suggests sensitive data")
		}
		// Optionally verify the X expression too (e.g. secretStruct.Field)
		checkOperand(sel.X, r, report, stringLiteralMsg)
	}
	// Check string literal parts of concatenation for sensitive keywords
	if basicLit, ok := expr.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
		val, err := strconv.Unquote(basicLit.Value)
		if err == nil && r.containsSensitiveInfo(val) {
			report(basicLit.Pos(), basicLit.End(), stringLiteralMsg)
		}
	}
	// Recurse for nested binary expressions (e.g. a + b + c)
	if binExpr, ok := expr.(*ast.BinaryExpr); ok && binExpr.Op == token.ADD {
		checkOperand(binExpr.X, r, report, stringLiteralMsg)
		checkOperand(binExpr.Y, r, report, stringLiteralMsg)
	}
}

func (r *Sensitive) containsSensitiveInfo(s string) bool {
	s = strings.ToLower(s)
	for _, kw := range r.keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}
