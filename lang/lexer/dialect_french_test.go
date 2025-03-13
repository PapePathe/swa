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
	"regexp"
	"testing"
)

func TestReservedFrench(t *testing.T) {
	expected := map[string]TokenKind{
		"si":        KeywordIf,
		"sinon":     KeywordElse,
		"structure": Struct,
		"variable":  Let,
		"constante": Const,
		"entier":    TypeInt,
	}

	French := French{}
	reserved := French.Reserved()

	if len(reserved) != len(expected) {
		t.Errorf("expected length of reserved keywords to match - actual: %d expected: eq %d", len(reserved), len(expected))
	}

	for name, kind := range expected {
		if kind != reserved[name] {
			t.Errorf("expected French reserved word for (%s) to be (%s)", kind, name)
		}
	}
}

func TestPatternsFrench(t *testing.T) {
	french := French{}
	patterns := french.Patterns()

	expectedPatternCount := 36
	if len(french.Patterns()) != expectedPatternCount {
		t.Errorf("Expected %d patterns, got %d", expectedPatternCount, len(patterns))
	}

	tests := []struct {
		pattern         *regexp.Regexp
		expectedPattern string
		handler         regexHandler
	}{
		{patterns[0].regex, `\s+`, skipHandler},
		{patterns[1].regex, `dialect`, defaultHandler(DialectDeclaration, "dialect")},
		{patterns[2].regex, `entier`, defaultHandler(TypeInt, "int")},
		{patterns[3].regex, `sinon`, defaultHandler(KeywordElse, "sinon")},
		{patterns[4].regex, `si`, defaultHandler(KeywordIf, "si")},
		{patterns[5].regex, `afficher`, defaultHandler(Struct, "afficher")},
		{patterns[6].regex, `structure`, defaultHandler(Struct, "structure")},
		{patterns[7].regex, `\/\/.*`, commentHandler},
		{patterns[8].regex, `"[^"]*"`, stringHandler},
		{patterns[9].regex, `[0-9]+(\.[0-9]+)?`, numberHandler},
		{patterns[10].regex, `[a-zA-Z_][a-zA-Z0-9_]*`, symbolHandler},
		{patterns[11].regex, `\[`, defaultHandler(OpenBracket, "[")},
		{patterns[12].regex, `\]`, defaultHandler(CloseBracket, "]")},
		{patterns[13].regex, `\{`, defaultHandler(OpenCurly, "{")},
		{patterns[14].regex, `\}`, defaultHandler(CloseCurly, "}")},
		{patterns[15].regex, `\(`, defaultHandler(OpenParen, "(")},
		{patterns[16].regex, `\)`, defaultHandler(CloseParen, ")")},
		{patterns[17].regex, `!=`, defaultHandler(NotEquals, "!=")},
		{patterns[18].regex, `\+=`, defaultHandler(PlusEquals, "+=")},
		{patterns[19].regex, `==`, defaultHandler(Equals, "==")},
		{patterns[20].regex, `=`, defaultHandler(Assignment, "=")},
		{patterns[21].regex, `!`, defaultHandler(Not, "!")},
		{patterns[22].regex, `<=`, defaultHandler(LessThanEquals, "<=")},
		{patterns[23].regex, `<`, defaultHandler(LessThan, "<")},
		{patterns[24].regex, `>=`, defaultHandler(GreaterThanEquals, ">=")},
		{patterns[25].regex, `>`, defaultHandler(GreaterThan, ">")},
		{patterns[26].regex, `\|\|`, defaultHandler(Or, "||")},
		{patterns[27].regex, `&&`, defaultHandler(And, "&&")},
		{patterns[28].regex, `\.`, defaultHandler(Dot, ".")},
		{patterns[29].regex, `;`, defaultHandler(SemiColon, ";")},
		{patterns[30].regex, `:`, defaultHandler(Colon, ":")},
		{patterns[31].regex, `,`, defaultHandler(Comma, ",")},
		{patterns[32].regex, `\+`, defaultHandler(Plus, "+")},
		{patterns[33].regex, `-`, defaultHandler(Minus, "-")},
		{patterns[34].regex, `/`, defaultHandler(Divide, "/")},
		{patterns[35].regex, `\*`, defaultHandler(Star, "*")},
	}

	for _, tt := range tests {
		t.Run(tt.expectedPattern, func(t *testing.T) {
			if tt.pattern.String() != tt.expectedPattern {
				t.Errorf("Expected pattern %s, got %s", tt.expectedPattern, tt.pattern.String())
			}

			if tt.handler == nil {
				t.Error("Expected handler to be non-nil")
			}
		})
	}
}
