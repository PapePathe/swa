package tests

import (
	"testing"
)

func TestWhileStatement(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		t.Run("Iterate from ten to zero", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./while/from-ten-to-zero.source.swa",
				T:                       t,
				ExpectedExecutionOutput: "10 9 8 7 6 5 4 3 2 1 ",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Iterate from zero to ten", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./while/from-zero-to-ten.source.swa",
				T:                       t,
				ExpectedExecutionOutput: "0 1 2 3 4 5 6 7 8 9 ",
			}

			req.AssertCompileAndExecute()
		})
	})
}
