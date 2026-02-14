package config

import (
	"fmt"
	"regexp"
)

type Config struct {
	Symbols   SymbolsConfig   `mapstructure:"symbols"`
	Sensitive SensitiveConfig `mapstructure:"sensitive"`
	Loggers   []LoggerConfig  `mapstructure:"loggers"`
}

func (c *Config) Validate() error {
	if err := c.Sensitive.Validate(); err != nil {
		return fmt.Errorf("sensitive config error: %w", err)
	}
	return nil
}

type SensitiveConfig struct {
	Keywords []string `mapstructure:"keywords"`
	Patterns []string `mapstructure:"patterns"`
}

func (c *SensitiveConfig) Validate() error {
	for _, p := range c.Patterns {
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("invalid sensitive pattern %q: %w", p, err)
		}
	}
	return nil
}

type SymbolsConfig struct {
	Allowed string `mapstructure:"allowed"`
}

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
