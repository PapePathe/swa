package compiler

import (
	"os"
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

type Package struct {
	D     lexer.Dialect
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
		dial := ctx.Dialect

		if pack.D != nil {
			dial = pack.D
		}

		tokens := lexer.TokenizeWithDialect(sourceCode, dial)

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
