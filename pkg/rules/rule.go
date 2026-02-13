package rules

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule defines a single log message validation rule.
type Rule interface {
	// Name returns a unique identifier for the rule (e.g. "lowercase").
	Name() string

	// Check validates the log message and returns diagnostics if violated.
	Check(msg string, pos, end token.Pos) []analysis.Diagnostic
}

// ExprRule is an optional interface for rules that need
// access to the full call expression for deeper AST analysis
// (e.g., checking structured logger attributes).
type ExprRule interface {
	Rule
	CheckCall(call *ast.CallExpr, pass *analysis.Pass) []analysis.Diagnostic
}
