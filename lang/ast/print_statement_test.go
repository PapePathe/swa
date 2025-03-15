package ast

import (
	"errors"
	"swahili/lang/values"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluatePrintStatement(t *testing.T) {
	t.Run("print static string", func(t *testing.T) {
		s := NewScope(nil)
		printExpr := PrintStatetement{
			Values: []Expression{StringExpression{Value: "Hola"}},
		}
		expectedValue := values.StringValue{Value: "Hola"}
		err, v := printExpr.Evaluate(s)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, v)
	})
	t.Run("print number", func(t *testing.T) {
		s := NewScope(nil)
		printExpr := PrintStatetement{
			Values: []Expression{NumberExpression{Value: 19}},
		}
		expectedValue := values.StringValue{Value: "19.000000"}
		err, v := printExpr.Evaluate(s)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, v)
	})
	t.Run("print symbol", func(t *testing.T) {
		s := NewScope(nil)
		s.Set("x", values.StringValue{Value: "Hola"})
		printExpr := PrintStatetement{
			Values: []Expression{SymbolExpression{Value: "x"}},
		}
		expectedValue := values.StringValue{Value: "Hola"}
		err, v := printExpr.Evaluate(s)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, v)
	})
	t.Run("errors when trying to print undefined symbol", func(t *testing.T) {
		s := NewScope(nil)
		printExpr := PrintStatetement{
			Values: []Expression{SymbolExpression{Value: "x"}},
		}
		err, _ := printExpr.Evaluate(s)

		assert.Equal(t, errors.New("Variable <x> does not exist"), err)
	})
}
