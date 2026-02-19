package compiler

import (
	"fmt"
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
	// passes []ast.CodeGenerator
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

	fmt.Println("header:")

	for i, v := range gen.Insts {
		var ArgOne string
		if v.ArgOne != nil {
			ArgOne = v.ArgOne.InstructionArg()
		}

		var ArgTwo string
		if v.ArgTwo != nil {
			ArgTwo = v.ArgTwo.InstructionArg()
		}

		fmt.Printf(" %5d %20s  %10s  %10s\n", i, v.Operation, ArgOne, ArgTwo)
	}

	fmt.Println()

	for _, proc := range gen.Procs {
		fmt.Printf("%s @%s( ", proc.Ret.String(), proc.Name)

		for _, arg := range proc.Args {
			fmt.Printf("%s:%s ", arg.Name, arg.ArgType.Value())
		}

		fmt.Println(")")

		for _, label := range proc.Labels {
			fmt.Printf("  %s:\n", label.Name)
			for i, v := range label.Insts {
				var ArgOne string
				if v.ArgOne != nil {
					ArgOne = v.ArgOne.InstructionArg()
				}

				var ArgTwo string
				if v.ArgTwo != nil {
					ArgTwo = v.ArgTwo.InstructionArg()
				}

				fmt.Printf(" %5d %20s  %10s  %10s\n", i, v.Operation, ArgOne, ArgTwo)
			}
		}

		fmt.Println()
	}

	fmt.Println("main:")
	for _, label := range gen.Main.Labels {
		fmt.Printf("  %s:\n", label.Name)

		for i, v := range label.Insts {
			var ArgOne string
			if v.ArgOne != nil {
				ArgOne = v.ArgOne.InstructionArg()
			}

			var ArgTwo string
			if v.ArgTwo != nil {
				ArgTwo = v.ArgTwo.InstructionArg()
			}

			fmt.Printf(" %5d %20s  %10s  %10s\n", i, v.Operation, ArgOne, ArgTwo)
		}
	}

	//	clang := findCommand("clang-19", "clang")
	//	objectCmd := exec.Command(clang, "-Wall", "-c", c.asmFileName(), "-o", c.objectFileName(), "-fPIE")
	//	objectCmd.Stdout = os.Stdout
	//	objectCmd.Stderr = os.Stderr
	//
	//	err = objectCmd.Run()
	//	if err != nil {
	//		return fmt.Errorf("Error durrng object creation <%w>", err)
	//	}
	//
	//	linkArgs := []string{c.objectFileName(), "-o", c.executableFileName()}
	//	linkCmd := exec.Command(clang, linkArgs...)
	//	linkCmd.Stdout = os.Stdout
	//	linkCmd.Stderr = os.Stderr
	//
	//	err = linkCmd.Run()
	//	if err != nil {
	//		return fmt.Errorf("Error durrng linking <%w>", err)
	//	}

	return nil
}

func (c ASMCompiler) asmFileName() string {
	return fmt.Sprintf("%s.asm", c.req.Target.Output)
}

func (c ASMCompiler) objectFileName() string {
	return fmt.Sprintf("%s.obj", c.req.Target.Output)
}

func (c ASMCompiler) executableFileName() string {
	return fmt.Sprintf("%s.exe", c.req.Target.Output)
}
