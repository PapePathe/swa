package parser

import (
	"swahili/lang/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrimaryExpression(t *testing.T) {
	t.Run("number expression", func(t *testing.T) {
		source := "dialect:french;10;"
		result := parseSourceCode(t, source)
		expected := ast.BlockStatement{
			Body: []ast.Statement{
				ast.ExpressionStatement{
					Exp: ast.NumberExpression{Value: 10.000000},
				},
			},
		}

		assert.Equal(t, expected, result)
	})
}
