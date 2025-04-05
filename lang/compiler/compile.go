package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"tinygo.org/x/go-llvm"
)

type BuildTarget struct {
	OperatingSystem string
	Architecture    string
}

const FilePerm = 0600

func Compile(tree ast.BlockStatement, target BuildTarget) {
	m := ir.NewModule()
	f := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	f.Sig.Variadic = true

	context := llvm.NewContext()
	defer context.Dispose() // Clean up when we're done.

	// Create a new module.  A module is a container for LLVM IR.
	module := llvm.GlobalContext().NewModule("my_module")
	defer module.Dispose()

	builder := context.NewBuilder()
	defer builder.Dispose()
	ctx := ast.CompilerCtx{
		Context:     &context,
		Builder:     &builder,
		Module:      &module,
		SymbolTable: map[string]llvm.Value{},
	}

	err, _ := tree.CompileLLVM(&ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(module.String())

	// err = os.WriteFile("./tmp/start.ll", []byte(m.String()), FilePerm)
	// if err != nil {
	// 	panic(err)
	// }

	// cmd := exec.Command("llc", "./tmp/start.ll", "-o", "./tmp/start.s")
	// if err := cmd.Run(); err != nil {
	// 	panic(err)
	// }

	// objectCmd := exec.Command("clang", "-c", "./tmp/start.s", "-o", "./tmp/start.o")
	// if err := objectCmd.Run(); err != nil {
	// 	panic(err)
	// }

	// linkCmd := exec.Command("clang", "./tmp/start.o", "-o", "./tmp/start.exe")
	// if err := linkCmd.Run(); err != nil {
	// 	err2 := fmt.Errorf("Error during linking <%s>", err)
	// 	panic(err2)
	// }
}
