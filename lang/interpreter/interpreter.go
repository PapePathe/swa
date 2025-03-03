package interpreter

import (
	"fmt"
	"swahili/lang/ast"
)

func Run(program ast.BlockStatement, globalScope *ast.Scope) {
	for _, v := range program.Body {
		result, _ := v.Evaluate(globalScope)

		fmt.Println(result)
	}
}
