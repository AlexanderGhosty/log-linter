package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func f() {
	// Setup zap
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	// zap: SugaredLogger (w methods)
	// 1. Normal key-value pairs
	sugar.Infow("msg", "key1", "val1") // want "func: Infow, args: Arg \"key1\" true, Arg \"val1\" false"

	// 2. Odd number of arguments
	sugar.Warnw("msg", "key1", "val1", "dangling") // want "func: Warnw, args: Arg \"key1\" true, Arg \"val1\" false, Arg \"dangling\" true"

	// 3. Non-string keys (should still be identified as keys by position)
	sugar.Errorw("msg", 123, "val1") // want "func: Errorw, args: Arg 123 true, Arg \"val1\" false"

	// zap: Logger (field constructors)
	// 4. Field constructors
	logger.Info("msg", zap.String("key1", "val1")) // want "func: Info, args: Arg \"key1\" true, Arg \"val1\" false"

	// 5. Nested calls
	sugar.Debugw("msg", "key1", "val1", "key2", "val2") // want "func: Debugw, args: Arg \"key1\" true, Arg \"val1\" false, Arg \"key2\" true, Arg \"val2\" false"

	// slog: Structured logging
	// 6. Key-value pairs
	slog.Info("msg", "key1", "val1") // want "func: Info, args: Arg \"key1\" true, Arg \"val1\" false"

	// 7. Attributes
	slog.Info("msg", slog.String("key1", "val1")) // want "func: Info, args: Arg \"key1\" true, Arg \"val1\" false"

	// 8. Log with mixed args
	slog.Info("msg", "key1", "val1", slog.Int("count", 42)) // want "func: Info, args: Arg \"key1\" true, Arg \"val1\" false, Arg \"count\" true, Arg 42 false"
}
