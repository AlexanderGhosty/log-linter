package configcheck

import "log/slog"

func Check() {
	slog.Info("magic_word") // want "log message may contain sensitive data"
	slog.Info("user@host")  // OK (allowed by config)
}
