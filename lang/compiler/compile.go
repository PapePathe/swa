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

func Compile(tree ast.BlockStatement, target BuildTarget, dialect lexer.Dialect, experimental bool) {
	var context llvm.Context
	var module llvm.Module
	var builder llvm.Builder

	if experimental {
		context = llvm.NewContext()
		module = context.NewModule("swa-main")
		builder = context.NewBuilder()
	} else {
		context = llvm.GlobalContext()
		module = context.NewModule("swa-main")
		builder = context.NewBuilder()
	}

	if experimental {
		defer context.Dispose()
	}
	defer module.Dispose()
	defer builder.Dispose()

	printArgTypes := []llvm.Type{llvm.PointerType(context.Int8Type(), 0)}
	printfFuncType := llvm.FunctionType(context.Int32Type(), printArgTypes, true)
	printfFunc := llvm.AddFunction(module, "printf", printfFuncType)
	printfFunc.SetLinkage(llvm.ExternalLinkage)

	ctx := ast.NewCompilerContext(
		&context,
		&builder,
		&module,
		dialect,
		nil,
	)

	var err error
	if experimental {
		generator := NewLLVMGenerator(ctx)
		err = tree.Accept(generator)
	} else {
		err, _ = tree.CompileLLVM(ctx)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = llvm.VerifyModule(*ctx.Module, llvm.ReturnStatusAction)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	llirFileName := fmt.Sprintf("%s.ll", target.Output)

	err = os.WriteFile(llirFileName, []byte(module.String()), FilePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	asmFileName := fmt.Sprintf("%s.s", target.Output)

	err = compileToAssembler(llirFileName, asmFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	objFilename := fmt.Sprintf("%s.o", target.Output)

	err = compileToObject(asmFileName, objFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exeFilename := fmt.Sprintf("%s.exe", target.Output)

	err = compileToExecutable(objFilename, exeFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func compileToAssembler(llirFileName string, assemblerFilename string) error {
	llc := findCommand("llc-19", "llc")
	cmd := exec.Command(llc, llirFileName, "-o", assemblerFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error compiling IR %w", err)
	}

	return nil
}

func compileToObject(assemblerFilename string, objectFilename string) error {
	clang := findCommand("clang-19", "clang")
	objectCmd := exec.Command(clang, "-c", assemblerFilename, "-o", objectFilename)
	objectCmd.Stdout = os.Stdout
	objectCmd.Stderr = os.Stderr

	err := objectCmd.Run()
	if err != nil {
		return fmt.Errorf("Error durrng object creation <%w>", err)
	}

	return nil
}

func compileToExecutable(objectFileName string, executableFileName string) error {
	clang := findCommand("clang-19", "clang")
	linkCmd := exec.Command(clang, objectFileName, "-o", executableFileName, "-no-pie")
	linkCmd.Stdout = os.Stdout
	linkCmd.Stderr = os.Stderr

	err := linkCmd.Run()
	if err != nil {
		return fmt.Errorf("Error durrng linking <%w>", err)
	}

	return nil
}

func findCommand(candidates ...string) string {
	for _, cmd := range candidates {
		_, err := exec.LookPath(cmd)
		if err == nil {
			return cmd
		}
	}

	return candidates[0]
}
