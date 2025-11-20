package tests

import (
	"testing"
)

func TestSignedNumbers(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numbers/signed.english.swa",
			ExpectedExecutionOutput: "x+y: 0, x*y: -4, x-y: -4, x/y: -1, x%y: 0",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./numbers/signed.french.swa",
			ExpectedExecutionOutput: "x+y: 0, x*y: -4, x-y: -4, x/y: -1, x%y: 0",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
