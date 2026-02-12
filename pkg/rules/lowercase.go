package rules

import (
	"fmt"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type Lowercase struct{}

func NewLowercase() Rule {
	return &Lowercase{}
}

func (r *Lowercase) Name() string {
	return "lowercase"
}

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
