package rules

import (
	"go/token"
	"testing"
)

func TestLowercase_Name(t *testing.T) {
	r := NewLowercase()
	if r.Name() != "lowercase" {
		t.Errorf("expected name 'lowercase', got %q", r.Name())
	}
}

func TestLowercase_Check(t *testing.T) {
	r := NewLowercase()

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		{name: "uppercase start", msg: "Starting server", wantDiag: true},
		{name: "lowercase start", msg: "starting server", wantDiag: false},
		{name: "empty string", msg: "", wantDiag: false},
		{name: "digit start", msg: "123 error", wantDiag: false},
		{name: "unicode uppercase", msg: "Über cool", wantDiag: true},
		{name: "single uppercase", msg: "A", wantDiag: true},
		{name: "single lowercase", msg: "a", wantDiag: false},
		{name: "space start", msg: " hello", wantDiag: false},
		{name: "symbol start", msg: "/path/to", wantDiag: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := r.Check(tt.msg, token.NoPos, token.NoPos)
			gotDiag := len(diags) > 0
			if gotDiag != tt.wantDiag {
				t.Errorf("Check(%q): got diagnostic=%v, want diagnostic=%v", tt.msg, gotDiag, tt.wantDiag)
			}
			if gotDiag && diags[0].Message != "log message should start with a lowercase letter" {
				t.Errorf("unexpected message: %q", diags[0].Message)
			}
		})
	}
}

func TestLowercase_SuggestedFix(t *testing.T) {
	r := NewLowercase()
	diags := r.Check("Hello world", token.NoPos, token.NoPos)

	if len(diags) == 0 {
		t.Fatal("expected diagnostic")
	}
	if len(diags[0].SuggestedFixes) == 0 {
		t.Fatal("expected suggested fix")
	}

	fix := diags[0].SuggestedFixes[0]
	if len(fix.TextEdits) == 0 {
		t.Fatal("expected text edit")
	}

	// The new text should be the quoted lowercase version
	newText := string(fix.TextEdits[0].NewText)
	if newText != `"hello world"` {
		t.Errorf("expected fix %q, got %q", `"hello world"`, newText)
	}
}
