package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructWithUnknownType(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.french.swa",
			ExpectedOutput: "struct proprerty type ({unknown_type}) not supported\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.english.swa",
			ExpectedOutput: "struct proprerty type ({unknown_type}) not supported\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestStructPropertyAssignment(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/assignment/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/assignment/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructAll(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.english.swa",
			OutputPath:              "15a6540e-5df1-4fbd-b3a1-a2edb12d117b",
			ExpectedExecutionOutput: "Name: (Pathe), TechStack: (Ruby, Rust, Go), Age: (40)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.french.swa",
			OutputPath:              "83b453cc-bd0b-432d-96d5-f0192dc2c2d3",
			ExpectedExecutionOutput: "Nom: (Pathe), Stack Technique: (Ruby, Rust, Go), Age: (40)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructPropertyInReturnExpression(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/return-expression/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/return-expression/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}
