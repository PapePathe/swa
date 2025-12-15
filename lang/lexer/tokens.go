package lexer

import (
	"fmt"
	"slices"
)

// Token ...
type Token struct {
	Value string
	Name  string
	Kind  TokenKind `json:"-"`
	Line  int
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

func (t Token) isOneOfMany(expectedTokens ...TokenKind) bool {
	return slices.ContainsFunc(expectedTokens, func(v TokenKind) bool {
		return t.Kind == v
	})
}
