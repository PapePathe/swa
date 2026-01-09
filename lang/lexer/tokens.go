package lexer

import (
	"fmt"
	"slices"
)

// Token ...
type Token struct {
	Value  string
	Name   string
	Kind   TokenKind `json:"-"`
	Line   int
	Column int
}

// NewToken ...
func NewToken(kind TokenKind, value string, line int) Token {
	return Token{Kind: kind, Value: value, Name: kind.String(), Line: line}
}

// Debug ...
func (t Token) Debug() {
	if t.isOneOfMany(Identifier, Number, String) {
		fmt.Printf("%s (%s)\n", t.Kind, t.Value)
	} else {
		fmt.Printf("%s\n", t.Kind)
	}
}

func (t Token) String() string {
	if t.isOneOfMany(Identifier, Number, String) {
		return fmt.Sprintf("%s (%s) (%d:%d)", t.Kind, t.Value, t.Line, t.Column)
	}

	return fmt.Sprintf("%s (%d:%d)", t.Kind, t.Line, t.Column)
}

func (t Token) isOneOfMany(expectedTokens ...TokenKind) bool {
	return slices.ContainsFunc(expectedTokens, func(v TokenKind) bool {
		return t.Kind == v
	})
}
