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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	t.Run("tokenize characters", func(t *testing.T) {
		tests := []struct {
			name     string
			src      string
			expected []Token
		}{
			{
				name: "with ascii character",
				src:  `dialect:french; 'a';`,
				expected: []Token{
					{Value: "dialect", Name: "DIALECT", Kind: 39},
					{Value: ":", Name: "COLON", Kind: 10},
					{Value: "french", Name: "IDENTIFIER", Kind: 35},
					{Value: ";", Name: "SEMI_COLON", Kind: 33},
					{Value: "'a'", Name: "CHARACTER", Kind: 41},
					{Value: ";", Name: "SEMI_COLON", Kind: 33},
					{Value: "EOF", Name: "EOF", Kind: 16},
				},
			},
			{
				name: "with utf8 character",
				src:  `dialect:french; 'è';`,
				expected: []Token{
					{Value: "dialect", Name: "DIALECT", Kind: 39},
					{Value: ":", Name: "COLON", Kind: 10},
					{Value: "french", Name: "IDENTIFIER", Kind: 35},
					{Value: ";", Name: "SEMI_COLON", Kind: 33},
					{Value: "'è'", Name: "CHARACTER", Kind: 41},
					{Value: ";", Name: "SEMI_COLON", Kind: 33},
					{Value: "EOF", Name: "EOF", Kind: 16},
				},
			},

			{
				name: "if else block",
				src: `
        	dialect:french;
          si(x>0) {
            width = 100;
            height = 100 + 400 - width;
          } sinon {
          	width += 300;
          }
        `,
				expected: []Token{
					{Name: "DIALECT", Kind: DialectDeclaration, Value: "dialect"},
					{Name: "COLON", Kind: Colon, Value: ":"},
					{Name: "IDENTIFIER", Kind: Identifier, Value: "french"},
					{Name: "SEMI_COLON", Kind: SemiColon, Value: ";"},
					{Name: "IF", Kind: KeywordIf, Value: "si"},
					{Name: "OPEN_PAREN", Kind: OpenParen, Value: "("},
					{Name: "IDENTIFIER", Kind: Identifier, Value: "x"},
					{Name: "GREATER_THAN", Kind: GreaterThan, Value: ">"},
					{Name: "NUMBER", Kind: Number, Value: "0"},
					{Name: "CLOSE_PAREN", Kind: CloseParen, Value: ")"},
					{Name: "OPEN_CURLY", Kind: OpenCurly, Value: "{"},
					{Name: "IDENTIFIER", Kind: Identifier, Value: "width"},
					{Name: "ASSIGNMENT", Kind: Assignment, Value: "="},
					{Name: "NUMBER", Kind: Number, Value: "100"},
					{Name: "SEMI_COLON", Kind: SemiColon, Value: ";"},
					{Name: "IDENTIFIER", Kind: Identifier, Value: "height"},
					{Name: "ASSIGNMENT", Kind: Assignment, Value: "="},
					{Name: "NUMBER", Kind: Number, Value: "100"},
					{Name: "PLUS", Kind: Plus, Value: "+"},
					{Name: "NUMBER", Kind: Number, Value: "400"},
					{Name: "MINUS", Kind: Minus, Value: "-"},
					{Name: "IDENTIFIER", Kind: Identifier, Value: "width"},
					{Name: "SEMI_COLON", Kind: SemiColon, Value: ";"},
					{Name: "CLOSE_CURLY", Kind: CloseCurly, Value: "}"},
					{Name: "ELSE", Kind: KeywordElse, Value: "sinon"},
					{Name: "OPEN_CURLY", Kind: OpenCurly, Value: "{"},
					{Name: "IDENTIFIER", Kind: Identifier, Value: "width"},
					{Name: "PLUS_EQUAL", Kind: PlusEquals, Value: "+="},
					{Name: "NUMBER", Kind: Number, Value: "300"},
					{Name: "SEMI_COLON", Kind: SemiColon, Value: ";"},
					{Name: "CLOSE_CURLY", Kind: CloseCurly, Value: "}"},
					{Name: "EOF", Kind: EOF, Value: "EOF"},
				},
			},
		}

		for _, test := range tests {
			result := Tokenize(test.src)
			assert.Equal(t, test.expected, result)
		}
	})
}
