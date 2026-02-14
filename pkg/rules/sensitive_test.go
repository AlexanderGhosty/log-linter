package rules

import (
	"go/token"
	"testing"

	"github.com/AlexanderGhosty/log-linter/pkg/logsupport"
)

func TestSensitive_Name(t *testing.T) {
	r := NewSensitive(logsupport.NewRegistry(nil), nil, nil)
	if r.Name() != "sensitive" {
		t.Errorf("expected name 'sensitive', got %q", r.Name())
	}
}

func TestSensitive_Check(t *testing.T) {
	r := NewSensitive(logsupport.NewRegistry(nil), nil, nil)

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		// Should flag — contains sensitive keywords
		{name: "contains password", msg: "user password is set", wantDiag: true},
		{name: "contains passwd", msg: "passwd reset complete", wantDiag: true},
		{name: "contains token", msg: "invalid token received", wantDiag: true},
		{name: "contains secret", msg: "loading secret from vault", wantDiag: true},
		{name: "contains api_key", msg: "api_key loaded", wantDiag: true},
		{name: "contains apikey", msg: "apikey is missing", wantDiag: true},
		{name: "contains access_key", msg: "access_key rotated", wantDiag: true},
		{name: "contains auth_token", msg: "auth_token expired", wantDiag: true},
		{name: "contains credential", msg: "credential not found", wantDiag: true},
		{name: "contains private_key", msg: "private_key loaded", wantDiag: true},
		{name: "case insensitive PASSWORD", msg: "PASSWORD set", wantDiag: true},
		{name: "case insensitive Token", msg: "Token expired", wantDiag: true},
		{name: "keyword in middle", msg: "user password reset", wantDiag: true},
		{name: "keyword substring", msg: "my_password_field", wantDiag: true},

		// Should NOT flag — no sensitive keywords
		{name: "clean message", msg: "starting server", wantDiag: false},
		{name: "empty string", msg: "", wantDiag: false},
		{name: "similar word pass", msg: "passenger boarded", wantDiag: false},
		{name: "user login", msg: "user logged in", wantDiag: false},
		{name: "request processed", msg: "request processed successfully", wantDiag: false},
		{name: "connection timeout", msg: "connection timed out", wantDiag: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := r.Check(tt.msg, token.NoPos, token.NoPos)
			gotDiag := len(diags) > 0
			if gotDiag != tt.wantDiag {
				t.Errorf("Check(%q): got diagnostic=%v, want diagnostic=%v", tt.msg, gotDiag, tt.wantDiag)
			}
			if gotDiag && diags[0].Message != "log message may contain sensitive data" {
				t.Errorf("unexpected message: %q", diags[0].Message)
			}
		})
	}
}

func TestSensitive_ContainsSensitiveInfo(t *testing.T) {
	s := &Sensitive{
		keywords: []string{"password", "token", "secret"},
	}

	tests := []struct {
		input string
		want  bool
	}{
		{"password", true},
		{"my_password", true},
		{"PASSWORD", true},
		{"Token", true},
		{"top_secret_value", true},
		{"hello", false},
		{"", false},
		{"pass", false}, // "pass" is not a keyword, only "password"
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := s.containsSensitiveInfo(tt.input)
			if got != tt.want {
				t.Errorf("containsSensitiveInfo(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
