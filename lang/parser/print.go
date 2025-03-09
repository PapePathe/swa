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
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrintStatement(p *Parser) ast.Statement {
	values := []ast.Expression{}

	p.expect(lexer.Print)

	p.expect(lexer.OpenParen)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		switch p.currentToken().Kind {
		case lexer.String:
			values = append(values, ast.StringExpression{Value: p.expect(lexer.String).Value})
		case lexer.Identifier:
			values = append(values, ast.SymbolExpression{Value: p.expect(lexer.Identifier).Value})
		default:
			panic(fmt.Sprintf("Token %s not supported in print statement", p.currentToken().Kind))
		}

		if p.currentToken().Kind == lexer.Comma {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseParen)
	p.expect(lexer.SemiColon)

	return ast.PrintStatetement{}
}
