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

package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseArrayInitialization(p *Parser) ast.Expression {
	var underlying ast.Type

	contents := []ast.Expression{}

	p.expect(lexer.OpenBracket)
	p.expect(lexer.CloseBracket)

	underlying = parseType(p, DefaultBindingPower)
	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		expr := parseExpression(p, Logical)
		contents = append(contents, expr)

		if p.currentToken().Kind != lexer.CloseCurly {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.ArrayInitializationExpression{
		Underlying: underlying,
		Contents:   contents,
	}
}
