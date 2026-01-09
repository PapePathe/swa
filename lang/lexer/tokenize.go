package lexer

import (
	"fmt"
	"os"
)

// Tokenize ...
func Tokenize(source string) ([]Token, Dialect) {
	lex, dialect, err := NewFastLexer(source)
	if err != nil {
		fmt.Println(err)
		// TODO this method should return error instead of exiting the program
		os.Exit(1)
	}

	tokens, err := lex.GetAllTokens()
	if err != nil {
		fmt.Println(err)
		// TODO this method should return error instead of exiting the program
		os.Exit(1)
	}

	return tokens, dialect
}
