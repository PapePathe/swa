package compiler

import (
	"fmt"
	"os"
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"sync"

	"tinygo.org/x/go-llvm"
)

type Package struct {
	Name  string
	Files []string
}

func (pack Package) Compile() error {
	context := llvm.NewContext()
	builder := context.NewBuilder()
	module := llvm.GlobalContext().NewModule(pack.Name)
	wg := sync.WaitGroup{}

	defer context.Dispose()
	defer module.Dispose()
	defer builder.Dispose()

	for _, source := range pack.Files {
		wg.Add(1)
		go func() {

			bytes, err := os.ReadFile(source)
			if err != nil {
				panic(err)
			}
			sourceCode := string(bytes)
			tokens, dialect := lexer.Tokenize(sourceCode)
			tree, err := parser.Parse(tokens)
			if err != nil {
				panic(err)
			}

			ctx := ast.NewCompilerContext(
				&context,
				&builder,
				&module,
				dialect,
				nil,
			)
			err, _ = tree.CompileLLVM(ctx)
			if err != nil {
				panic(err)
			}

			fmt.Println(module.String())

			wg.Done()

		}()
	}

	wg.Wait()

	return nil
}
