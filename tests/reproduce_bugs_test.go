package tests

import (
	"testing"
)

func TestGlobals(t *testing.T) {
	t.Parallel()

	expectedExecutionOutput := "LANG: Swahili,PI: 3.14, SPEED_OF_LIGHT: 300000, COMPILED: 1"

	t.Run("Wolof", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./globals/1.wolof.swa",
			expectedExecutionOutput)
	})

	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./globals/1.french.swa",
			expectedExecutionOutput)
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./globals/1.soussou.swa",
			expectedExecutionOutput)
	})

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./globals/1.english.swa",
			expectedExecutionOutput)
	})

	t.Run("Errors", func(t *testing.T) {
		t.Run("2", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./globals/2.english.swa",
				"struct initialization should happen inside a function\n")
		})

		t.Run("3", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./globals/3.english.swa",
				"array initialization should happen inside a function\n")
		})
	})
}

func TestBugInvalidArrayAccess(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-array-access/source.english.swa",
			"Property age is not an array\n")
	})

	t.Run("Wolof", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-array-access/source.wolof.swa",
			"Property age is not an array\n")
	})

	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-array-access/source.french.swa",
			"ArrayAccessExpression property age is not an array\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-array-access/source.soussou.swa",
			"Property age mu tableau ra\n")
	})
}

func TestBugInvalidFieldAccess(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-field-access/source.english.swa",
			"variable i is not a struct instance\n")
	})

	t.Run("Wolof", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-field-access/source.wolof.swa",
			"struct named int does not exist in symbol table\n")
	})

	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-field-access/source.french.swa",
			"variable i is not a struct instance\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./bugs/invalid-field-access/source.soussou.swa",
			"variable i mu fokhi ra\n")
	})
}

func TestBugArrayOfStructs(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/array-of-structs/source.english.swa",
			"")
	})

	t.Run("Wolof", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/array-of-structs/source.wolof.swa",
			"")
	})

	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/array-of-structs/source.french.swa",
			"")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/array-of-structs/source.soussou.swa",
			"")
	})
}

func TestBugStructAssignment(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/struct-assignment/source.english.swa",
			"999")
	})

	t.Run("Wolof", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/struct-assignment/source.wolof.swa",
			"999")
	})

	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/struct-assignment/source.french.swa",
			"999")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./bugs/struct-assignment/source.soussou.swa",
			"999")
	})
}
