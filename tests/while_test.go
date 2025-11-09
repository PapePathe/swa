package tests

import (
	"testing"
)

func TestWhileStatement(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./while/simple.source.swa",
				T:                       t,
				ExpectedExecutionOutput: "0 1 2 3 4 5 6 7 8 9 ",
			}

			req.AssertCompileAndExecute()
		})
	})
}
