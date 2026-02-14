package rules

import (
	"go/token"
	"testing"

	"github.com/AlexanderGhosty/log-linter/pkg/logsupport"
)

func TestSymbols_Name(t *testing.T) {
	r := NewSymbols(logsupport.NewRegistry(nil), "")
	if r.Name() != "symbols" {
		t.Errorf("expected name 'symbols', got %q", r.Name())
	}
}

func TestSymbols_Check(t *testing.T) {
	r := NewSymbols(logsupport.NewRegistry(nil), "")

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		{name: "clean message", msg: "starting server", wantDiag: false},
		{name: "exclamation mark", msg: "error!", wantDiag: true},
		{name: "emoji rocket", msg: "deployed 🚀", wantDiag: true},
		{name: "emoji heart", msg: "love ❤️", wantDiag: true},
		{name: "multiple exclamations", msg: "failed!!!", wantDiag: true},
		{name: "at sign", msg: "user@host", wantDiag: true},
		{name: "hash", msg: "issue #123", wantDiag: true},
		{name: "curly braces", msg: "data: {key}", wantDiag: true},
		{name: "square brackets", msg: "list: [1,2]", wantDiag: true},
		{name: "empty string", msg: "", wantDiag: false},

		// Allowed characters
		{name: "period", msg: "server stopped.", wantDiag: false},
		{name: "comma", msg: "a, b, c", wantDiag: false},
		{name: "hyphen", msg: "re-deploy", wantDiag: false},
		{name: "underscore", msg: "user_id", wantDiag: false},
		{name: "colon", msg: "error: timeout", wantDiag: false},
		{name: "slash", msg: "/api/v1/users", wantDiag: false},
		{name: "equals", msg: "key=value", wantDiag: false},
		{name: "percent", msg: "100% done", wantDiag: false},
		{name: "parentheses", msg: "retry (3 times)", wantDiag: false},
		{name: "apostrophe", msg: "can't connect", wantDiag: false},
		{name: "digits", msg: "port 8080", wantDiag: false},
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

func TestSymbols_SuggestedFix(t *testing.T) {
	r := NewSymbols(logsupport.NewRegistry(nil), "")
	diags := r.Check("hello!🚀", token.NoPos, token.NoPos)

	if len(diags) == 0 {
		t.Fatal("expected diagnostic")
	}
	if len(diags[0].SuggestedFixes) == 0 {
		t.Fatal("expected suggested fix")
	}

	fix := diags[0].SuggestedFixes[0]
	newText := string(fix.TextEdits[0].NewText)
	// Should strip '!' and '🚀', leaving "hello"
	if newText != `"hello"` {
		t.Errorf("expected fix %q, got %q", `"hello"`, newText)
	}
}
