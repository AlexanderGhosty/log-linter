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

	// Determine message index to skip it (it is checked by Check method if constant)
	pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
	msgIndex := -1
	if ok {
		msgIndex = utils.GetMessageIndex(pkgPath, funcName)
	}

	// Check the message argument itself if it's NOT a constant string (concatenation etc.)
	// If it IS a constant string, the basic Check method handles it.
	if msgIndex >= 0 && msgIndex < len(call.Args) {
		targetArg := call.Args[msgIndex]
		tv, ok := pass.TypesInfo.Types[targetArg]
		if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
			checkOperand(targetArg, r, report, "log message may contain sensitive data")
		}
	}

	utils.InspectLogArgs(pass, call, msgIndex, func(arg ast.Expr, isKey bool) {
		if isKey {
			// Check if key is a constant string
			tv, ok := pass.TypesInfo.Types[arg]
			if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
				val := constant.StringVal(tv.Value)
				if r.containsSensitiveInfo(val) {
					report(arg.Pos(), arg.End(), "log field key may contain sensitive data")
				}
			} else {
				// Key is not a constant string, analyze the expression
				checkOperand(arg, r, report, "log attribute contains sensitive data")
			}
		} else {
			// Value - check recursively
			checkOperand(arg, r, report, "log attribute contains sensitive data")
		}
	})

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
