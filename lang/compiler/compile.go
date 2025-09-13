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
	Output          string
}

const FilePerm = 0600

func Compile(tree ast.BlockStatement, target BuildTarget) {
	context := llvm.NewContext()
	defer context.Dispose()

	module := llvm.GlobalContext().NewModule("swa-main")
	defer module.Dispose()

	llvm.AddFunction(
		module,
		"printf",
		llvm.FunctionType(context.Int32Type(), []llvm.Type{llvm.PointerType(context.Int8Type(), 0)}, true),
	)

	builder := context.NewBuilder()
	defer builder.Dispose()
	ctx := ast.CompilerCtx{
		Context:           &context,
		Builder:           &builder,
		Module:            &module,
		SymbolTable:       map[string]llvm.Value{},
		StructSymbolTable: map[string]ast.StructSymbolTableEntry{},
		FuncSymbolTable:   map[string]llvm.Type{},
	}

	err, _ := tree.CompileLLVM(&ctx)
	if err != nil {
		panic(err)
	}

	llirFileName := fmt.Sprintf("%s.ll", target.Output)
	err = os.WriteFile(llirFileName, []byte(module.String()), FilePerm)
	if err != nil {
		panic(err)
	}

	asmFileName := fmt.Sprintf("%s.s", target.Output)
	cmd := exec.Command("llc-19", llirFileName, "-o", asmFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		err := fmt.Errorf("Error compiling IR %s", err)
		panic(err)
	}

	objFilename := fmt.Sprintf("%s.o", target.Output)
	objectCmd := exec.Command("clang-19", "-c", asmFileName, "-o", objFilename)
	if err := objectCmd.Run(); err != nil {
		panic(err)
	}

	exeFilename := fmt.Sprintf("%s.exe", target.Output)
	linkCmd := exec.Command("clang-19", objFilename, "-o", exeFilename)
	if err := linkCmd.Run(); err != nil {
		err2 := fmt.Errorf("Error durrng linking <%s>", err)
		panic(err2)
	}
}
