package lexer

import (
	"fmt"
	"os"
)

// Tokenize ...
func Tokenize(source string) ([]Token, Dialect) {
	lex, dialect, err := New(source)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lex.tokenizeLoop()

	return lex.Tokens, dialect
}

func TokenizeWithDialect(source string, d Dialect) []Token {
	lex, err := NewWithDialect(source, d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lex.tokenizeLoop()

	return lex.Tokens
}
