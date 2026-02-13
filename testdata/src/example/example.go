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
	slog.Info("user password: " + password) // want "sensitive data" "sensitive data"
	slog.Debug("api_key=" + apiKey)         // want "sensitive data" "sensitive data"
	slog.Info("token: " + token)            // want "sensitive data" "sensitive data"

	// Slog structured attributes
	slog.Info("login", "password", password)           // want "sensitive data" "sensitive data"
	slog.Info("login", slog.String("token", "secret")) // want "sensitive data" "sensitive data"

	// Zap fields
	logger := zap.NewExample()
	logger.Info("login", zap.String("password", "secret")) // want "sensitive data" "sensitive data"

	// Correct examples (no diagnostic)
	slog.Info("user authenticated successfully") // OK
	slog.Debug("api request completed")          // OK

	// Direct variable usage (should be flagged)
	slog.Info(password) // want "sensitive data"

	// Struct field access
	type User struct {
		Password string
		Token    string
	}
	u := User{}
	slog.Info("logging user", "pass", u.Password)             // want "sensitive data"
	slog.Info("logging value", slog.String("token", u.Token)) // want "sensitive data" "sensitive data"

	// Sensitive values in constructors
	slog.Info("config", slog.String("key", "my_secret_token"))  // want "sensitive data"
	logger.Info("config", zap.String("key", "my_secret_token")) // want "sensitive data"
}

func ContextFuncs() {
	ctx := context.Background()
	slog.InfoContext(ctx, "Starting context func") // want "log message should start with a lowercase letter"
}

func StructuredKeys() {
	// Structured logging (slog)
	slog.Info("server started", "key", "value")
	slog.Info("server started", "ключ", "value") // want "log message should be in English"
	slog.Info("server started", "key!", "value") // want "log message should not contain special characters or emoji"

	// Structured logging (zap)
	logger := zap.NewExample()
	sugar := zap.S()

	sugar.Infow("server started",
		"key", "value",
		"ключ", "value", // want "log message should be in English"
		"key!", "value", // want "log message should not contain special characters or emoji"
	)

	logger.Info("server started",
		zap.String("key", "value"),
		zap.String("ключ", "value"), // want "log message should be in English"
		zap.String("key!", "value"), // want "log message should not contain special characters or emoji"
	)
}
