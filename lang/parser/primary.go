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
	"strconv"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParsePrimaryExpression ...
func ParsePrimaryExpression(p *Parser) ast.Expression {
	switch p.currentToken().Kind {
	case lexer.Number:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)

		return ast.NumberExpression{
			Value: number,
		}
	case lexer.String:
		value := p.advance().Value

		return ast.StringExpression{
			Value: value[1 : len(value)-1],
		}
	case lexer.Identifier:
		return ast.SymbolExpression{
			Value: p.advance().Value,
		}

	default:
		panic(fmt.Sprintf("Cannot create PrimaryExpression from %s", p.currentToken().Kind))
	}
}
