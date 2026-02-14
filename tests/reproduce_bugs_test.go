package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobals(t *testing.T) {
	t.Parallel()

	t.Run("Wolof", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./globals/1.wolof.swa",
			ExpectedExecutionOutput: "LANG: Swahili,PI: 3.14, SPEED_OF_LIGHT: 300000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./globals/1.french.swa",
			ExpectedExecutionOutput: "LANG: Swahili,PI: 3.14, SPEED_OF_LIGHT: 300000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./globals/1.soussou.swa",
			ExpectedExecutionOutput: "LANG: Swahili,PI: 3.14, SPEED_OF_LIGHT: 300000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {

		t.Run("1", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./globals/1.english.swa",
				ExpectedExecutionOutput: "LANG: Swahili,PI: 3.14, SPEED_OF_LIGHT: 300000",
				T:                       t,
			}

			req.AssertCompileAndExecute()
		})

		t.Run("2", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./globals/2.english.swa",
				ExpectedOutput: "struct initialization should happen inside a function\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})

		t.Run("3", func(t *testing.T) {
			req := CompileRequest{
				InputPath:      "./globals/3.english.swa",
				ExpectedOutput: "array initialization should happen inside a function\n",
				T:              t,
			}

			assert.Error(t, req.Compile())
		})
	})
}

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

	t.Run("Wolof", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./bugs/invalid-array-access/source.wolof.swa",
			ExpectedOutput: "Property age is not an array\n", // Expected correct error
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./bugs/invalid-array-access/source.french.swa",
			ExpectedOutput: "ArrayAccessExpression property age is not an array\n", // Expected correct error
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./bugs/invalid-array-access/source.soussou.swa",
			ExpectedOutput: "Property age mu tableau ra\n", // Expected correct error
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

	t.Run("Wolof", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./bugs/invalid-field-access/source.wolof.swa",
			ExpectedOutput: "struct named int does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./bugs/invalid-field-access/source.french.swa",
			ExpectedOutput: "variable i is not a struct instance\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./bugs/invalid-field-access/source.soussou.swa",
			ExpectedOutput: "variable i mu fokhi ra\n",
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

	t.Run("Wolof", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath: "./bugs/array-of-structs/source.wolof.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath: "./bugs/array-of-structs/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath: "./bugs/array-of-structs/source.soussou.swa",
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

	t.Run("Wolof", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./bugs/struct-assignment/source.wolof.swa",
			ExpectedExecutionOutput: "999",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./bugs/struct-assignment/source.french.swa",
			ExpectedExecutionOutput: "999",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./bugs/struct-assignment/source.soussou.swa",
			ExpectedExecutionOutput: "999",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
