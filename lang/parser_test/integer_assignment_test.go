package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestIntegerAssigment(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
      nombre = 4;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "=", Kind: 9},
					Assignee: ast.SymbolExpression{Value: "nombre"},
					Value:    ast.NumberExpression{Value: 4},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithInfixOperator(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
      nombre = -4;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "=", Kind: 9},
					Assignee: ast.SymbolExpression{Value: "nombre"},
					Value: ast.PrefixExpression{
						Operator:        lexer.Token{Value: "-", Kind: 22},
						RightExpression: ast.NumberExpression{Value: 4},
					},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithAdditionToSelf(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
      nombre += 4;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "+=", Kind: 31},
					Assignee: ast.SymbolExpression{Value: "nombre"},
					Value:    ast.NumberExpression{Value: 4},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithMultiplicationToSelf(t *testing.T) {
	// result := parser.Parse(lexer.Tokenize(`
	//     // dialect=malinke;
	//     nombre *= 4;
	// `))

	// expected := ast.BlockStatement{
	// 	Body: []ast.Statement{
	// 		ast.ExpressionStatement{
	// 			Exp: ast.AssignmentExpression{
	// 				Operator: lexer.Token{Value: "+=", Kind: 31},
	// 				Assignee: ast.SymbolExpression{Value: "nombre"},
	// 				Value:    ast.NumberExpression{Value: 4},
	// 			},
	// 		},
	// 	},
	// }

	// assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithSubstractionToSelf(t *testing.T) {
	// result := parser.Parse(lexer.Tokenize(`
	//     // dialect=malinke;
	//     nombre *= 4;
	// `))

	// expected := ast.BlockStatement{
	// 	Body: []ast.Statement{
	// 		ast.ExpressionStatement{
	// 			Exp: ast.AssignmentExpression{
	// 				Operator: lexer.Token{Value: "+=", Kind: 31},
	// 				Assignee: ast.SymbolExpression{Value: "nombre"},
	// 				Value:    ast.NumberExpression{Value: 4},
	// 			},
	// 		},
	// 	},
	// }

	// assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithExpression(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
	     // dialect=malinke;
       const multiply = 45.2 + 5 * 4;
	 `))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "multiply",
				IsConstant: true,
				Value: ast.BinaryExpression{
					Left: ast.NumberExpression{Value: 45.2},
					Right: ast.BinaryExpression{
						Left:     ast.NumberExpression{Value: 5},
						Right:    ast.NumberExpression{Value: 4},
						Operator: lexer.Token{Value: "*", Kind: 32},
					},
					Operator: lexer.Token{Value: "+", Kind: 30},
				},
				ExplicitType: ast.Type(nil),
			},
		},
	}

	assert.Equal(t, result, expected)
}
