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

func (lex *Lexer) Patterns() []RegexpPattern {
	return lex.patterns
}

func (lex *Lexer) advanceN(n int) {
	lex.position += n
}

func (lex *Lexer) newLine() {
	lex.line += 1
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

func getDialect(source string) (Dialect, error) {
	for _, d := range dialects {
		matches := d.DetectionPattern().FindStringSubmatch(source)

		if len(matches) > 0 {
			return d, nil
		}
	}

	return nil, fmt.Errorf("you must define your dialect")
}
