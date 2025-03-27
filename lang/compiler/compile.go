package compiler

import (
	"log"
	"os"
	"os/exec"
	"swahili/lang/ast"

	"github.com/llir/llvm/ir"
)

type BuildTarget struct {
	OperatingSystem string
	Architecture    string
}

func Compile(tree ast.BlockStatement, target BuildTarget) {
	m := ir.NewModule()

	for _, stmt := range tree.Body {
		err := stmt.Compile(m, nil)
		if err != nil {
			panic(err)
		}
	}

	err := os.WriteFile("./tmp/start.ll", []byte(m.String()), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("llc", "./tmp/start.ll", "-o", "./tmp/start.s")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	objectCmd := exec.Command("clang", "-c", "./tmp/start.s", "-o", "./tmp/start.o")
	if err := objectCmd.Run(); err != nil {
		log.Fatal(err)
	}

	linkCmd := exec.Command("clang", "./tmp/start.o", "-o", "./tmp/start.exe")
	if err := linkCmd.Run(); err != nil {
		log.Fatal(err)
	}
}
