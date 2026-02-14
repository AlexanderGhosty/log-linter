package rules

import (
	"go/token"
	"testing"
)

func TestEnglish_Name(t *testing.T) {
	r := NewEnglish(nil)
	if r.Name() != "english" {
		t.Errorf("expected name 'english', got %q", r.Name())
	}
}

func TestEnglish_Check(t *testing.T) {
	r := NewEnglish(nil)

	tests := []struct {
		name     string
		msg      string
		wantDiag bool
	}{
		{name: "english only", msg: "starting server", wantDiag: false},
		{name: "cyrillic", msg: "запуск сервера", wantDiag: true},
		{name: "mixed english and cyrillic", msg: "server запуск", wantDiag: true},
		{name: "chinese characters", msg: "服务器启动", wantDiag: true},
		{name: "digits and symbols", msg: "error 404: /path", wantDiag: false},
		{name: "empty string", msg: "", wantDiag: false},
		{name: "only digits", msg: "12345", wantDiag: false},
		{name: "accented latin", msg: "café latte", wantDiag: true},
		{name: "german umlaut", msg: "über alles", wantDiag: true},
		{name: "japanese", msg: "テスト", wantDiag: true},
		{name: "emoji only", msg: "🚀", wantDiag: false}, // emoji is not a letter
		{name: "ascii punctuation", msg: "hello, world!", wantDiag: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := r.Check(tt.msg, token.NoPos, token.NoPos)
			gotDiag := len(diags) > 0
			if gotDiag != tt.wantDiag {
				t.Errorf("Check(%q): got diagnostic=%v, want diagnostic=%v", tt.msg, gotDiag, tt.wantDiag)
			}
			if gotDiag && diags[0].Message != "log message should be in English" {
				t.Errorf("unexpected message: %q", diags[0].Message)
			}
		})
	}
}
