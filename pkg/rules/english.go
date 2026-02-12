package rules

import (
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type English struct{}

func NewEnglish() Rule {
	return &English{}
}

func (r *English) Name() string {
	return "english"
}

func (r *English) Check(msg string, pos, end token.Pos) []analysis.Diagnostic {
	for _, ch := range msg {
		if unicode.IsLetter(ch) && ch > 127 {
			return []analysis.Diagnostic{{
				Pos:     pos,
				End:     end,
				Message: "log message should be in English",
			}}
		}
	}
	return nil
}
