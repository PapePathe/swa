package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayOfStructsWithUndefinedPropertyAccess(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./arrays/structs/undefined-property-access.french.swa",
			ExpectedOutput: "property (Name) not found in struct Engineer\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructsWithUndefinedPropertyInitialization(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./arrays/structs/undefined-property-initialization.french.swa",
			ExpectedOutput: "Property with name (Name) does not exist on struct Engineer\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructsWithUndefined(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./arrays/structs/undefined-struct.english.swa",
			ExpectedOutput: "struct named Engineer does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructs(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./arrays/structs/source.french.swa",
			ExpectedExecutionOutput: "(nom: Pathe, age: 40, taille: 1.80, technos: Ruby, Rust, Go) (nom: Lucien, age: 24, taille: 1.81, technos: Typescript, HTML, Css) (nom: Manel, age: 25, taille: 1.82, technos: Typescript, Ruby) (nom: Bintou, age: 28, technos: Javascript, Css, HTML)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./arrays/structs/source.english.swa",
			ExpectedExecutionOutput: "(nom: Pathe, age: 40, taille: 1.80, stack: Ruby, Rust, Go) (nom: Lucien, age: 24, taille: 1.81, stack: Typescript, HTML, Css) (nom: Manel, age: 25, taille: 1.82, stack: Typescript, Ruby) (nom: Bintou, age: 28, taille: 1.83, stack: Javascript, Css, HTML) ",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArrayOfStrings(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./arrays/strings/source.french.swa",
			ExpectedExecutionOutput: "valeurs dans le tableau: (abc),(efg),(ijk)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./arrays/strings/source.english.swa",
			ExpectedExecutionOutput: "values in the array: (abc),(efg),(ijk)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArraysInPrintStatement(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./arrays/print/source.french.swa",
			ExpectedExecutionOutput: "Les valeurs dans le tableau sont: 1 2 3 4 5",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./arrays/print/source.english.swa",
			ExpectedExecutionOutput: "Array values are: 1 2 3 4 5",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestArrays(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath: "./arrays/ok/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath: "./arrays/ok/source.english.swa",
			T:         t,
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
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./arrays/out-of-bounds/source.english.swa",
			ExpectedOutput: "Element at index (%!s(int=5)) does not exist in array (array)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayAccessWithPrefixExpression(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.french.swa",
			ExpectedOutput: "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (-4)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.english.swa",
			ExpectedOutput: "Only numbers are supported as array index, current: (-4)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestArrayOfStructsAsFuncParam(t *testing.T) {
	// TODO find and fix the bug
	t.Skip("THERE IS A KNOWN BUG HERE THAT MUST BE FIXED")
	t.Run("Search Occurrence", func(t *testing.T) {
		t.Parallel()
		NewSuccessfulCompileRequest(
			t,
			"./arrays/structs/as-function-params.swa",
			"X=1 occurs 2 times, X=5 occurs 2 times, X=99 occurs 0 times",
		)
	})

	t.Run("High Earners", func(t *testing.T) {
		t.Parallel()
		NewSuccessfulCompileRequest(
			t,
			"./arrays/structs/count-high-earners.swa",
			"Employees with salary above 50000: 2\n",
		)
	})
}

func TestSortingAlgorithms(t *testing.T) {
	t.Run("QSort", func(t *testing.T) {
		t.Parallel()
		NewSuccessfulCompileRequest(
			t,
			"./arrays/sorting/qsort.swa",
			"Sorted List: 1 2 3 4 5 ",
		)
	})

	t.Run("Selection", func(t *testing.T) {
		t.Parallel()

		NewSuccessfulCompileRequest(
			t,
			"./arrays/sorting/selection.swa",
			"Sorted Array: 5 8 10 11 22 25 34 64 77 90 \n",
		)
	})

	t.Run("Bubble", func(t *testing.T) {
		t.Parallel()

		NewSuccessfulCompileRequest(
			t,
			"./arrays/sorting/bubble.swa",
			"Sorted Array: 4 8 11 12 22 25 34 64 77 90 \n",
		)
	})
}

func TestArrayAccessWithSymbolExpression(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.french.swa",
			ExpectedOutput: "variable named x does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.english.swa",
			ExpectedOutput: "variable named x does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}
