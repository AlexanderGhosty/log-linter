package rules

import (
	"go/token"
	"testing"
)

func TestSensitive_Config(t *testing.T) {
	r := NewSensitive([]string{"beer", "wine"}, nil)

	tests := []struct {
		name     string
		input    string
		wantDiag bool
	}{
		{name: "beer", input: "beer", wantDiag: true},
		{name: "wine", input: "wine", wantDiag: true},
		{name: "water", input: "water", wantDiag: false},
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

func TestSymbols_Config(t *testing.T) {
	r := NewSymbols("@#")

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		// Custom allowed symbols
		{name: "allowed @", msg: "user@host", wantDiag: false},
		{name: "allowed #", msg: "issue #1", wantDiag: false},
		{name: "disallowed !", msg: "error!", wantDiag: true},
		// Default allowed symbols should still work
		{name: "default period", msg: "done.", wantDiag: false},
		{name: "default hyphen", msg: "re-try", wantDiag: false},
		{name: "default underscore", msg: "user_id", wantDiag: false},
		{name: "default colon", msg: "key: value", wantDiag: false},
		{name: "default slash", msg: "/api/v1", wantDiag: false},
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
