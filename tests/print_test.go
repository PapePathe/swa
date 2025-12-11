package tests

import (
	"testing"
)

func TestPrint(t *testing.T) {
	t.Run("static string", func(t *testing.T) {
		t.Parallel()

		NewSuccessfulCompileRequest(
			t,
			"./print/binexpr.swa",
			"2",
		)
	})

	t.Run("table display", func(t *testing.T) {
		t.Parallel()

		NewSuccessfulCompileRequest(
			t,
			"./print/table.swa",
			`   id|                dept| salary
   20|     human resources|1000000
   21|          accounting|1000000
   22|             finance|1000000
   23|         engineering|1000000
`,
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
