package parser_test

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedASTForDualCondition = ast.BlockStatement{
	Body: []ast.Statement{
		ast.ConditionalStatetement{
			Condition: ast.BinaryExpression{
				Left:     ast.SymbolExpression{Value: "x"},
				Right:    ast.NumberExpression{Value: 0},
				Operator: lexer.Token{Value: ">", Kind: 18},
			},
			Success: ast.BlockStatement{
				Body: []ast.Statement{
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9},
							Assignee: ast.SymbolExpression{Value: "width"},
							Value:    ast.NumberExpression{Value: 100},
						},
					},
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9},
							Assignee: ast.SymbolExpression{Value: "height"},
							Value: ast.BinaryExpression{
								Left: ast.BinaryExpression{
									Left:     ast.NumberExpression{Value: 100},
									Right:    ast.NumberExpression{Value: 400},
									Operator: lexer.Token{Value: "+", Kind: 30},
								},
								Right:    ast.SymbolExpression{Value: "width"},
								Operator: lexer.Token{Value: "-", Kind: 22},
							},
						},
					},
				},
			},
			Failure: ast.BlockStatement{
				Body: []ast.Statement{
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "+=", Kind: 31},
							Assignee: ast.SymbolExpression{Value: "width"},
							Value:    ast.NumberExpression{Value: 300},
						},
					},
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "+=", Kind: 31},
							Assignee: ast.SymbolExpression{Value: "height"},
							Value: ast.BinaryExpression{
								Left:     ast.SymbolExpression{Value: "rand"},
								Right:    ast.NumberExpression{Value: 14},
								Operator: lexer.Token{Value: "/", Kind: 15},
							},
						},
					},
				},
			},
		},
	},
}

//func TestDualBranchConditionalFrench(t *testing.T) {
//	result := parser.Parse(lexer.Tokenize(`
//  	dialect:french;
//    si(x>0) {
//      width = 100;
//      height = 100 + 400 - width;
//    } sinon {
//    	width += 300;
//    	height += rand / 14;
//    }
//	`))
//
//	assert.Equal(t, expectedASTForDualCondition, result)
//}

func TestDualBranchConditionalMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  	dialect:malinke;
    ni(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } nii {
    	width += 300;
    	height += rand / 14;
    }
	`))

	assert.Equal(t, expectedASTForDualCondition, result)
}

func TestDualBranchConditionalEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  	dialect:english;
    if(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } else {
    	width += 300;
    	height += rand / 14;
    }
	`))

	assert.Equal(t, expectedASTForDualCondition, result)
}

var expectedAST = ast.BlockStatement{
	Body: []ast.Statement{
		ast.ConditionalStatetement{
			Condition: ast.BinaryExpression{
				Left: ast.SymbolExpression{Value: "x"},
				Right: ast.NumberExpression{
					Value: 0,
				}, Operator: lexer.Token{Value: ">", Kind: 18},
			},
			Success: ast.BlockStatement{
				Body: []ast.Statement{
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9},
							Assignee: ast.MemberExpression{
								Object:   ast.SymbolExpression{Value: "r1"},
								Property: ast.SymbolExpression{Value: "width"},
								Computed: false,
							},
							Value: ast.NumberExpression{Value: 100},
						},
					},
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9},
							Assignee: ast.MemberExpression{
								Object: ast.SymbolExpression{Value: "r1"},
								Property: ast.SymbolExpression{
									Value: "height",
								},
								Computed: false,
							},
							Value: ast.BinaryExpression{
								Left: ast.BinaryExpression{
									Left: ast.NumberExpression{
										Value: 100,
									}, Right: ast.NumberExpression{Value: 400},
									Operator: lexer.Token{Value: "+", Kind: 30},
								},
								Right: ast.MemberExpression{
									Object: ast.SymbolExpression{
										Value: "r1",
									}, Property: ast.SymbolExpression{Value: "width"}, Computed: false,
								},
								Operator: lexer.Token{Value: "-", Kind: 22},
							},
						},
					},
				},
			},
			Failure: ast.BlockStatement{
				Body: []ast.Statement(nil),
			},
		},
	},
}

func TestSingleBranchConditionalFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
 		 dialect:french;
 		 si(x>0) {
 		   r1.width = 100;
 		   r1.height = 100 + 400 - r1.width;
 		 }
	`))

	assert.Equal(t, expectedAST, result)
}

func TestSingleBranchConditionalMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
 		 dialect:malinke;
 		 ni(x>0) {
 		   r1.width = 100;
 		   r1.height = 100 + 400 - r1.width;
 		 }
	`))

	assert.Equal(t, expectedAST, result)
}

func TestSingleBranchConditionalEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
 		 dialect:english;
 		 if(x>0) {
 		   r1.width = 100;
 		   r1.height = 100 + 400 - r1.width;
 		 }
	`))

	assert.Equal(t, expectedAST, result)
}
