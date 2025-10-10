package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type CompileRequest struct {
	Folder       string
	InputPath    string
	OutputPath   string
	ExpectedLLIR string

	// The expected output after the program is compiled
	ExpectedOutput string

	// The expected output when the program is executed on the host machine
	ExpectedExecutionOutput string
	T                       *testing.T
}

func (cr *CompileRequest) Compile() error {
	cr.T.Helper()

	cr.OutputPath = uuid.New().String()
	cmd := exec.Command("./swa", "compile", "-s", cr.InputPath, "-o", cr.OutputPath)

	output, err := cmd.CombinedOutput()

	assert.Equal(cr.T, cr.ExpectedOutput, string(output))

	return err
}

func (cr *CompileRequest) RunProgram() error {
	cr.T.Helper()

	cmd := exec.Command(fmt.Sprintf("./%s.exe", cr.OutputPath))

	output, err := cmd.CombinedOutput()
	if cr.ExpectedExecutionOutput != string(output) {
		cr.T.Fatalf(
			"Execution error want: %s, has: %s",
			cr.ExpectedExecutionOutput,
			string(output),
		)
	}

	return err
}

func (cr *CompileRequest) AssertCompileAndExecute() {
	if err := cr.Compile(); err != nil {
		cr.T.Fatalf("Compiler error (%s)", err)
	}

	if err := cr.RunProgram(); err != nil {
		cr.T.Fatalf("Runtime error (%s)", err)
	}
}

func (cr *CompileRequest) Cleanup() {
	cr.T.Helper()

	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.ll", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.s", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.o", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.exe", cr.OutputPath)))
}
