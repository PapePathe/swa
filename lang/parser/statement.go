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

// ParseStatement ...
func ParseStatement(p *Parser) ast.Statement {
	statementFn, exists := statementLookup[p.currentToken().Kind]

	if exists {
		return statementFn(p)
	}

	expression := parseExpression(p, DefaultBindingPower)
	p.expect(lexer.SemiColon)

	return ast.ExpressionStatement{
		Exp: expression,
	}
}

func ParseStructDeclarationStatement(p *Parser) ast.Statement {
	p.expect(lexer.Struct)
	structName := p.expect(lexer.Identifier).Value

	propertes := map[string]ast.StructProperty{}

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		var propertyName string

		if p.currentToken().Kind == lexer.Identifier {
			propertyName = p.expect(lexer.Identifier).Value
			p.expectError(lexer.Colon, "Expected to find colon following struct property name")
			propType := parseType(p, DefaultBindingPower)
			p.expect(lexer.Comma)

			if _, exists := propertes[propertyName]; exists {
				panic(fmt.Sprintf("property %s has already been defined", propertyName))
			}

			propertes[propertyName] = ast.StructProperty{
				PropType: propType,
			}

			continue
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.StructDeclarationStatement{
		Name:       structName,
		Properties: propertes,
	}
}
