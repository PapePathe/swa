package tests

import (
	"testing"
)

func TestBugFixes(t *testing.T) {
	t.Parallel()

	t.Run("116-keyword-regex-patterns-lack-word-boundaries", func(t *testing.T) {
		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath: "./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.french.swa",
				T:         t,
			}

			req.AssertCompileAndExecute()
		})
	})
}
