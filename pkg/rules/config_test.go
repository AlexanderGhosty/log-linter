package rules

import (
	"go/token"
	"testing"
)

func TestSensitive_Config(t *testing.T) {
	r := NewSensitive([]string{"beer", "wine"})

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		{name: "custom keyword beer", msg: "free beer", wantDiag: true},
		{name: "custom keyword wine", msg: "red wine", wantDiag: true},
		// Default keywords are REPLACED, not merged in recent implementation.
		{name: "default keyword password", msg: "user password", wantDiag: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := r.Check(tt.msg, token.NoPos, token.NoPos)
			gotDiag := len(diags) > 0
			if gotDiag != tt.wantDiag {
				t.Errorf("Check(%q): got diagnostic=%v, want diagnostic=%v", tt.msg, gotDiag, tt.wantDiag)
			}
		})
	}
}

func TestSymbols_Config(t *testing.T) {
	r := NewSymbols("@#")

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		{name: "allowed @", msg: "user@host", wantDiag: false},
		{name: "allowed #", msg: "issue #1", wantDiag: false},
		{name: "disallowed !", msg: "error!", wantDiag: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := r.Check(tt.msg, token.NoPos, token.NoPos)
			gotDiag := len(diags) > 0
			if gotDiag != tt.wantDiag {
				t.Errorf("Check(%q): got diagnostic=%v, want diagnostic=%v", tt.msg, gotDiag, tt.wantDiag)
			}
		})
	}
}
