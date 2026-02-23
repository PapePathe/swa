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
	inputFile := "comprehensive.swa"
	dialects := []string{
		"french",
		"wolof",
		"malinke",
		"soussou",
	}

	expectedOutput := "Flag is true\nSum is greater than 25\n012"

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
