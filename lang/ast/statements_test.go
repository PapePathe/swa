package ast

import (
	"swahili/lang/values"
	"testing"
)

func TestEvaluateVariableDeclarationStatement(t *testing.T) {
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
