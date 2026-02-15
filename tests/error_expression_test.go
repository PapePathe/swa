package tests

import "testing"

func TestErrorExpression(t *testing.T) {
	t.Run("With error", func(t *testing.T) {
		t.Run("english", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./error-expression/with-error.english.swa",
				"Division by zero error dividend is zero",
			)
		})
	})
}
