package compiler

import (
	"os"
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

type Package struct {
	Name  string
	Files []string
}

func (pack Package) Compile(ctx *ast.CompilerCtx) error {
	for _, source := range pack.Files {
		bytes, err := os.ReadFile(source)
		if err != nil {
			return err
		}
		sourceCode := string(bytes)
		tokens := lexer.TokenizeWithDialect(sourceCode, ctx.Dialect)
		tree, _, err := parser.Parse(tokens, true)
		if err != nil {
			return err
		}

		err, _ = tree.CompileLLVM(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
