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
			ExpectedLLIR:   "./arrays/ok/source.french.ll",
			OutputPath:     "474c1abf-c877-455b-b6b9-09bdbc9b3d47",
			ExpectedOutput: "",
			T:              t,
		}
		defer req.Cleanup()

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/ok/source.english.swa",
			ExpectedLLIR:   "./arrays/ok/source.english.ll",
			OutputPath:     "f32e0fc5-a990-491e-8c5b-89a787866f1a",
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
			ExpectedLLIR:   "./arrays/out-of-bounds/source.french.ll",
			OutputPath:     "78cd171d-69e5-44c9-84c1-85c3939d00a8",
			ExpectedOutput: "Element at index (5) does not exist in array (tableau)\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/out-of-bounds/source.english.swa",
			ExpectedLLIR:   "./arrays/out-of-bounds/source.english.ll",
			OutputPath:     "66ebfb11-c40a-41d5-84c1-de0b710c9a88",
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
			ExpectedLLIR:   "./arrays/prefix-expression/source.french.ll",
			OutputPath:     "a6fe34d3-83f8-4dc5-a03e-c786ce604da2",
			ExpectedOutput: "Only numbers are supported as array index, current: ({{- MINUS MINUS} {4}})\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/prefix-expression/source.english.swa",
			ExpectedLLIR:   "./arrays/prefix-expression/source.english.ll",
			OutputPath:     "89e28b67-a4e4-4a45-855e-eb14ca790ba6",
			ExpectedOutput: "Only numbers are supported as array index, current: ({{- MINUS MINUS} {4}})\n",
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
			ExpectedLLIR:   "./arrays/symbol-expression/source.french.ll",
			OutputPath:     "31e357ba-aa48-4999-b437-289b3cf1cac4",
			ExpectedOutput: "Only numbers are supported as array index, current: ({x})\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./arrays/symbol-expression/source.english.swa",
			ExpectedLLIR:   "./arrays/symbol-expression/source.english.ll",
			OutputPath:     "602748f7-8ced-4180-8247-f27c54524196",
			ExpectedOutput: "Only numbers are supported as array index, current: ({x})\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}
