package rules

import (
	"fmt"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Lowercase checks that log messages start with a lowercase letter.
type Lowercase struct{}

// NewLowercase creates a new Lowercase rule.
func NewLowercase() Rule {
	return &Lowercase{}
}

// Name returns the name of the rule.
func (r *Lowercase) Name() string {
	return "lowercase"
}

// Check validates a single log message string.
func (r *Lowercase) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	if msg == "" {
		return nil
	}

	runes := []rune(msg)
	firstRune := runes[0]
	if !unicode.IsUpper(firstRune) {
		return nil
	}

	// Suggest replacement: lowercased first letter
	newMsg := string(unicode.ToLower(firstRune)) + string(runes[1:])

	return []analysis.Diagnostic{{
		Pos:     pos,
		End:     end,
		Message: "log message should start with a lowercase letter",
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: fmt.Sprintf("change to %q", newMsg),
			TextEdits: []analysis.TextEdit{{
				Pos:     pos,
				End:     end,
				NewText: []byte(strconv.Quote(newMsg)),
			}},
		}},
	}}
}
