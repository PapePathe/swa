package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestArrayInstantiation(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:malinke; 
			const numbers = []number{1,2,3,4,5,6};
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "numbers",
				IsConstant: true,
				Value: ast.ArrayInitializationExpression{
					Underlying: ast.SymbolType{Name: "number"},
					Contents: []ast.Expression{
						ast.NumberExpression{Value: 1},
						ast.NumberExpression{Value: 2},
						ast.NumberExpression{Value: 3},
						ast.NumberExpression{Value: 4},
						ast.NumberExpression{Value: 5},
						ast.NumberExpression{Value: 6},
					},
				},
				ExplicitType: ast.Type(nil),
			},
		},
	}

	assert.Equal(t, result, expected)
}
