// Package config defines the configuration structure for the log-linter.
package config

import (
	"fmt"
	"regexp"
)

// Config holds the main configuration for the linter.
type Config struct {
	Symbols   SymbolsConfig   `mapstructure:"symbols"`
	Sensitive SensitiveConfig `mapstructure:"sensitive"`
	Loggers   []LoggerConfig  `mapstructure:"loggers"`
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	if err := c.Sensitive.Validate(); err != nil {
		return fmt.Errorf("sensitive config error: %w", err)
	}
	return nil
}

// SensitiveConfig holds configuration for sensitive data detection.
type SensitiveConfig struct {
	Keywords []string `mapstructure:"keywords"`
	Patterns []string `mapstructure:"patterns"`
}

// Validate checks the sensitive configuration for errors.
func (c *SensitiveConfig) Validate() error {
	for _, p := range c.Patterns {
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("invalid sensitive pattern %q: %w", p, err)
		}
	}
	return nil
}

// SymbolsConfig holds configuration for symbol restrictions.
type SymbolsConfig struct {
	Allowed string `mapstructure:"allowed"`
}

// LoggerConfig defines a custom logger configuration.
type LoggerConfig struct {
	// Package path of the logger (e.g. "log/slog", "go.uber.org/zap")
	Package string `mapstructure:"package"`
	// Implementation type: "slog", "zap", "generic"
	UserType string `mapstructure:"user_type"`
	// Names of field constructors (e.g. "String", "Int")
	FieldConstructors []string `mapstructure:"field_constructors"`
	// Index of the message argument in the log call
	MessageIndex int `mapstructure:"message_index"`
}
