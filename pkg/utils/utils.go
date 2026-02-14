package utils

import (
	"go/ast"
	"strings"

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

// NormalizeVendor strips the vendor prefix from a package path if present.
// This is useful when analyzing code that vendors its dependencies, as the
// package path might look like "github.com/org/repo/vendor/log/slog".
//
// Example:
//
//	NormalizeVendor("github.com/org/repo/vendor/log/slog") // returns "log/slog"
//	NormalizeVendor("log/slog")                            // returns "log/slog"
func NormalizeVendor(pkgPath string) string {
	if i := strings.Index(pkgPath, "/vendor/"); i >= 0 {
		return pkgPath[i+len("/vendor/"):]
	}
	return pkgPath
}
