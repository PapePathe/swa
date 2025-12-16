package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBugInvalidArrayAccess(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./bugs/invalid-array-access/source.english.swa",
			ExpectedOutput: "Property age is not an array\n", // Expected correct error
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestBugInvalidFieldAccess(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./bugs/invalid-field-access/source.english.swa",
			ExpectedOutput: "variable i is not a struct instance\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestBugArrayOfStructs(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./bugs/array-of-structs/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestBugStructAssignment(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./bugs/struct-assignment/source.english.swa",
			ExpectedExecutionOutput: "999",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
