package utils

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// ResolveCallPackagePath uses type information to resolve the package path of the called function.
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

// IsSupportedLogger returns true if the package path and function name correspond to a supported logger.
func IsSupportedLogger(pkgPath, funcName string) bool {
	// Remove vendor prefix if present
	if i := strings.Index(pkgPath, "/vendor/"); i >= 0 {
		pkgPath = pkgPath[i+len("/vendor/"):]
	}

	switch pkgPath {
	case "log/slog":
		switch funcName {
		case "Info", "Warn", "Error", "Debug", "Log":
			return true
		case "InfoContext", "WarnContext", "ErrorContext", "DebugContext", "LogAttrs":
			return true
		}
	case "go.uber.org/zap":
		// Check for Zap methods
		if strings.HasSuffix(funcName, "f") || strings.HasSuffix(funcName, "w") {
			return true
		}
		switch funcName {
		case "Info", "Warn", "Error", "Debug", "Panic", "Fatal", "DPanic":
			return true
		}
	}
	return false
}

// GetMessageIndex returns the index of the log message argument.
// For slog:
//   - Info, Warn, Error, Debug: msg is at index 0
//   - InfoContext, WarnContext, ErrorContext, DebugContext: msg is at index 1 (ctx, msg, ...)
//   - Log, LogAttrs: msg is at index 2 (ctx, level, msg, ...)
//
// For zap: typically 0.
func GetMessageIndex(pkgPath, funcName string) int {
	if i := strings.Index(pkgPath, "/vendor/"); i >= 0 {
		pkgPath = pkgPath[i+len("/vendor/"):]
	}

	if pkgPath == "log/slog" {
		switch funcName {
		case "Log", "LogAttrs":
			return 2
		default:
			if strings.HasSuffix(funcName, "Context") {
				return 1
			}
		}
	}
	// Default to 0 (slog.Info/Warn/Error/Debug, zap.*)
	return 0
}
