package lexer

import (
	"fmt"
)

// Tokenize ...
func Tokenize(source string) ([]Token, Dialect, []error) {
	lex, dialect, err := New(source)
	if err != nil {
		return nil, nil, []error{err}
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
			lex.Errors = append(lex.Errors, fmt.Errorf("Lexer::Error -> unrecognized token near %s", lex.remainder()))
			lex.advanceN(1)
		}
	}

	lex.push(NewToken(EOF, "EOF", lex.line))

	return lex.Tokens, dialect, lex.Errors
}
