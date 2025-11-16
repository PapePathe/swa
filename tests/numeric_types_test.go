package tests

import (
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numeric-types/integer-arithmetic/source.english.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestFloatArithmetic(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numeric-types/float-arithmetic/source.english.swa",
			ExpectedExecutionOutput: "okokokokokok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestTypeDeclarations(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numeric-types/type-declarations/source.english.swa",
			ExpectedExecutionOutput: "okok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestIntegerArrays(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numeric-types/integer-arrays/source.english.swa",
			ExpectedExecutionOutput: "okokok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestFloatArrays(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numeric-types/float-arrays/source.english.swa",
			ExpectedExecutionOutput: "okokok",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestIntegerFloatSeparation(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numeric-types/integer-float-separation/source.english.swa",
			ExpectedExecutionOutput: "",
			T:                       t,
		}

		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}
