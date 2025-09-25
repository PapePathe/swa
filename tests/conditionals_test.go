package tests

import (
	"testing"
)

func TestGreaterThanEquals(t *testing.T) {
	t.Parallel()
	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals/source.english.swa",
			ExpectedLLIR:            "./conditionals/greater-than-equals/source.english.ll",
			OutputPath:              "64e945b0-5e88-49ce-848c-6dfb4af57412",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals/source.french.swa",
			ExpectedLLIR:            "./conditionals/greater-than-equals/source.french.ll",
			OutputPath:              "5d3bad62-1a59-42db-9308-505dcab59f02",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestGreaterThanEqualsWithPointerAndInt(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals-pointer-and-int/source.english.swa",
			ExpectedLLIR:            "./conditionals/greater-than-equals-pointer-and-int/source.english.ll",
			OutputPath:              "98f0c855-3a27-488f-980c-e6a375d5a627",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/greater-than-equals-pointer-and-int/source.french.swa",
			ExpectedLLIR:            "./conditionals/greater-than-equals-pointer-and-int/source.french.ll",
			OutputPath:              "5e51c0c9-e8bc-494f-aa01-2afc0a238b91",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestEquals(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.english.swa",
			OutputPath:              "530955fd-7a49-4264-a832-53801a81b019",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./conditionals/equals/source.french.swa",
			OutputPath:              "10a66ac2-024f-4717-be0d-6ccb5379d490",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}
