package config

import "testing"

func TestSensitiveConfig_Validate(t *testing.T) {
	tests := []struct {
		cfg     *SensitiveConfig
		name    string
		wantErr bool
	}{
		{
			name: "valid patterns",
			cfg: &SensitiveConfig{
				Patterns: []string{`\d+`, `^foo`},
			},
			wantErr: false,
		},
		{
			name: "invalid pattern",
			cfg: &SensitiveConfig{
				Patterns: []string{`(`}, // unclosed group
			},
			wantErr: true,
		},
		{
			name: "empty patterns",
			cfg: &SensitiveConfig{
				Patterns: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
