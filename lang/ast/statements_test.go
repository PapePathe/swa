package ast

import (
	"swahili/lang/values"
	"testing"
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

	if err, _ = statement.Evaluate(scope); err == nil {
		t.Errorf("declaring an existing variable should error with <%s>", err)
	}
}
