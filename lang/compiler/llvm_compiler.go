package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"swahili/lang/ast"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type LLVMCompiler struct {
	req     *LLVMCompilerRequest
	passes  []ast.CodeGenerator
	context *CompilerCtx
}

type LLVMCompilerRequest struct {
	Tree     *ast.BlockStatement
	Target   BuildTarget
	Dialect  lexer.Dialect
	Filename string
}

func NewLLVMCompiler(req *LLVMCompilerRequest) *LLVMCompiler {
	context := llvm.GlobalContext()
	module := context.NewModule(req.Filename)
	builder := context.NewBuilder()
	ctx := NewCompilerContext(&context, &builder, &module, req.Dialect, nil)

	return &LLVMCompiler{
		req: req,
		passes: []ast.CodeGenerator{
			NewLLVMGenerator(ctx),
			NewLLVMTypeChecker(ctx),
		},
		context: ctx,
	}
}

func (c *LLVMCompiler) Run() error {
	defer c.context.Module.Dispose()
	defer c.context.Builder.Dispose()

	printArgTypes := []llvm.Type{llvm.PointerType(c.context.Context.Int8Type(), 0)}
	printfFuncType := llvm.FunctionType(c.context.Context.Int32Type(), printArgTypes, true)
	printfFunc := llvm.AddFunction(*c.context.Module, "printf", printfFuncType)
	printfFunc.SetLinkage(llvm.ExternalLinkage)

	strcmpArgTypes := []llvm.Type{
		llvm.PointerType(c.context.Context.Int8Type(), 0),
		llvm.PointerType(c.context.Context.Int8Type(), 0),
	}
	strcmpFuncType := llvm.FunctionType(c.context.Context.Int32Type(), strcmpArgTypes, false)
	strcmpFunc := llvm.AddFunction(*c.context.Module, "strcmp", strcmpFuncType)
	strcmpFunc.SetLinkage(llvm.ExternalLinkage)

	for _, v := range c.passes {
		err := c.req.Tree.Accept(v)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(c.llirFileName(), []byte(c.context.Module.String()), FilePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	count := c.context.MainFuncOccurrences
	switch count {
	case 1:
	// we good
	case 0:
		key := "LLVMCompiler.MissingProgramEntrypoint"

		return c.context.Dialect.Error(key)
	default:
		key := "LLVMCompiler.TooManyProgramEntrypoints"

		return c.context.Dialect.Error(key, count)
	}

	err = llvm.VerifyModule(*c.context.Module, llvm.ReturnStatusAction)
	if err != nil {
		return err
	}

	if c.context.Debugging {
		fmt.Println("Started llir compilation")
	}

	llc := findCommand("llc-19", "llc")
	cmd := exec.Command(llc, c.llirFileName(), "-o", c.asmFileName(), "-relocation-model=pic")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error compiling IR %w", err)
	}

	if c.context.Debugging {
		fmt.Println("Finished llir compilation")
	}

	clang := findCommand("clang-19", "clang")
	objectCmd := exec.Command(clang, "-c", c.asmFileName(), "-o", c.objectFileName(), "-fPIE")
	objectCmd.Stdout = os.Stdout
	objectCmd.Stderr = os.Stderr

	err = objectCmd.Run()
	if err != nil {
		return fmt.Errorf("Error durrng object creation <%w>", err)
	}

	linkArgs := []string{c.objectFileName(), "-o", c.executableFileName()}
	linkCmd := exec.Command(clang, linkArgs...)
	linkCmd.Stdout = os.Stdout
	linkCmd.Stderr = os.Stderr

	err = linkCmd.Run()
	if err != nil {
		return fmt.Errorf("Error durrng linking <%w>", err)
	}

	return nil
}

func (c LLVMCompiler) llirFileName() string {
	return fmt.Sprintf("%s.ll", c.req.Target.Output)
}

func (c LLVMCompiler) asmFileName() string {
	return fmt.Sprintf("%s.s", c.req.Target.Output)
}

func (c LLVMCompiler) objectFileName() string {
	return fmt.Sprintf("%s.o", c.req.Target.Output)
}

func (c LLVMCompiler) executableFileName() string {
	return fmt.Sprintf("%s.exe", c.req.Target.Output)
}
