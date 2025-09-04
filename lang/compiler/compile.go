package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type BuildTarget struct {
	OperatingSystem string
	Architecture    string
}

const FilePerm = 0600

func Compile(tree ast.BlockStatement, target BuildTarget) {
	//	m := ir.NewModule()
	// f := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	// f.Sig.Variadic = true

	context := llvm.NewContext()
	defer context.Dispose() // Clean up when we're done.

	module := llvm.GlobalContext().NewModule("my_module")
	defer module.Dispose()

	llvm.AddFunction(
		module,
		"printf",
		llvm.FunctionType(context.Int32Type(), []llvm.Type{llvm.PointerType(context.Int8Type(), 0)}, true),
	)

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

	err = os.WriteFile("start.ll", []byte(module.String()), FilePerm)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("llc-19", "start.ll", "-o", "start.s")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		err := fmt.Errorf("Error compiling IR %s", err)
		panic(err)
	}

	objectCmd := exec.Command("clang-19", "-c", "start.s", "-o", "start.o")
	if err := objectCmd.Run(); err != nil {
		panic(err)
	}

	linkCmd := exec.Command("clang-19", "start.o", "-o", "start.exe")
	if err := linkCmd.Run(); err != nil {
		err2 := fmt.Errorf("Error durrng linking <%s>", err)
		panic(err2)
	}
}
