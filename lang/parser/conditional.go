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

func ParseBlockStatement(p *Parser) ast.BlockStatement {
	body := []ast.Statement{}

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		body = append(body, ParseStatement(p))
	}

	p.expect(lexer.CloseCurly)

	return ast.BlockStatement{
		Body: body,
	}
}

func ParseConditionalExpression(p *Parser) ast.Statement {
	failBlock := ast.BlockStatement{}

	p.expect(lexer.KeywordIf)
	p.expect(lexer.OpenParen)

	condition := parseExpression(p, DefaultBindingPower)

	p.expect(lexer.CloseParen)

	successBlock := ParseBlockStatement(p)

	if p.currentToken().Kind == lexer.KeywordElse {
		p.expect(lexer.KeywordElse)

		failBlock = ParseBlockStatement(p)
	}

	return ast.ConditionalStatetement{
		Condition: condition,
		Success:   successBlock,
		Failure:   failBlock,
	}
}
