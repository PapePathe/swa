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
