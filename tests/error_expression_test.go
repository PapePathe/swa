package tests

import "testing"

func TestErrorExpression(t *testing.T) {
	t.Run("With error", func(t *testing.T) {
		t.Run("english", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./error-expression/with-error.english.swa",
				ExpectedExecutionOutput: "Division by zero error dividend is zero",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})
	})
}
