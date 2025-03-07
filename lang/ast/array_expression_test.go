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
	"swahili/lang/values"
	"testing"
)

func TestArrayExpression(t *testing.T) {
	array := ArrayInitializationExpression{
		Underlying: NumberType{},
		Contents:   []Expression{NumberExpression{Value: 10}},
	}
	s := NewScope(nil)

	err, result := array.Evaluate(s)
	if err != nil {
		t.Errorf("Unexpected error <%s> while evaluating array initialization expression", err)
	}

	vals, ok := result.GetValue().([]values.Value)

	if !ok {
		t.Errorf("Evaluation should return an array value <%s>", vals)
	}

	if len(vals) != 1 {
		t.Errorf("Array should only have one value <%s>", vals)
	}
}
