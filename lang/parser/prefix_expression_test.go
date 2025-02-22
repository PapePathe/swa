package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func TestParsePrefixExpression(t *testing.T) {
	st := Parse(lexer.Tokenize("-44;"))
	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.PrefixExpression{
					Operator:        lexer.Token{Value: "-", Kind: 22},
					RightExpression: ast.NumberExpression{Value: 44},
				},
			},
		},
	}

	assert.Equal(t, st, expected)
}
