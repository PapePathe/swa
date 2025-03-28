package compiler

import (
	"os"
	"os/exec"
	"swahili/lang/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type BuildTarget struct {
	OperatingSystem string
	Architecture    string
}

func Compile(tree ast.BlockStatement, target BuildTarget) {
	m := ir.NewModule()
	m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))

	err := tree.Compile(ast.NewContext(nil, m))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./tmp/start.ll", []byte(m.String()), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("llc", "./tmp/start.ll", "-o", "./tmp/start.s")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	objectCmd := exec.Command("clang", "-c", "./tmp/start.s", "-o", "./tmp/start.o")
	if err := objectCmd.Run(); err != nil {
		panic(err)
	}

	linkCmd := exec.Command("clang", "./tmp/start.o", "-o", "./tmp/start.exe")
	if err := linkCmd.Run(); err != nil {
		panic(err)
	}
}
