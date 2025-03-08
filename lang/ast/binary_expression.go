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

package ast

import (
	"fmt"
	"swahili/lang/lexer"
	"swahili/lang/values"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (BinaryExpression) expression() {}
func (be BinaryExpression) Evaluate(s *Scope) (error, values.Value) {
	fmt.Println("Evaluating BinaryExpression", be)

	_, left := be.Left.Evaluate(s)
	_, right := be.Right.Evaluate(s)

	leftVal, _ := left.GetValue().(float64)
	rightVal, _ := right.GetValue().(float64)

	switch be.Operator.Kind {
	case lexer.GreaterThan:
		return nil, values.BooleaValue{Value: leftVal > rightVal}
	}

	return nil, nil
}
