package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBugFixes(t *testing.T) {
	t.Parallel()

	t.Run("102-missing-arity-check-in-function-calls", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/102-missing-arity-check-in-function-calls.english.swa",
				ExpectedOutput: "function add expect 2 arguments but was given 1\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

	t.Run("90-function-calls-with-variables-pass-pointers-instead-of-values", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.english.swa",
				ExpectedExecutionOutput: "z: 30",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/90-function-calls-with-variables-pass-pointers-instead-of-values.french.swa",
				ExpectedExecutionOutput: "z: 30",
				T:                       t,
			}
			req.AssertCompileAndExecute()
		})
	})

	t.Run("94-block-statements-do-not-create-new-scopes-shadowing-fails", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.english.swa",
				ExpectedExecutionOutput: "Inner x: 20, Outer x: 10",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("French", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/94-block-statements-do-not-create-new-scopes-shadowing-fails.french.swa",
				ExpectedExecutionOutput: "Inner x: 20, Outer x: 10",
				T:                       t,
			}
			req.AssertCompileAndExecute()
		})
	})

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
