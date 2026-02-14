package logsupport

import (
	"go/ast"
	"strings"

	"github.com/AlexanderGhosty/log-linter/pkg/config"
	"github.com/AlexanderGhosty/log-linter/pkg/utils"
	"golang.org/x/tools/go/analysis"
)

// Registry holds the configuration for supported loggers.
type Registry struct {
	configs []config.LoggerConfig
}

// NewRegistry creates a new Registry with the given matching configurations.
func NewRegistry(customConfigs []config.LoggerConfig) *Registry {
	// If custom configurations are provided (non-nil), use them exclusively.
	// This allows users to disable default loggers by providing an empty list,
	// or to define exactly the set of loggers they want.
	if customConfigs != nil {
		return &Registry{
			configs: customConfigs,
		}
	}

	// Otherwise, use defaults
	defaults := []config.LoggerConfig{
		{
			Package:      "log/slog",
			UserType:     "slog",
			MessageIndex: 0,
			FieldConstructors: []string{
				"String", "Int", "Int64", "Float64", "Bool", "Time", "Duration", "Any", "Group", "Attr",
			},
		},
		{
			Package:      "go.uber.org/zap",
			UserType:     "zap",
			MessageIndex: 0,
			FieldConstructors: []string{
				"String", "Int", "Int64", "Float64", "Bool", "Time", "Duration", "Any",
				"Binary", "ByteString", "Error", "NamedError", "Stringer",
				"Strings", "Ints", "Float64s", "Bools", "Times", "Durations",
				"Object", "Array", "Reflect", "Namespace", "Stack",
				"Int8", "Int16", "Int32", "Uint", "Uint8", "Uint16", "Uint32", "Uint64",
				"Float32", "Complex64", "Complex128", "Uintptr",
			},
		},
	}

	return &Registry{
		configs: defaults,
	}
}

// IsSupportedLogger returns true if the package and function correspond to a supported logger.
func (r *Registry) IsSupportedLogger(pkgPath, funcName string) bool {
	// Remove vendor prefix for matching
	pkgPath = normalizeVendor(pkgPath)

	for _, cfg := range r.configs {
		if cfg.Package == pkgPath {
			switch cfg.UserType {
			case "slog":
				switch funcName {
				case "Info", "Warn", "Error", "Debug", "Log":
					return true
				case "InfoContext", "WarnContext", "ErrorContext", "DebugContext", "LogAttrs":
					return true
				}
			case "zap":
				switch funcName {
				// Logger methods
				case "Info", "Warn", "Error", "Debug", "Panic", "Fatal", "DPanic":
					return true
				// SugaredLogger printf-style methods
				case "Infof", "Warnf", "Errorf", "Debugf", "Panicf", "Fatalf", "DPanicf":
					return true
				// SugaredLogger key-value methods
				case "Infow", "Warnw", "Errorw", "Debugw", "Panicw", "Fatalw", "DPanicw":
					return true
				}
			default:
				// For generic loggers, we might need more specific configuration
				// For now, assume if the package matches, it's supported
				return true
			}
		}
	}
	return false
}

// MessageIndex returns the index of the log message argument.
func (r *Registry) MessageIndex(pkgPath, funcName string) int {
	cleanPkgPath := normalizeVendor(pkgPath)

	for _, cfg := range r.configs {
		if cfg.Package == cleanPkgPath {
			if cfg.UserType == "slog" {
				switch funcName {
				case "Log", "LogAttrs":
					return 2
				default:
					if strings.HasSuffix(funcName, "Context") {
						return 1
					}
				}
				return 0
			}
			return cfg.MessageIndex
		}
	}
	return 0
}

// IsFieldConstructor returns true if the function is a field constructor.
func (r *Registry) IsFieldConstructor(pkgPath, funcName string) bool {
	cleanPkgPath := normalizeVendor(pkgPath)

	for _, cfg := range r.configs {
		if cfg.Package == cleanPkgPath {
			for _, fc := range cfg.FieldConstructors {
				if fc == funcName {
					return true
				}
			}
		}
	}
	return false
}

// InspectLogArgs iterates over the arguments of a log call.
func (r *Registry) InspectLogArgs(pass *analysis.Pass, call *ast.CallExpr, msgIndex int, fn func(arg ast.Expr, isKey bool)) {
	pkgPath, funcName, ok := utils.ResolveCallPackagePath(pass, call)
	if !ok {
		return
	}

	// Determine UserType
	cleanPkgPath := normalizeVendor(pkgPath)

	var userType string
	for _, cfg := range r.configs {
		if cfg.Package == cleanPkgPath {
			userType = cfg.UserType
			break
		}
	}

	for i, arg := range call.Args {
		if i <= msgIndex {
			continue
		}

		// Check for nested field constructors
		if callExpr, ok := arg.(*ast.CallExpr); ok {
			constructorPkg, constructorFunc, resolved := utils.ResolveCallPackagePath(pass, callExpr)
			if resolved {

				if r.IsFieldConstructor(constructorPkg, constructorFunc) {

					if len(callExpr.Args) > 0 {
						// The first arg is the key
						fn(callExpr.Args[0], true)
						// Subsequent args are values
						for _, valArg := range callExpr.Args[1:] {
							fn(valArg, false)
						}
					}
					continue
				}
			}
		}

		// Handle key-value pairs
		isSlog := userType == "slog"
		// For zap, only "w" suffixed methods are key-value pairs (sugared)
		isZapSugared := userType == "zap" && strings.HasSuffix(funcName, "w")

		if isSlog || isZapSugared {
			relativeIndex := i - msgIndex
			// In key-value pairs, odd indices (1, 3, ...) relative to message are keys
			isKey := relativeIndex%2 == 1

			fn(arg, isKey)
		}
	}
}

// normalizeVendor strips the vendor prefix from a package path if present.
func normalizeVendor(pkgPath string) string {
	if i := strings.Index(pkgPath, "/vendor/"); i >= 0 {
		return pkgPath[i+len("/vendor/"):]
	}
	return pkgPath
}
