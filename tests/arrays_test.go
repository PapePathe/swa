package tests

import (
	"testing"
)

func TestArrayOfStructsWithUndefinedPropertyAccess(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/structs/undefined-property-access.french.swa",
			"property (Name) not found in struct Engineer\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/structs/undefined-property-access.soussou.swa",
			"se (Name) mu na fokhi Engineer kui\n")
	})
}

func TestArrayOfStructsWithUndefinedPropertyInitialization(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/structs/undefined-property-initialization.french.swa",
			"Property with name (Name) does not exist on struct Engineer\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/structs/undefined-property-initialization.soussou.swa",
			"Property with name (Name) does not exist on struct Engineer\n")
	})
}

func TestArrayOfStructsWithUndefined(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/structs/undefined-struct.english.swa",
			"struct named Engineer does not exist in symbol table\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/structs/undefined-struct.soussou.swa",
			"fokhi Engineer maboroma\n")
	})
}

func TestArrayOfStructs(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/structs/source.french.swa",
			"(nom: Pathe, age: 40, taille: 1.80, technos: Ruby, Rust, Go) (nom: Lucien, age: 24, taille: 1.81, technos: Typescript, HTML, Css) (nom: Manel, age: 25, taille: 1.82, technos: Typescript, Ruby) (nom: Bintou, age: 28, technos: Javascript, Css, HTML)")
	})

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/structs/source.english.swa",
			"(nom: Pathe, age: 40, taille: 1.80, stack: Ruby, Rust, Go) (nom: Lucien, age: 24, taille: 1.81, stack: Typescript, HTML, Css) (nom: Manel, age: 25, taille: 1.82, stack: Typescript, Ruby) (nom: Bintou, age: 28, taille: 1.83, stack: Javascript, Css, HTML) ")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/structs/source.soussou.swa",
			"(nom: Pathe, age: 40, taille: 1.80, technos: Ruby, Rust, Go) (nom: Lucien, age: 24, taille: 1.81, technos: Typescript, HTML, Css) (nom: Manel, age: 25, taille: 1.82, technos: Typescript, Ruby) (nom: Bintou, age: 28, technos: Javascript, Css, HTML)")
	})
}

func TestArrayOfStrings(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/strings/source.french.swa",
			"valeurs dans le tableau: (abc),(efg),(ijk)")
	})

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/strings/source.english.swa",
			"values in the array: (abc),(efg),(ijk)")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/strings/source.soussou.swa",
			"valeurs dans le tableau: (abc),(efg),(ijk)")
	})
}

func TestArraysInPrintStatement(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/print/source.french.swa",
			"Les valeurs dans le tableau sont: 1 2 3 4 5")
	})

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/print/source.english.swa",
			"Array values are: 1 2 3 4 5")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/print/source.soussou.swa",
			"Les valeurs dans le tableau sont: 1 2 3 4 5")
	})
}

func TestArrays(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/ok/source.french.swa",
			"")
	})

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/ok/source.english.swa",
			"")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./arrays/ok/source.soussou.swa",
			"")
	})
}

func TestArrayIndexOutOfBounds(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/out-of-bounds/source.french.swa",
			"L'element a la position (%!s(int=5)) depasse les limites du tableau (tableau)\n")
	})

	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/out-of-bounds/source.english.swa",
			"Element at index (%!s(int=5)) does not exist in array (array)\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/out-of-bounds/source.soussou.swa",
			"Se kui (%!s(int=5)) mu na tableau (tableau) kui\n")
	})
}

func TestArrayAccessWithPrefixExpression(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/prefix-expression/source.french.swa",
			"Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (-4)\n")
	})

	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/prefix-expression/source.english.swa",
			"Only numbers are supported as array index, current: (-4)\n")
	})

	t.Run("Soussou", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/prefix-expression/source.soussou.swa",
			"Konti nan lanné tableau kui, yakosi: (-4)\n")
	})
}

func TestArrayOfStructsAsFuncParam(t *testing.T) {
	// TODO find and fix the bug
	t.Skip("THERE IS A KNOWN BUG HERE THAT MUST BE FIXED")
	t.Run("Search Occurrence", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/structs/as-function-params.swa",
			"X=1 occurs 2 times, X=5 occurs 2 times, X=99 occurs 0 times",
		)
	})

	t.Run("High Earners", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/structs/count-high-earners.swa",
			"Employees with salary above 50000: 2\n",
		)
	})
}

func TestSortingAlgorithms(t *testing.T) {
	t.Run("QSort", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/sorting/qsort.swa",
			"Sorted List: 1 2 3 4 5 ",
		)
	})

	t.Run("Selection", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/sorting/selection.swa",
			"Sorted Array: 5 8 10 11 22 25 34 64 77 90 \n",
		)
	})

	t.Run("Bubble", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/sorting/bubble.swa",
			"Sorted Array: 4 8 11 12 22 25 34 64 77 90 \n",
		)
	})
}
func TestArrayZeroValues(t *testing.T) {
	t.Run("Array of ints", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/zero-values/int.swa",
			"- [3]int with zero value\n0 0 0\n\n",
		)
	})

	t.Run("Array of floats", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/zero-values/float.swa",
			"- [3]float with zero value\n0.00 0.00 0.00\n\n",
		)
	})

	t.Run("Array of strings", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/zero-values/string.swa",
			"- [3]string with zero value\n'' '' ''\n\n",
		)
	})

	t.Run("Array of structs", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./arrays/zero-values/struct.swa",
			"- [3]Dimension with zero value\n(Index: 0) Length: 0.00 Height: 0\n(Index: 1) Length: 0.00 Height: 0\n(Index: 2) Length: 0.00 Height: 0\n",
		)
	})

	t.Run("Size Too large", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/zero-values/size-limits.swa",
			"ArraySize (10000) too big for zero value initialization, max is 1000\n",
		)
	})
}

func TestArrayAccessWithSymbolExpression(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/symbol-expression/source.french.swa",
			"variable named x does not exist in symbol table\n")
	})

	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./arrays/symbol-expression/source.english.swa",
			"variable named x does not exist in symbol table\n")
	})
}
