package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBugFixes(t *testing.T) {
	t.Parallel()

	t.Run("114-variable-redeclaration-allowed-in-same-scope", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.english.swa",
				ExpectedOutput: "variable x is aleady defined\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/114-variable-redeclaration-allowed-in-same-scope.french.swa",
				ExpectedOutput: "variable x is aleady defined\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("116-keyword-regex-patterns-lack-word-boundaries", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath: "./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.english.swa",
				T:         t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath: "./bug-fixes/116-keyword-regex-patterns-lack-word-boundaries.french.swa",
				T:         t,
			}

			req.AssertCompileAndExecute()
		})
	})
}
