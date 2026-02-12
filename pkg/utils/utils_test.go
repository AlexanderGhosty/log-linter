package utils

import "testing"

func TestIsSupportedLogger(t *testing.T) {
	tests := []struct {
		name     string
		pkgPath  string
		funcName string
		want     bool
	}{
		// slog
		{"slog Info", "log/slog", "Info", true},
		{"slog Warn", "log/slog", "Warn", true},
		{"slog Error", "log/slog", "Error", true},
		{"slog Debug", "log/slog", "Debug", true},
		{"slog Log", "log/slog", "Log", true},
		{"slog InfoContext", "log/slog", "InfoContext", true},
		{"slog WarnContext", "log/slog", "WarnContext", true},
		{"slog ErrorContext", "log/slog", "ErrorContext", true},
		{"slog DebugContext", "log/slog", "DebugContext", true},
		{"slog LogAttrs", "log/slog", "LogAttrs", true},
		{"slog unsupported", "log/slog", "With", false},
		// zap
		{"zap Info", "go.uber.org/zap", "Info", true},
		{"zap Warn", "go.uber.org/zap", "Warn", true},
		{"zap Error", "go.uber.org/zap", "Error", true},
		{"zap Debug", "go.uber.org/zap", "Debug", true},
		{"zap Fatal", "go.uber.org/zap", "Fatal", true},
		{"zap Panic", "go.uber.org/zap", "Panic", true},
		{"zap DPanic", "go.uber.org/zap", "DPanic", true},
		{"zap sugared Infof", "go.uber.org/zap", "Infof", true},
		{"zap sugared Infow", "go.uber.org/zap", "Infow", true},
		{"zap sugared Warnf", "go.uber.org/zap", "Warnf", true},
		{"zap unsupported", "go.uber.org/zap", "Named", false},
		// other packages
		{"fmt Println", "fmt", "Println", false},
		{"log Printf", "log", "Printf", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSupportedLogger(tt.pkgPath, tt.funcName)
			if got != tt.want {
				t.Errorf("IsSupportedLogger(%q, %q) = %v, want %v", tt.pkgPath, tt.funcName, got, tt.want)
			}
		})
	}
}

func TestGetMessageIndex(t *testing.T) {
	tests := []struct {
		name     string
		pkgPath  string
		funcName string
		want     int
	}{
		// slog: bare functions — msg at index 0
		{"slog Info", "log/slog", "Info", 0},
		{"slog Warn", "log/slog", "Warn", 0},
		{"slog Error", "log/slog", "Error", 0},
		{"slog Debug", "log/slog", "Debug", 0},
		// slog: context functions — msg at index 1 (ctx, msg, ...)
		{"slog InfoContext", "log/slog", "InfoContext", 1},
		{"slog WarnContext", "log/slog", "WarnContext", 1},
		{"slog ErrorContext", "log/slog", "ErrorContext", 1},
		{"slog DebugContext", "log/slog", "DebugContext", 1},
		// slog: Log/LogAttrs — msg at index 2 (ctx, level, msg, ...)
		{"slog Log", "log/slog", "Log", 2},
		{"slog LogAttrs", "log/slog", "LogAttrs", 2},
		// zap: msg at index 0
		{"zap Info", "go.uber.org/zap", "Info", 0},
		{"zap Infof", "go.uber.org/zap", "Infof", 0},
		// vendored
		{"vendored slog InfoContext", "myproject/vendor/log/slog", "InfoContext", 1},
		{"vendored slog LogAttrs", "myproject/vendor/log/slog", "LogAttrs", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMessageIndex(tt.pkgPath, tt.funcName)
			if got != tt.want {
				t.Errorf("GetMessageIndex(%q, %q) = %v, want %v", tt.pkgPath, tt.funcName, got, tt.want)
			}
		})
	}
}
