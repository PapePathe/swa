package tests

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MultiDialectTest struct {
	name                    string
	inputPath               string
	expectedOutput          string
	expectedExecutionOutput string
}

type CompileRequest struct {
	InputPath  string
	OutputPath string
	// The expected output after the program is compiled or parsed
	ExpectedOutput string
	// The expected output when the program is executed on the host machine
	ExpectedExecutionOutput string
	T                       *testing.T
}

func NewSuccessfulCompileRequest(t *testing.T, input string, output string) {
	req := CompileRequest{InputPath: input, ExpectedExecutionOutput: output, T: t}
	req.AssertCompileAndExecute()
}

func (cr *CompileRequest) Parse(format string) error {
	cr.T.Helper()

	cmd := exec.Command("./swa", "parse", "-s", cr.InputPath, "-o", format)
	output, err := cmd.CombinedOutput()

	if cr.ExpectedOutput != string(output) {
		cr.T.Fatalf(
			"Parser error want: %s, has: %s \n Source file %s",
			fmt.Sprintf("%q", cr.ExpectedOutput),
			fmt.Sprintf("%q", string(output)),
			cr.InputPath,
		)
	}

	return err
}

func (cr *CompileRequest) Compile() error {
	cr.T.Helper()

	cr.OutputPath = uuid.New().String()

	cmd := exec.Command("./swa", "compile", "-s", cr.InputPath, "-o", cr.OutputPath)
	output, err := cmd.CombinedOutput()

	if cr.ExpectedOutput != string(output) {
		cr.T.Fatalf(
			"Compilation error want: %s, has: %s \n Source file %s",
			fmt.Sprintf("%q", cr.ExpectedOutput),
			fmt.Sprintf("%q", string(output)),
			cr.InputPath,
		)
	}

	return err
}

func (cr *CompileRequest) RunProgram() error {
	cr.T.Helper()

	cmd := exec.Command(fmt.Sprintf("./%s.exe", cr.OutputPath))

	output, err := cmd.CombinedOutput()
	if cr.ExpectedExecutionOutput != string(output) {
		cr.T.Fatalf(
			"Execution error want: %s, has: %s \n Source file %s",
			fmt.Sprintf("%q", cr.ExpectedExecutionOutput),
			fmt.Sprintf("%q", string(output)),
			cr.InputPath,
		)
	}

	return err
}

func (cr *CompileRequest) AssertCompileAndExecute() {
	if err := cr.Compile(); err != nil {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "Compiler error (%s)\n Source file %s", err, cr.InputPath)

		data, err := os.ReadFile(fmt.Sprintf("%s.ll", cr.OutputPath))
		if err == nil {
			sb.WriteString("\n")
			sb.WriteString(string(data))
		}

		cr.T.Fatal(sb.String())
	}

	defer cr.Cleanup()

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
