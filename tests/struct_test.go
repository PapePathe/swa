package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructWithUnknownType(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.french.swa",
			ExpectedOutput: "struct proprerty type ({unknown_type}) not supported\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.english.swa",
			ExpectedOutput: "struct proprerty type ({unknown_type}) not supported\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestStructPropertyAssignment(t *testing.T) {
	t.Parallel()

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/assignment/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/assignment/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructAll(t *testing.T) {
	t.Parallel()

	t.Run("Soussou", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.soussou.swa",
			ExpectedExecutionOutput: "Nom: (Pathe), Stack Technique: (Ruby, Rust, Go), Age: (40)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.english.swa",
			ExpectedExecutionOutput: "Name: (Pathe), TechStack: (Ruby, Rust, Go), Age: (40), Height: (1.80)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.french.swa",
			ExpectedExecutionOutput: "Nom: (Pathe), Stack Technique: (Ruby, Rust, Go), Age: (40), Taille: (1.80)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructPropertyInReturnExpression(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/return-expression/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath: "./structs/return-expression/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}
