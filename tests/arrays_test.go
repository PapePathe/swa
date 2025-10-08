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
