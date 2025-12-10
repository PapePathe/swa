package tests

import (
	"testing"
)

func TestPrint(t *testing.T) {
	t.Parallel()
	t.Run("static string", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./print/variable.string.french.swa",
			"contenu de la variable: french",
		)
	})

	t.Run("French", func(t *testing.T) {
		t.Run("static string", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./print/variable.string.french.swa",
				ExpectedExecutionOutput: "contenu de la variable: french",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("float expression", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./print/expression.float.french.swa",
				ExpectedExecutionOutput: "a: 10.005000, b: -10.005000",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("float variable", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./print/variable.float.french.swa",
				ExpectedExecutionOutput: "a: 10.005000, b: -10.005000",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Number variable", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./print/variable.number.french.swa",
				ExpectedExecutionOutput: "contenu de la variable: 10",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})
	})
}
