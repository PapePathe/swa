package tests

import (
	"os"
	"testing"
)

func TestConvertAndCompileFrench(t *testing.T) {
	input := `dialect:english;

start() int {
	print("Bonjour le monde");
	return 0;
}
`
	inputFile := "test_french_input.swa"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)

	req := CompileRequest{
		InputPath:               inputFile,
		ExpectedExecutionOutput: "Bonjour le monde",
		T:                       t,
	}

	req.AssertConvertAndCompile("french")
}

func TestComprehensiveDialectConversion(t *testing.T) {
	// Source file created previously
	inputFile := "comprehensive.swa"

	// Ensure the source file exists (it should have been created by the tool, but for the test runner...)
	// Since the tool created it in real FS, we can use it. But for reproducibility in test suite,
	// good practice might be to write it here or expect it in testdata.
	// For now, assuming it's in the current directory (tests/).

	dialects := []string{
		"french",
		"wolof",
		"malinke",
		"soussou",
		// "igbo",
	}

	expectedOutput := "Flag is trueSum is greater than 25012"

	for _, dialect := range dialects {
		t.Run(dialect, func(t *testing.T) {
			req := CompileRequest{
				InputPath:               inputFile,
				ExpectedExecutionOutput: expectedOutput,
				T:                       t,
			}
			req.AssertConvertAndCompile(dialect)
		})
	}
}
