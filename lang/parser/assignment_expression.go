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

import "swahili/lang/ast"

func ParseAssignmentExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	rightHandSide := parseExpression(p, DefaultBindingPower)

	return ast.AssignmentExpression{
		Operator: operatorToken,
		Value:    rightHandSide,
		Assignee: left,
	}
}
