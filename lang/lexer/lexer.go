package lexer

import (
	"fmt"
)

// Lexer ...
type Lexer struct {
	Tokens        []Token // The tokens
	source        string  // The source code
	position      int     // The current position of the lexer
	line          int
	patterns      []RegexpPattern      // the list of patterns of the language
	reservedWords map[string]TokenKind // list of reserved words
}

func NewWithDialect(source string, d Dialect) (*Lexer, error) {
	return &Lexer{
		Tokens:        make([]Token, 0),
		patterns:      d.Patterns(),
		reservedWords: d.Reserved(),
		source:        source,
	}, nil
}

func New(source string) (*Lexer, Dialect, error) {
	dialect, err := getDialect(source)
	if err != nil {
		return nil, nil, err
	}

	return &Lexer{
		Tokens:        make([]Token, 0),
		patterns:      dialect.Patterns(),
		reservedWords: dialect.Reserved(),
		source:        source,
	}, dialect, nil
}

func (lex *Lexer) advanceN(n int) {
	lex.position += n
}

func (lex *Lexer) newLine() {
	lex.line += 1
}

func (lex *Lexer) Patterns() []RegexpPattern {
	return lex.patterns
}

func (lex *Lexer) remainder() string {
	return lex.source[lex.position:]
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) atEOF() bool {
	return lex.position >= len(lex.source)
}

func (lex *Lexer) tokenizeLoop() {
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
}

func getDialect(source string) (Dialect, error) {
	for _, d := range dialects {
		matches := d.DetectionPattern().FindStringSubmatch(source)

		if len(matches) > 0 {
			return d, nil
		}
	}

	return nil, fmt.Errorf("you must define your dialect")
}
