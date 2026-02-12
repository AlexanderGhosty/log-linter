package example

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func LowercaseRule() {
	// Task examples (Rule 1)
	slog.Info("Starting server on port 8080")   // want "log message should start with a lowercase letter"
	slog.Error("Failed to connect to database") // want "log message should start with a lowercase letter"
	slog.Info("starting server on port 8080")   // OK
	slog.Error("failed to connect to database") // OK

	// Zap
	logger := zap.NewExample()
	logger.Info("Starting server") // want "log message should start with a lowercase letter"
	logger.Info("starting server") // OK
}

func EnglishRule() {
	// Task examples (Rule 2)
	slog.Info("запуск сервера")                    // want "log message should be in English"
	slog.Error("ошибка подключения к базе данных") // want "log message should be in English"
	slog.Info("starting server")                   // OK
	slog.Error("failed to connect to database")    // OK
}

func SymbolsRule() {
	// Task examples (Rule 3)
	slog.Info("server started!🚀")      // want "log message should not contain special characters or emoji"
	slog.Error("connection failed!!!") // want "log message should not contain special characters or emoji"

	// Allowed chars
	slog.Info("user_id=123 (test)") // OK
	slog.Info("can't connect")      // OK
	slog.Info("server started")     // OK
	slog.Error("connection failed") // OK
}

func SensitiveRule() {
	password := "secret123"
	apiKey := "abc123"
	token := "xyz"

	// Task examples (Rule 4) — concatenation with sensitive data
	slog.Info("user password: " + password) // want "log message may contain sensitive data" "variable name suggests sensitive data"
	slog.Debug("api_key=" + apiKey)         // want "log message may contain sensitive data" "variable name suggests sensitive data"
	slog.Info("token: " + token)            // want "log message may contain sensitive data" "variable name suggests sensitive data"

	// Slog structured attributes
	slog.Info("login", "password", password)           // want "log attribute contains sensitive data"
	slog.Info("login", slog.String("token", "secret")) // want "log field key may contain sensitive data"

	// Zap fields
	logger := zap.NewExample()
	logger.Info("login", zap.String("password", "secret")) // want "log field key may contain sensitive data"

	// Correct examples (no diagnostic)
	slog.Info("user authenticated successfully") // OK
	slog.Debug("api request completed")          // OK
}

func ContextFuncs() {
	ctx := context.Background()
	slog.InfoContext(ctx, "Starting context func") // want "log message should start with a lowercase letter"
}
