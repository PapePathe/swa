package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBugFixes(t *testing.T) {
	t.Parallel()

	t.Run("99-missing-type-check-in-variable-declaration", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.1.swa",
				ExpectedOutput: "expected DataTypeString got DoubleTypeKind\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.2.swa",
				ExpectedOutput: "expected DataTypeString got IntegerTypeKind\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("3", func(t *testing.T) {
			//			TODO: fix validation when var type is string
			_ = CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.3.swa",
				ExpectedOutput: "",
				T:              t,
			}

			//			assert.Error(t, req.Compile())
		})

		t.Run("4", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.4.swa",
				ExpectedOutput: "expected DataTypeNumber got pointer of %!v(PANIC=String method: unreachable)\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("5", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.5.swa",
				ExpectedOutput: "expected DataTypeNumber got DoubleTypeKind\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("6", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/99-missing-type-check-in-variable-declaration.english.6.swa",
				ExpectedOutput: "expected DataTypeNumber got ArrayTypeKind\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

	})

	t.Run("103-missing-argument-type-check-in-function-calls", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./bug-fixes/103-missing-argument-type-check-in-function-calls.english.swa",
				ExpectedOutput: "expected argument of type IntegerType(32 bits) expected but got ArrayType(IntegerType(8 bits)[3])\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})

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

	t.Run("91-array-of-structs-indexing-does-not-support-variables", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./bug-fixes/91-array-of-structs-indexing-does-not-support-variables.english.swa",
				ExpectedExecutionOutput: "Name: Alice",
				T:                       t,
			}

			req.AssertCompileAndExecute()
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
