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
			ExpectedLLIR:   "./structs/unknown-property-type/source.french.ll",
			OutputPath:     "5e9377cb-904d-484b-9b62-28e531340079",
			ExpectedOutput: "struct proprerty type ({unknown_type}) not supported\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.english.swa",
			ExpectedLLIR:   "./structs/unknown-property-type/source.english.ll",
			OutputPath:     "5e9377cb-904d-484b-9b62-28e531340079",
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
			InputPath:    "./structs/assignment/source.french.swa",
			ExpectedLLIR: "./structs/assignment/source.french.ll",
			OutputPath:   "27a8a248-e3d8-4fd0-b99a-15c9160d30f5",
			T:            t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:    "./structs/assignment/source.english.swa",
			ExpectedLLIR: "./structs/assignment/source.english.ll",
			OutputPath:   "27a8a248-e3d8-4fd0-b99a-15c9160d30f5",
			T:            t,
		}

		defer req.Cleanup()

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

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.french.swa",
			OutputPath:              "83b453cc-bd0b-432d-96d5-f0192dc2c2d3",
			ExpectedExecutionOutput: "Nom: (Pathe), Stack Technique: (Ruby, Rust, Go), Age: (40)",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestStructPropertyInReturnExpression(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:    "./structs/return-expression/source.english.swa",
			ExpectedLLIR: "./structs/return-expression/source.english.ll",
			OutputPath:   "9d87488a-e878-4943-aaf0-7e8a4b09ea1f",
			T:            t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:    "./structs/return-expression/source.french.swa",
			ExpectedLLIR: "./structs/return-expression/source.french.ll",
			OutputPath:   "ca2cd604-76a7-40d8-992d-186adedbe12a",
			T:            t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}
