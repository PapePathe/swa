package tests

import (
	"testing"
)

func TestConditionals(t *testing.T) {
	t.Parallel()

	// t.Run("English", func(t *testing.T) {
	// 	req := CompileRequest{
	// 		InputPath:    "./structs/return-expression/source.english.swa",
	// 		ExpectedLLIR: "./structs/return-expression/source.english.ll",
	// 		OutputPath:   "uuid",
	// 		T:            t,
	// 	}

	// 	defer req.Cleanup()

	// 	if err := req.Compile(); err != nil {
	// 		t.Fatalf("Compiler error (%s)", err)
	// 	}

	// 	req.AssertGeneratedLLIR()

	// 	if err := req.RunProgram(); err != nil {
	// 		t.Fatalf("Runtime error (%s)", err)
	// 	}
	// })

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals/source.french.swa",
			ExpectedLLIR:            "./conditionals/greater-than-equals/source.french.ll",
			OutputPath:              "5d3bad62-1a59-42db-9308-505dcab59f02",
			ExpectedExecutionOutput: "okok",
			T:                       t,
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

func TestGreaterThanEqualsWithPointerAndInt(t *testing.T) {
	t.Parallel()

	// t.Run("English", func(t *testing.T) {
	// 	req := CompileRequest{
	// 		InputPath:    "./structs/return-expression/source.english.swa",
	// 		ExpectedLLIR: "./structs/return-expression/source.english.ll",
	// 		OutputPath:   "uuid",
	// 		T:            t,
	// 	}

	// 	defer req.Cleanup()

	// 	if err := req.Compile(); err != nil {
	// 		t.Fatalf("Compiler error (%s)", err)
	// 	}

	// 	req.AssertGeneratedLLIR()

	// 	if err := req.RunProgram(); err != nil {
	// 		t.Fatalf("Runtime error (%s)", err)
	// 	}
	// })

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals-pointer-and-int/source.french.swa",
			ExpectedLLIR:            "./conditionals/greater-than-equals-pointer-and-int/source.french.ll",
			OutputPath:              "5e51c0c9-e8bc-494f-aa01-2afc0a238b91",
			ExpectedExecutionOutput: "okok",
			T:                       t,
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
