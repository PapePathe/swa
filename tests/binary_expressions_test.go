package tests

import (
	"testing"
)

func TestBinaryExpressionsModulo(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		t.Run("Modulo with symbol expressions", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/variable-declaration-2.english.swa",
				T:                       t,
				ExpectedExecutionOutput: "v1: 0, v2: 1",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Modulo in conditionals", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/conditionals.english.swa",
				T:                       t,
				ExpectedExecutionOutput: "(4 modulo 3 is equal to 1)(4 modulo 2 is equal to 0)",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Modulo with integer expressions", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/variable-declaration.english.swa",
				T:                       t,
				ExpectedExecutionOutput: "x: 1, y: 0",
			}

			req.AssertCompileAndExecute()
		})
	})
}
