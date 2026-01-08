package lexer

import (
	"fmt"
	"os"
)

// Tokenize2 ...
func Tokenize2(source string) ([]Token, Dialect) {
	lex, dialect, err := NewFastLexer(source)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens, err := lex.GetAllTokens()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return tokens, dialect
}

// Tokenize ...
func Tokenize(source string) ([]Token, Dialect) {
	lex, dialect, err := New(source)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for !lex.atEOF() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())

			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)

				matched = true

				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("Lexer::Error -> unrecognized token near %s", lex.remainder()))
		}
	}

	lex.push(NewToken(EOF, "EOF", lex.line))

	return lex.Tokens, dialect
}
