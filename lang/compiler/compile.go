package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"swahili/lang/ast"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type BuildTarget struct {
	OperatingSystem string
	Architecture    string
	Output          string
}

const FilePerm = 0600

func Compile(tree ast.BlockStatement, target BuildTarget, dialect lexer.Dialect) {
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
		Dialect:           dialect,
		SymbolTable:       map[string]ast.SymbolTableEntry{},
		StructSymbolTable: map[string]ast.StructSymbolTableEntry{},
		FuncSymbolTable:   map[string]llvm.Type{},
		ArraysSymbolTable: map[string]ast.ArraySymbolTableEntry{},
	}

	err, _ := tree.CompileLLVM(&ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	llirFileName := fmt.Sprintf("%s.ll", target.Output)

	err = os.WriteFile(llirFileName, []byte(module.String()), FilePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	asmFileName := fmt.Sprintf("%s.s", target.Output)

	if err := compileToAssembler(llirFileName, asmFileName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	objFilename := fmt.Sprintf("%s.o", target.Output)
	if err := compileToObject(asmFileName, objFilename); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exeFilename := fmt.Sprintf("%s.exe", target.Output)
	if err := compileToExecutable(objFilename, exeFilename); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func compileToAssembler(llirFileName string, assemblerFilename string) error {
	cmd := exec.Command("llc-19", llirFileName, "-o", assemblerFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error compiling IR %w", err)
	}

	return nil
}

func compileToObject(assemblerFilename string, objectFilename string) error {
	objectCmd := exec.Command("clang-19", "-c", assemblerFilename, "-o", objectFilename)
	objectCmd.Stdout = os.Stdout
	objectCmd.Stderr = os.Stderr

	if err := objectCmd.Run(); err != nil {
		return fmt.Errorf("Error durrng object creation <%w>", err)
	}
	return nil
}

func compileToExecutable(objectFileName string, executableFileName string) error {
	linkCmd := exec.Command("clang-19", objectFileName, "-o", executableFileName)
	linkCmd.Stdout = os.Stdout
	linkCmd.Stderr = os.Stderr

	if err := linkCmd.Run(); err != nil {
		return fmt.Errorf("Error durrng linking <%w>", err)
	}
	return nil
}
