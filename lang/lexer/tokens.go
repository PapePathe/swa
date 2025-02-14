package lexer

import "fmt"

// Token ...
type Token struct {
	Value string
	Kind  TokenKind
}

// NewToken ...
func NewToken(kind TokenKind, value string) Token {
	return Token{Kind: kind, Value: value}
}

func (t Token) isOneOfMany(expectedTokens ...TokenKind) bool {
	for _, tk := range expectedTokens {
		if tk == t.Kind {
			return true
		}
	}

	return false
}

// Debug ...
func (t Token) Debug() {
	if t.isOneOfMany(Identifier, Number, String) {
		fmt.Printf("%s (%s)\n", t.Kind, t.Value)
	} else {
		fmt.Printf("%s\n", t.Kind)
	}
}
