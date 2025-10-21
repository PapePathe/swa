package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArraysInPrintStatement(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./arrays/print/source.french.swa",
			ExpectedOutput:          "",
			ExpectedExecutionOutput: "La valeur a la position 0 est 1,La valeur a la position 1 est 2",
			T:                       t,
		}
		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	// t.Run("English", func(t *testing.T) {
	// 	req := CompileRequest{
	// 		InputPath:      "./arrays/ok/source.english.swa",
	// 		ExpectedOutput: "",
	// 		T:              t,
	// 	}
	// 	defer req.Cleanup()

	// 	req.AssertCompileAndExecute()
	// })
}

func TestArrays(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/ok/source.french.swa",
			ExpectedOutput: "",
			T:              t,
		}
		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/ok/source.english.swa",
			ExpectedOutput: "",
			T:              t,
		}
		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})
}

func TestArrayIndexOutOfBounds(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/out-of-bounds/source.french.swa",
			ExpectedOutput: "L'element a la position (5) depasse les limites du tableau (tableau)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/out-of-bounds/source.english.swa",
			ExpectedOutput: "Element at index (5) does not exist in array (array)\n",
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
			ExpectedOutput: "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (- 4)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.english.swa",
			ExpectedOutput: "Only numbers are supported as array index, current: (- 4)\n",
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
			ExpectedOutput: "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (x)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.english.swa",
			ExpectedOutput: "Only numbers are supported as array index, current: (x)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}
