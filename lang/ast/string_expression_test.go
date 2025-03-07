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

import "testing"

func TestEvaluateString(t *testing.T) {
	s := NewScope(nil)
	sExpr := StringExpression{Value: "my string"}
	err, result := sExpr.Evaluate(s)

	if err != nil {
		t.Errorf("string evaluation should not error <%s>", err)
	}

	stringResult, ok := result.GetValue().(string)

	if !ok {
		t.Errorf("Expected value to be of type string")
	}

	if sExpr.Value != stringResult {
		t.Errorf("Expected %s to eq %s", sExpr, result)
	}
}
