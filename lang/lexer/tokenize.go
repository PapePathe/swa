package lexer

import (
	"fmt"
	"os"
)

// Tokenize ...
func Tokenize(source string) []Token {
	lex, err := New(source)
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

	lex.push(NewToken(EOF, "EOF"))

	return lex.Tokens
}
