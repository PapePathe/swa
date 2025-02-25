package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestStructInstantiation(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  		dialect:  malinke; 
      const r1: Rectangle = Rectangle {
        width: 10.2,
        height: 45.2 + 6 - 11 * 999,
      };
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "r1",
				IsConstant: true,
				Value: ast.StructInitializationExpression{
					Name: "Rectangle",
					Properties: map[string]ast.Expression{
						"height": ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left:     ast.NumberExpression{Value: 45.2},
								Right:    ast.NumberExpression{Value: 6},
								Operator: lexer.Token{Value: "+", Kind: 30},
							},
							Right: ast.BinaryExpression{
								Left:     ast.NumberExpression{Value: 11},
								Right:    ast.NumberExpression{Value: 999},
								Operator: lexer.Token{Value: "*", Kind: 32},
							},
							Operator: lexer.Token{Value: "-", Kind: 22},
						},
						"width": ast.NumberExpression{Value: 10.2},
					},
				},
				ExplicitType: ast.SymbolType{Name: "Rectangle"},
			},
		},
	}

	assert.Equal(t, result, expected)
}
