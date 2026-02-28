package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"swahili/lang/assembly/generator"
	"swahili/lang/assembly/tac"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

type ASMCompilerRequest struct {
	Tree     *ast.BlockStatement
	Target   BuildTarget
	Dialect  lexer.Dialect
	Filename string
}

type ASMCompiler struct {
	req *ASMCompilerRequest
}

func NewASMCompiler(req *ASMCompilerRequest) *ASMCompiler {
	return &ASMCompiler{
		req: req,
	}
}

func (c *ASMCompiler) Run() error {
	gen := tac.NewTripleGenerator()

	err := c.req.Tree.Accept(gen)
	if err != nil {
		return err
	}

	linux := generator.NewLinux64AsmGen()
	err = gen.Gen(linux)
	if err != nil {
		return err
	}

	err = os.WriteFile(c.asmFileName(), []byte(linux.Code()), FilePerm)
	if err != nil {
		return err
	}

	clang := findCommand("clang-19", "clang")
	objectCmd := exec.Command(clang, "-Wall", "-c", c.asmFileName(), "-o", c.objectFileName(), "-fPIE")
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

func (c ASMCompiler) asmFileName() string {
	return fmt.Sprintf("%s.swa-asm", c.req.Target.Output)
}

func (c ASMCompiler) objectFileName() string {
	return fmt.Sprintf("%s.swa-obj", c.req.Target.Output)
}

func (c ASMCompiler) executableFileName() string {
	return fmt.Sprintf("%s.swa-exe", c.req.Target.Output)
}
