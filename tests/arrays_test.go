package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			ExpectedOutput: "Element at index (5) does not exist in array (tableau)\n",
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
			ExpectedOutput: "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: ([[- 4]])\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.english.swa",
			ExpectedOutput: "Only numbers are supported as array index, current: ([[- 4]])\n",
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
			ExpectedOutput: "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: ([[x]])\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.english.swa",
			ExpectedOutput: "Only numbers are supported as array index, current: ([[x]])\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}
