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

import "fmt"

// Token ...
type Token struct {
	Value string
	Kind  TokenKind
}

// NewToken ...
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

// Debug ...
func (t Token) Debug() {
	if t.isOneOfMany(Identifier, Number, String) {
		fmt.Printf("%s (%s)\n", t.Kind, t.Value)
	} else {
		fmt.Printf("%s\n", t.Kind)
	}
}
