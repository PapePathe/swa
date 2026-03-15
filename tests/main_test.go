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

func NewFailedCompileRequest(t *testing.T, input string, output string) {
	t.Parallel()

	req := CompileRequest{InputPath: input, ExpectedOutput: output, T: t}
	assert.Error(t, req.Compile())
}

func SuccessFulMultidialectCompileRequest(t *testing.T, tests []MultiDialectTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req := CompileRequest{
				InputPath:               test.inputPath,
				T:                       t,
				ExpectedExecutionOutput: test.expectedExecutionOutput,
			}

			req.AssertCompileAndExecute()
		})
	}
}

type OptName int

const (
	OptCompileNative OptName = iota
)

func NewSuccessfulCompileRequest(t *testing.T, input string, output string, opts ...OptName) {
	t.Parallel()

	req := CompileRequest{InputPath: input, ExpectedExecutionOutput: output, T: t}
	req.AssertCompileAndExecute()

	//	for _, opt := range opts {
	//		switch opt {
	//		case OptCompileNative:
	//			req.T.Run("CompileX", func(t *testing.T) {
	//				t.Parallel()
	//
	//				req.AssertCompileAndExecuteX()
	//			})
	//		default:
	//			panic("Unknown native compile opt name")
	//		}
	//	}
}

func NewSuccessfulXCompileRequest(req CompileRequest) {
	req.T.Parallel()

	req.AssertCompileAndExecuteX()
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

func (cr *CompileRequest) CompileX() error {
	cr.T.Helper()

	cr.OutputPath = uuid.New().String()

	cmd := exec.Command("./swa", "compile", "-g", "swa", "-s", cr.InputPath, "-o", cr.OutputPath, "-g", "swa")
	output, err := cmd.CombinedOutput()

	fmt.Println(string(output))

	if cr.ExpectedOutput != string(output) {
		cr.T.Fatalf(
			"Compilation error want: %s, has: %s \n Source file %s",
			fmt.Sprintf("%q", cr.ExpectedOutput),
			fmt.Sprintf("%q", string(output)),
			cr.InputPath,
		)
	}

	if err != nil {
		return err
	}

	return nil
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

	if err != nil {
		return err
	}

	//	cmd = exec.Command("./swa", "parse", "-s", cr.InputPath, "-o", "json")
	//	_, err = cmd.CombinedOutput()
	//
	//	assert.NoError(cr.T, err)
	//
	//	cmd = exec.Command("./swa", "parse", "-s", cr.InputPath, "-o", "tree")
	//	_, err = cmd.CombinedOutput()
	//
	//	assert.NoError(cr.T, err)

	return nil
}

func (cr *CompileRequest) RunProgramX() error {
	cr.T.Helper()

	cmd := exec.Command(fmt.Sprintf("./%s.native.exe", cr.OutputPath))

	output, err := cmd.CombinedOutput()
	if cr.ExpectedExecutionOutput != string(output) {
		cr.T.Fatalf(
			"Execution error want: %s, has: %s \n Source file %s",
			fmt.Sprintf("%q", cr.ExpectedExecutionOutput),
			fmt.Sprintf("%q", string(output)),
			cr.InputPath,
		)
	}

	//	defer cr.CleanupX()

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

	//	defer cr.Cleanup()

	return err
}

func (cr *CompileRequest) AssertCompileAndExecuteX() {
	if err := cr.CompileX(); err != nil {
		cr.T.Fatal(err.Error())
	}

	if err := cr.RunProgramX(); err != nil {
		cr.T.Fatalf("Runtime error (%s)", err)
	}
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

	if err := cr.RunProgram(); err != nil {
		cr.T.Fatalf("Runtime error (%s)", err)
	}
}

func (cr *CompileRequest) AssertConvertAndCompile(targetDialect string) {
	cr.T.Helper()

	// 1. Convert
	convertedPath := strings.Replace(cr.InputPath, ".swa", fmt.Sprintf(".%s.swa", targetDialect), 1)
	cmd := exec.Command("./swa", "convert", "--source", cr.InputPath, "--target-dialect", targetDialect, "--out", convertedPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		cr.T.Fatalf("Conversion failed: %v\n%s", err, output)
	}

	// 2. Update InputPath to point to converted file
	cr.InputPath = convertedPath

	// 3. Compile and Execute
	if err := cr.Compile(); err != nil {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "Compiler error (%s)\n Source file %s", err, cr.InputPath)

		data, err := os.ReadFile(convertedPath)
		if err == nil {
			sb.WriteString("\nConverted Content:\n")
			sb.WriteString(string(data))
		}

		cr.T.Fatal(sb.String())
	}

	//	defer cr.Cleanup()

	if err := cr.RunProgram(); err != nil {
		cr.T.Fatalf("Runtime error (%s)", err)
	}

	// Cleanup converted file
	// os.Remove(convertedPath)
}

func (cr *CompileRequest) Cleanup() {
	cr.T.Helper()

	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.ll", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.s", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.o", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.exe", cr.OutputPath)))
}

func (cr *CompileRequest) CleanupX() {
	cr.T.Helper()

	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.native.s", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.native.o", cr.OutputPath)))
	assert.NoError(cr.T, os.Remove(fmt.Sprintf("%s.native.exe", cr.OutputPath)))
}
