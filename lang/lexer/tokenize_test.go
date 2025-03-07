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

import "testing"

func TestTokenize(t *testing.T) {
	result := Tokenize(`
  	dialect:french;
    si(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } sinon {
    	width += 300;
    }
	`)

	expected := []Token{
		{Kind: DialectDeclaration, Value: "dialect"},
		{Kind: Colon, Value: ":"},
		{Kind: Identifier, Value: "french"},
		{Kind: SemiColon, Value: ";"},
		{Kind: KeywordIf, Value: "si"},
		{Kind: OpenParen, Value: "("},
		{Kind: Identifier, Value: "x"},
		{Kind: GreaterThan, Value: ">"},
		{Kind: Number, Value: "0"},
		{Kind: CloseParen, Value: ")"},
		{Kind: OpenCurly, Value: "{"},
		{Kind: Identifier, Value: "width"},
		{Kind: Assignment, Value: "="},
		{Kind: Number, Value: "100"},
		{Kind: SemiColon, Value: ";"},
		{Kind: Identifier, Value: "height"},
		{Kind: Assignment, Value: "="},
		{Kind: Number, Value: "100"},
		{Kind: Plus, Value: "+"},
		{Kind: Number, Value: "400"},
		{Kind: Minus, Value: "-"},
		{Kind: Identifier, Value: "width"},
		{Kind: SemiColon, Value: ";"},
		{Kind: CloseCurly, Value: "}"},
		{Kind: KeywordElse, Value: "sinon"},
		{Kind: OpenCurly, Value: "{"},
		{Kind: Identifier, Value: "width"},
		{Kind: PlusEquals, Value: "+="},
		{Kind: Number, Value: "300"},
		{Kind: SemiColon, Value: ";"},
		{Kind: CloseCurly, Value: "}"},
		{Kind: EOF, Value: "EOF"},
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %s to eq %s at index %d", expected[i], v, i)
		}
	}
}
