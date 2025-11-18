package tests

import (
	"testing"
)

func TestPrint(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		t.Run("static string", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./print/static-string.french.swa",
				ExpectedExecutionOutput: "contenu de la variable: french",
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
