package tests

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		t.Run("Substract", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./functions/substract.source.french.swa",
				ExpectedOutput:          "",
				ExpectedExecutionOutput: "10 - 5 = 5",
				T:                       t,
			}
			defer req.Cleanup()

			req.AssertCompileAndExecute()
		})

		t.Run("Add", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./functions/add.source.french.swa",
				ExpectedOutput:          "",
				ExpectedExecutionOutput: "10 + 5 = 15",
				T:                       t,
			}
			defer req.Cleanup()

			req.AssertCompileAndExecute()
		})

		t.Run("Divide", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./functions/divide.source.french.swa",
				ExpectedOutput:          "",
				ExpectedExecutionOutput: "10 / 5 = 2",
				T:                       t,
			}
			defer req.Cleanup()

			req.AssertCompileAndExecute()
		})

		t.Run("Multiply", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./functions/multiply.source.french.swa",
				ExpectedOutput:          "",
				ExpectedExecutionOutput: "10 * 5 = 50",
				T:                       t,
			}
			defer req.Cleanup()

			req.AssertCompileAndExecute()
		})
	})
}
