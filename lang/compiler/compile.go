package compiler

import (
	"os/exec"
)

type BuildTarget struct {
	OperatingSystem string
	Architecture    string
	Output          string
}

const FilePerm = 0600

func Compile(req LLVMCompilerRequest) error {
	cmp := NewLLVMCompiler(req)

	err := cmp.Run()
	if err != nil {
		return err
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
