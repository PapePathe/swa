package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructZeroValues(t *testing.T) {
	t.Run("Structs", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/zero-values.swa",
			ExpectedExecutionOutput: "- Dimension with zero value\nLength: 0.00, Width: 0.00, Height: 0, Name: \n\n- Object with zero value\nLength: 0.00, Width: 0.00, Height: 0, Name: , Shape: \n\n- Block with zero value\nTop(Length: 0.00, Width: 0.00, Height: 0, Name: , Shape: )\nBottom(Length: 0.00, Width: 0.00, Height: 0, Name: , Shape: )\n",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructWithUnknownType(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.french.swa",
			ExpectedOutput: "struct named unknown_type does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:      "./structs/unknown-property-type/source.english.swa",
			ExpectedOutput: "struct named unknown_type does not exist in symbol table\n",
			T:              t,
		}

		assert.Error(t, req.Compile())
	})
}

func TestStructPropertyAssignment(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath: "./structs/assignment/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath: "./structs/assignment/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructAll(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/all/source.english.swa",
			OutputPath:              "15a6540e-5df1-4fbd-b3a1-a2edb12d117b",
			ExpectedExecutionOutput: "Name: (Pathe), TechStack: (Ruby, Rust, Go), Age: (40), Height: (1.80)",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()

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
	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath: "./structs/return-expression/source.english.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})

	t.Run("French", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath: "./structs/return-expression/source.french.swa",
			T:         t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructBasic(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/nested/basic/source.english.swa",
			ExpectedExecutionOutput: "Name: Pathe, City: Dakar, ZipCode: 12000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructInline(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./structs/nested/inline/source.english.swa",
			ExpectedExecutionOutput: "Name: Sene, City: Thies, ZipCode: 21000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructFieldAccess(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./structs/nested/field-access/source.english.swa",
			ExpectedExecutionOutput: "Name: Pathe, Age: 30, City: Dakar, ZipCode: 12000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructTriple(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/nested/triple/source.english.swa",
			ExpectedExecutionOutput: "X: 10, Y: 20, getX: 10, setX: 999",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructAssignment(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./structs/nested/assignment/source.english.swa",
			ExpectedExecutionOutput: "City: Dakar, ZipCode: 99999, Code+1: 100000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestNestedStructArrays(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/nested/arrays/source.english.swa",
			ExpectedExecutionOutput: "Person 1 City: Dakar, Person 2 ZipCode: 21000",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructArrayEmbedded(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/array-embedded/source.english.swa",
			ExpectedExecutionOutput: "Hello: 100 200 300",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructArrayPointer(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/array-pointer/source.english.swa",
			ExpectedExecutionOutput: "Hello: 100 200 300",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructMultiArrays(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()
		req := CompileRequest{
			InputPath:               "./structs/multi-arrays/source.english.swa",
			ExpectedExecutionOutput: "arr1[0]=10 arr1[1]=20 arr2[0]=100 arr2[2]=300 value=42",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructMultiPointers(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/multi-pointers/source.english.swa",
			ExpectedExecutionOutput: "ptr1[0]=10 ptr1[1]=20 ptr2[0]=100 ptr2[2]=300 value=99",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructMixedArrays(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/mixed-arrays/source.english.swa",
			ExpectedExecutionOutput: "embedded[0]=10 embedded[1]=20 pointer[0]=100 pointer[2]=300",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructFuncParamEmbedded(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/func-param-embedded/source.english.swa",
			ExpectedExecutionOutput: "Result: 65",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructFuncParamPointer(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		req := CompileRequest{
			InputPath:               "./structs/func-param-pointer/source.english.swa",
			ExpectedExecutionOutput: "Result: 65",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructWithEmbeddedArrayOfStructs(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Parallel()

		req := CompileRequest{
			InputPath:               "./structs/array-of-structs/source.english.swa",
			ExpectedExecutionOutput: "p[0].x=1 p[0].y=1\np[1].x=2 p[1].y=2\np[2].x=3 p[2].y=3\n",
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}

func TestStructInitialization(t *testing.T) {
	t.Run("Init array field with a symbol", func(t *testing.T) {
		t.Parallel()

		input := "./structs/all/init-field-from-array.swa"
		expected := "Backing[0]: 10, Wrapper.Data[0]: 10 \nBacking[1]: 20, Wrapper.Data[1]: 20 \nBacking[2]: 30, Wrapper.Data[2]: 30 \nBacking[3]: 0, Wrapper.Data[3]: 0 \nBacking[4]: 0, Wrapper.Data[4]: 0 \nBacking[5]: 0, Wrapper.Data[5]: 0 \nBacking[6]: 0, Wrapper.Data[6]: 0 \nBacking[7]: 0, Wrapper.Data[7]: 0 \nBacking[8]: 0, Wrapper.Data[8]: 0 \nBacking[9]: 0, Wrapper.Data[9]: 0 \n"
		req := CompileRequest{
			InputPath:               input,
			ExpectedExecutionOutput: expected,
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
	t.Run("Init int field with a symbol", func(t *testing.T) {
		t.Parallel()

		input := "./structs/all/init-field-from-number.swa"
		expected := "1000"
		req := CompileRequest{
			InputPath:               input,
			ExpectedExecutionOutput: expected,
			T:                       t,
		}

		req.AssertCompileAndExecute()
	})
}
