package parser_test

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedASTForArrayDecl = ast.BlockStatement{
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

func TestArrayInstantiationEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:english; 
			const numbers = []number{1,2,3,4,5,6};
	`))
	assert.Equal(t, result, expectedASTForArrayDecl)
}

func TestArrayInstantiationFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:french; 
			constante numbers = []number{1,2,3,4,5,6};
	`))

	assert.Equal(t, result, expectedASTForArrayDecl)
}

func TestArrayInstantiationMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:malinke; 
			const numbers = []number{1,2,3,4,5,6};
	`))

	assert.Equal(t, result, expectedASTForArrayDecl)
}
