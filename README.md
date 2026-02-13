# Log-Linter

A Go linter designed to enforce logging standards and style guidelines. Compatible with `golangci-lint`.

## Features

Enforces the following rules on log messages (for `log/slog` and `go.uber.org/zap`):

1. **Lowercase**: Log messages should start with a lowercase letter.
   - ❌ `log.Info("Starting server")`
   - ✅ `log.Info("starting server")` (suggests auto-fix)

2. **English Only**: Log messages should be in English (ASCII only).
   - ❌ `log.Info("запуск сервера")`
   - ✅ `log.Info("starting server")`

3. **No Special Characters**: Log messages should not contain special characters or emojis.
   - ❌ `log.Info("server started!🚀")`
   - ✅ `log.Info("server started")` (suggests auto-fix)
   - Whitelist: Letters, digits, space, `.`, `,`, `-`, `_`, `:`, `/`, `=`, `%`, `(`, `)`, `'`

4. **No Sensitive Data**: Log messages and attributes/fields should not contain sensitive information.
   - Checks message content for keywords (`password`, `token`, etc.).
   - Checks structured logging keys (`slog` string keys, `zap.String` field keys) for sensitive names.
   - Checks variable names in string concatenation (legacy style).
   - ❌ `slog.Info("user password: " + password)`
   - ❌ `slog.Info("login", "password", p)`

## Requirements

- Go 1.23+

## Installation & Usage

### 1. Standalone CLI

You can build and run the linter as a standalone tool:

```bash
make build
./loglinter ./path/to/your/package/...
```

To run on included example:

```bash
make lint-example
```

### 2. Integration with golangci-lint

This linter is designed to work as a [Module Plugin](https://golangci-lint.run/plugins/module-plugins/) for `golangci-lint`.

#### Step 1: Create a custom build configuration `.custom-gcl.yml`

```yaml
version: v2.8.0
plugins:
  - module: 'github.com/AlexanderGhosty/log-linter'
    import: 'github.com/AlexanderGhosty/log-linter/plugin'
    version: v0.1.0 # example tag; replace with a specific version or commit
    path: . # If using local source, or omit for remote module
```

#### Step 2: Build custom binary

```bash
golangci-lint custom
```

This will produce a `custom-gcl` binary.

#### Step 3: Configure `.golangci.yml`

In your project configuration:

```yaml
linters-settings:
  custom:
    loglinter:
      path: ./custom-gcl 
      description: Check log messages for style guide compliance
      original-url: github.com/AlexanderGhosty/log-linter

linters:
  enable:
    - loglinter
```

## Supported Loggers

- `log/slog`: `Info`, `Warn`, `Error`, `Debug`, `Log`, `LogAttrs`, and `*Context` variants.
- `go.uber.org/zap`: `Info`, `Warn`, `Error`, `Debug`, `Fatal`, `Panic`, `DPanic`, and `*f`, `*w` variants.
