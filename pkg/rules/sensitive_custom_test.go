package rules

import (
	"go/token"
	"testing"

	"github.com/AlexanderGhosty/log-linter/pkg/logsupport"
)

func TestSensitive_CustomPatterns(t *testing.T) {
	// Define custom patterns: e.g. Credit Card-like pattern
	patterns := []string{
		`\d{4}-\d{4}-\d{4}-\d{4}`, // Simple CC pattern
	}

	r := NewSensitive(logsupport.NewRegistry(nil), nil, patterns)

	tests := []struct {
		name     string
		input    string
		wantDiag bool
	}{
		{
			name:     "match custom pattern",
			input:    "Payment processed: 1234-5678-9012-3456",
			wantDiag: true,
		},
		{
			name:     "no match",
			input:    "Payment processed: success",
			wantDiag: false,
		},
		{
			name:     "match default keyword",
			input:    "password: 123",
			wantDiag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := r.Check(tt.input, token.NoPos, token.NoPos)
			gotDiag := len(diags) > 0
			if gotDiag != tt.wantDiag {
				t.Errorf("Check(%q) diag = %v, want %v", tt.input, gotDiag, tt.wantDiag)
			}
		})
	}
}
