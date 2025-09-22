package tests

import (
	"testing"
)

func TestStructs(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
	})

	t.Run("English", func(t *testing.T) {
	})
}

func TestStructWithUnknownType(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
	})

	t.Run("English", func(t *testing.T) {
	})
}

func TestStructPropertyInComplexExpression(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
	})

	t.Run("English", func(t *testing.T) {
	})
}

func TestStructPropertyInConditional(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
	})

	t.Run("English", func(t *testing.T) {
	})
}

func TestStructPropertyAssignment(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
	})

	t.Run("English", func(t *testing.T) {
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

		if err := req.Compile(); err != nil {
			t.Fatalf("Compiler error (%s)", err)
		}

		req.AssertGeneratedLLIR()

		if err := req.RunProgram(); err != nil {
			t.Fatalf("Runtime error (%s)", err)
		}
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:    "./structs/return-expression/source.french.swa",
			ExpectedLLIR: "./structs/return-expression/source.french.ll",
			OutputPath:   "ca2cd604-76a7-40d8-992d-186adedbe12a",
			T:            t,
		}

		defer req.Cleanup()

		if err := req.Compile(); err != nil {
			t.Fatalf("Compiler error (%s)", err)
		}

		req.AssertGeneratedLLIR()

		if err := req.RunProgram(); err != nil {
			t.Fatalf("Runtime error (%s)", err)
		}
	})
}
