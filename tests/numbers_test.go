package tests

import (
	"testing"
)

func TestMixedNumbers(t *testing.T) {
	t.Run("Mixed", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numbers/mixed.english.swa",
			ExpectedExecutionOutput: "(x + y = 7.50, y + x = 7.50)(x * y = 12.50, y * x = 12.50)(x / y = 2.00, y / x = 0.50)(x - y = 2.50, y - x = -2.50)",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})
}

func TestSignedNumbers(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {

		t.Run("Signed", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./numbers/signed.english.swa",
				ExpectedExecutionOutput: "x+y: 0, x*y: -4, x-y: -4, x/y: -1, x%y: 0",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("integer arithmetic", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./numbers/integer-arithmetic/source.english.swa",
				ExpectedExecutionOutput: "okokokokokok",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("float-arithmetic", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./numbers/float-arithmetic/source.english.swa",
				ExpectedExecutionOutput: "okokokokokok",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("integer-arrays", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./numbers/integer-arrays/source.english.swa",
				ExpectedExecutionOutput: "okokok",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("float-arrays", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./numbers/float-arrays/source.english.swa",
				ExpectedExecutionOutput: "okokok",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("integer-float-separation", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./numbers/integer-float-separation/source.english.swa",
				ExpectedExecutionOutput: "",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numbers/signed.french.swa",
			ExpectedExecutionOutput: "x+y: 0, x*y: -4, x-y: -4, x/y: -1, x%y: 0",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
