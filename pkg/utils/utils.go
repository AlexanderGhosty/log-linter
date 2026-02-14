package utils

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ResolveCallPackagePath resolves the package path and function name of a call expression.
// It uses the type information available in the analysis.Pass to determine the called object.
//
// Returns:
//   - pkgPath: The full package path of the called function (e.g., "log/slog").
//   - funcName: The name of the called function (e.g., "Info").
//   - ok: True if resolution was successful, false otherwise.
func ResolveCallPackagePath(pass *analysis.Pass, call *ast.CallExpr) (string, string, bool) {
	fun := call.Fun

	// Handle method calls: logger.Info(...) or pkg.Func(...)
	if selExpr, ok := fun.(*ast.SelectorExpr); ok {
		// Identify the object being called (the function/method itself)
		selObj := pass.TypesInfo.ObjectOf(selExpr.Sel)
		if selObj == nil {
			return "", "", false
		}

		if pkg := selObj.Pkg(); pkg != nil {
			return pkg.Path(), selObj.Name(), true
		}
	}

	return "", "", false
}
