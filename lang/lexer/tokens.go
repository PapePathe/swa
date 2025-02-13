package lexer

import "fmt"

type Token struct {
	Value string
	Kind  TokenKind
}

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

func (t Token) Debug() {
	if t.isOneOfMany(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", t.Kind, t.Value)
	} else {
		fmt.Printf("%s ()\n", t.Kind)
	}
}
