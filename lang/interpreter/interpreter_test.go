package interpreter

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"
)

func TestRun(t *testing.T) {
	tree := parser.Parse(lexer.Tokenize(`
  	dialect:english;
  	let number = 10;
 		`))

	scope := ast.NewScope(nil)
	Run(tree, scope)
}
