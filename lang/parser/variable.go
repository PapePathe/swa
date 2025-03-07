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

// ParseVarDeclarationStatement ...
func ParseVarDeclarationStatement(p *Parser) ast.Statement {
	var explicitType ast.Type

	var assigedValue ast.Expression

	isConstant := p.advance().Kind == lexer.Const
	errStr := "Inside variable declaration expected to find variable name"
	variableName := p.expectError(lexer.Identifier, errStr).Value

	if p.currentToken().Kind == lexer.Colon {
		p.advance()
		explicitType = parseType(p, DefaultBindingPower)
	}

	if p.currentToken().Kind != lexer.SemiColon {
		p.expect(lexer.Assignment)
		assigedValue = parseExpression(p, Assignment)
	} else if explicitType == nil {
		panic("Missing either right hand side in var declaration or exlicit type")
	}

	if isConstant && assigedValue == nil {
		panic("Cannot define constant wihtout a value")
	}

	p.expect(lexer.SemiColon)

	return ast.VarDeclarationStatement{
		IsConstant:   isConstant,
		Value:        assigedValue,
		Name:         variableName,
		ExplicitType: explicitType,
	}
}
