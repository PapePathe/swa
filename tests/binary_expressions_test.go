package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryExpressionsModulo(t *testing.T) {

	t.Run("Strings in binary-expressions", func(t *testing.T) {
		t.Parallel()

		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./binary-expressions/strings/1.english.swa",
				T:              t,
				ExpectedOutput: "Strings are not supported in symbol right of binary expression\n",
			}

			assert.Error(t, req.Compile())
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./binary-expressions/strings/2.english.swa",
				T:              t,
				ExpectedOutput: "Strings are not supported in symbol left of binary expression\n",
			}

			assert.Error(t, req.Compile())
		})

		t.Run("3", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./binary-expressions/strings/3.english.swa",
				T:              t,
				ExpectedOutput: "Strings are not supported in left of binary expression\n",
			}

			assert.Error(t, req.Compile())
		})

		t.Run("4", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./binary-expressions/strings/4.english.swa",
				T:              t,
				ExpectedOutput: "Strings are not supported in right of binary expression\n",
			}

			assert.Error(t, req.Compile())
		})

	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()

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
