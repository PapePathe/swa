/*
* swahili/lang
* Copyright (C) 2025  Papa Pathe SENE
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package lexer

import (
	"fmt"
	"regexp"
)

// Lexer ...
type Lexer struct {
	Tokens        []Token              // The tokens
	source        string               // The source code
	position      int                  // The current position of the lexer
	patterns      []RegexpPattern      // the list of patterns of the language
	reservedWords map[string]TokenKind // list of reserved words
}

func New(source string) (*Lexer, error) {
	dialect, err := getDialect(source)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		Tokens:        make([]Token, 0),
		patterns:      dialect.Patterns(),
		reservedWords: dialect.Reserved(),
		source:        source,
	}, nil
}

func (lex *Lexer) advanceN(n int) {
	lex.position += n
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

func getDialect(source string) (Dialect, error) {
	re := regexp.MustCompile(`dialect:([a-zA-Z]+);`)
	matches := re.FindStringSubmatch(source)

	if len(matches) > 1 {
		dialect, ok := dialects[matches[1]]
		if !ok {
			return nil, fmt.Errorf("dialect <%s> is not supported", matches[1])
		}

		return dialect, nil
	} else {
		return nil, fmt.Errorf("you must define your dialect")
	}
}
