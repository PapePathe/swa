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
			ExpectedOutput: "struct named unknown_type does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.english.swa",
			ExpectedOutput: "struct named unknown_type does not exist in symbol table\n",
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

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.english.swa",
			OutputPath:              "15a6540e-5df1-4fbd-b3a1-a2edb12d117b",
			ExpectedExecutionOutput: "Name: (Pathe), TechStack: (Ruby, Rust, Go), Age: (40), Height: (1.80)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/all/source.french.swa",
			OutputPath:              "83b453cc-bd0b-432d-96d5-f0192dc2c2d3",
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

func TestNestedStructBasic(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/nested/basic/source.english.swa",
			ExpectedExecutionOutput: "Name: Pathe, City: Dakar, ZipCode: 12000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructInline(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/nested/inline/source.english.swa",
			ExpectedExecutionOutput: "Name: Sene, City: Thies, ZipCode: 21000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructFieldAccess(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/nested/field-access/source.english.swa",
			ExpectedExecutionOutput: "Name: Pathe, Age: 30, City: Dakar, ZipCode: 12000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructTriple(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/nested/triple/source.english.swa",
			ExpectedExecutionOutput: "X: 10, Y: 20, getX: 10, setX: 999",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructAssignment(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/nested/assignment/source.english.swa",
			ExpectedExecutionOutput: "City: Dakar, ZipCode: 99999, Code+1: 100000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
