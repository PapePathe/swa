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

	"github.com/stretchr/testify/assert"
)

func TestEvaluateIntegerVariableDeclarationStatement(t *testing.T) {
	scope := NewScope(nil)
	statement := VarDeclarationStatement{
		Name:       "x",
		IsConstant: false,
		Value:      NumberExpression{Value: 10},
	}

	err, val := statement.Evaluate(scope)

	if err != nil {
		t.Errorf("Evaluating var decl should not error. current: %s", err)
	}

	if val != nil {
		t.Errorf("return value of var decl statement should be nil")
	}

	content, ok := scope.Variables["x"]

	if !ok {
		t.Errorf("Variable X should be defined")
	}

	expectedContent := values.NumberValue{Value: 10}

	if content != expectedContent {
		t.Errorf("%s should be equal to %s", expectedContent, content)
	}
}

func TestEvaluateStringVariableDeclarationStatement(t *testing.T) {
	scope := NewScope(nil)
	statement := VarDeclarationStatement{
		Name:       "string_var",
		IsConstant: false,
		Value:      StringExpression{Value: "a simple string"},
	}

	err, val := statement.Evaluate(scope)

	if err != nil {
		t.Errorf("Evaluating var decl should not error. current: %s", err)
	}

	if val != nil {
		t.Errorf("return value of var decl statement should be nil")
	}

	content, ok := scope.Variables["string_var"]

	if !ok {
		t.Errorf("Variable string_var should be defined")
	}

	expectedContent := values.StringValue{Value: "a simple string"}

	if content != expectedContent {
		t.Errorf("%s should be equal to %s", expectedContent, content)
	}
}

func TestEvaluatArrayVariableDeclarationStatement(t *testing.T) {
	scope := NewScope(nil)
	statement := VarDeclarationStatement{
		Name:       "array",
		IsConstant: false,
		Value: ArrayInitializationExpression{
			Contents: []Expression{StringExpression{Value: "a simple string"}},
		},
	}

	err, val := statement.Evaluate(scope)

	if err != nil {
		t.Errorf("Evaluating var decl should not error. current: %s", err)
	}

	if val != nil {
		t.Errorf("return value of var decl statement should be nil")
	}

	content, ok := scope.Variables["array"]

	if !ok {
		t.Errorf("Variable string_var should be defined")
	}

	expectedContent := values.ArrayValue{Values: []values.Value{values.StringValue{Value: "a simple string"}}}

	assert.Equal(t, content, expectedContent)
}
