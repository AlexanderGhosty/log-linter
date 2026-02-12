package rules

import (
	"fmt"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type Symbols struct{}

func NewSymbols() Rule {
	return &Symbols{}
}

func (r *Symbols) Name() string {
	return "symbols"
}

func (r *Symbols) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	var cleanMsg []rune
	hasBadChars := false

	for _, r := range msg {
		if isAllowed(r) {
			cleanMsg = append(cleanMsg, r)
		} else {
			hasBadChars = true
		}
	}

	if !hasBadChars {
		return nil
	}

	// Suggest replacement
	newMsg := string(cleanMsg)

	// If the message becomes empty, we probably shouldn't suggest replacing with empty string
	// unless that's intended logging. But removing all chars means the message was all invalid.

	return []analysis.Diagnostic{{
		Pos:     pos,
		End:     end,
		Message: "log message should not contain special characters or emoji",
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: fmt.Sprintf("remove special characters -> %q", newMsg),
			TextEdits: []analysis.TextEdit{{
				Pos:     pos,
				End:     end,
				NewText: []byte(strconv.Quote(newMsg)),
			}},
		}},
	}}
}

func isAllowed(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
		return true
	}
	switch r {
	case '.', ',', '-', '_', ':', '/', '=', '%', '(', ')', '\'':
		return true
	}
	return false
}
