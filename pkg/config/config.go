package config

type Config struct {
	Symbols   SymbolsConfig   `mapstructure:"symbols"`
	Sensitive SensitiveConfig `mapstructure:"sensitive"`
}

type SensitiveConfig struct {
	Keywords []string `mapstructure:"keywords"`
}

type SymbolsConfig struct {
	Allowed string `mapstructure:"allowed"`
}
