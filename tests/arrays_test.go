package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayOfStructsWithUndefinedPropertyAccess(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/structs/undefined-property-access.french.swa",
			ExpectedOutput: "ArrayOfStructsAccessExpression: property Name not found\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructsWithUndefinedPropertyInitialization(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/structs/undefined-property-initialization.french.swa",
			ExpectedOutput: "StructInitializationExpression: property Name not found\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructsWithUndefined(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/structs/undefined-struct.english.swa",
			ExpectedOutput: "Type ({{Engineer}}) is not a valid struct\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructs(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/structs/source.french.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "(nom: Pathe, age: 40, stack: Ruby, Rust, Go) (nom: Lucien, age: 24, stack: Typescript, HTML, Css) (nom: Manel, age: 25, stack: Typescript, Ruby) (nom: Bintou, age: 28, stack: Javascript, Css, HTML)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/structs/source.english.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "(nom: Pathe, age: 40, stack: Ruby, Rust, Go) (nom: Lucien, age: 24, stack: Typescript, HTML, Css) (nom: Manel, age: 25, stack: Typescript, Ruby) (nom: Bintou, age: 28, stack: Javascript, Css, HTML)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArrayOfStrings(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/strings/source.french.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "valeurs dans le tableau: (abc),(efg),(ijk)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/strings/source.english.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "values in the array: (abc),(efg),(ijk)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArraysInPrintStatement(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/print/source.french.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "Les valeurs dans le tableau sont: 1 2 3 4 5",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/print/source.english.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "Array values are: 1 2 3 4 5",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArrays(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/ok/source.french.swa",
			ExpectedOutput: "",
			T:              t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/ok/source.english.swa",
			ExpectedOutput: "",
			T:              t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArrayIndexOutOfBounds(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/out-of-bounds/source.french.swa",
			ExpectedOutput: "L'element a la position (%!s(int=5)) depasse les limites du tableau (tableau)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/out-of-bounds/source.english.swa",
			ExpectedOutput: "Element at index (%!s(int=5)) does not exist in array (array)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayAccessWithPrefixExpression(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.french.swa",
			ExpectedOutput: "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (-4)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.english.swa",
			ExpectedOutput: "Only positive numbers are supported as array index, current: (-4)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayAccessWithSymbolExpression(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.french.swa",
			ExpectedOutput: "Variable x does not exist\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.english.swa",
			ExpectedOutput: "Variable x does not exist\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}
