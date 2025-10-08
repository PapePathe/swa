package tests

import (
	"testing"
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
