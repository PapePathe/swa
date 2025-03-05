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
