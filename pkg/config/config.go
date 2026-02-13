package config

import (
	"fmt"
	"regexp"
)

type Config struct {
	Symbols   SymbolsConfig   `mapstructure:"symbols"`
	Sensitive SensitiveConfig `mapstructure:"sensitive"`
}

func (c *Config) Validate() error {
	if err := c.Sensitive.Validate(); err != nil {
		return fmt.Errorf("sensitive config error: %w", err)
	}
	return nil
}

type SensitiveConfig struct {
	Keywords []string `mapstructure:"keywords"`
<<<<<<< HEAD
=======
	Patterns []string `mapstructure:"patterns"`
	Ignore   []string `mapstructure:"ignore"`
>>>>>>> 2f570fa (feat: enhance sensitive configuration with custom patterns and validation)
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
