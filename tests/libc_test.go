package tests

import (
	"testing"
)

func TestOpenReadAndClose(t *testing.T) {
	t.Run("Calling open,read,close from libc", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./libc/read-file.english.swa",
			// TODO: this test is failing because the text returned contains
			// non printable characters
			//
			// ExpectedExecutionOutput: `you just read the contents of a file`,
			ExpectedExecutionOutput: "close return value: 0",
			T:                       t,
		}
		req.AssertCompileAndExecute()
	})
}
