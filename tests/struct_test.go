package tests

import (
	"testing"
)

func TestEmptyStruct(t *testing.T) {
	tests := []MultiDialectTest{
		{
			name:           "English",
			inputPath:      "./structs/empty.swa",
			expectedOutput: "Struct (Point) must have at least one field\n",
		},
		{
			name:           "french",
			inputPath:      "./structs/empty.french.swa",
			expectedOutput: "La structure (Point) doit avoir au moins un champ\n",
		},
		// TODO
		// * Add test cas for wolof
		// * Add test cas for soussou
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			NewFailedCompileRequest(t,
				test.inputPath,
				test.expectedOutput)
		})
	}

}

func TestStructZeroValues(t *testing.T) {
	expected := "- Dimension with zero value\nLength: 0.00, Width: 0.00, Height: 0, Name: \n\n- Object with zero value\nLength: 0.00, Width: 0.00, Height: 0, Name: , Shape: \n\n- Block with zero value\nTop(Length: 0.00, Width: 0.00, Height: 0, Name: , Shape: )\nBottom(Length: 0.00, Width: 0.00, Height: 0, Name: , Shape: )\n"

	t.Run("Structs with explicit zero initializer", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/zero-values.swa",
			expected)
	})

	t.Run("Structs with implicit zero initializer", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/zero-values-implicit.swa",
			expected)
	})

	t.Run("Array of structs", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/zero-values-array.swa",
			"- [1]Path with zero value\nx1: 0 y1:0")
	})
}

func TestStructWithUnknownType(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./structs/unknown-property-type/source.french.swa",
			"struct named unknown_type does not exist in symbol table\n")
	})

	t.Run("English", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./structs/unknown-property-type/source.english.swa",
			"struct named unknown_type does not exist in symbol table\n")
	})
}

func TestStructPropertyAssignment(t *testing.T) {
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
	t.Run("English", func(t *testing.T) {
		expected := "Name: (Pathe), TechStack: (Ruby, Rust, Go), Age: (40), Height: (1.80)"
		NewSuccessfulCompileRequest(t,
			"./structs/all/source.english.swa",
			expected)
	})

	t.Run("French", func(t *testing.T) {
		expected := "Nom: (Pathe), Stack Technique: (Ruby, Rust, Go), Age: (40), Taille: (1.80)"
		NewSuccessfulCompileRequest(t,
			"./structs/all/source.french.swa",
			expected)
	})
}

func TestStructPropertyInReturnExpression(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/return-expression/source.english.swa",
			"")
	})

	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/return-expression/source.french.swa",
			"")
	})
}

func TestNestedStructBasic(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/nested/basic/source.english.swa",
			"Name: Pathe, City: Dakar, ZipCode: 12000")
	})
}

func TestNestedStructInline(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/nested/inline/source.english.swa",
			"Name: Sene, City: Thies, ZipCode: 21000")
	})
}

func TestNestedStructFieldAccess(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/nested/field-access/source.english.swa",
			"Name: Pathe, Age: 30, City: Dakar, ZipCode: 12000")
	})
}

func TestNestedStructTriple(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/nested/triple/source.english.swa",
			"X: 10, Y: 20, getX: 10, setX: 999")
	})
}

func TestNestedStructAssignment(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/nested/assignment/source.english.swa",
			"City: Dakar, ZipCode: 99999, Code+1: 100000")
	})
}

func TestNestedStructArrays(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/nested/arrays/source.english.swa",
			"Person 1 City: Dakar, Person 2 ZipCode: 21000")
	})
}

func TestStructArrayEmbedded(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/array-embedded/source.english.swa",
			"Hello: 100 200 300")
	})
}

func TestStructArrayPointer(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/array-pointer/source.english.swa",
			"Hello: 100 200 300")
	})
}

func TestStructMultiArrays(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/multi-arrays/source.english.swa",
			"arr1[0]=10 arr1[1]=20 arr2[0]=100 arr2[2]=300 value=42")
	})
}

func TestStructMultiPointers(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/multi-pointers/source.english.swa",
			"ptr1[0]=10 ptr1[1]=20 ptr2[0]=100 ptr2[2]=300 value=99")
	})
}

func TestStructMixedArrays(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/mixed-arrays/source.english.swa",
			"embedded[0]=10 embedded[1]=20 pointer[0]=100 pointer[2]=300")
	})
}

func TestStructFuncParamEmbedded(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/func-param-embedded/source.english.swa",
			"Result: 65")
	})
}

func TestStructFuncParamPointer(t *testing.T) {
	t.Parallel()

	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/func-param-pointer/source.english.swa",
			"Result: 65")
	})
}

func TestStructWithEmbeddedArrayOfStructs(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/array-of-structs/source.english.swa",
			"p[0].x=1 p[0].y=1\np[1].x=2 p[1].y=2\np[2].x=3 p[2].y=3\n")
	})
}

func TestStructInitialization(t *testing.T) {
	t.Run("Init array field with a symbol", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/all/init-field-from-array.swa",
			"Backing[0]: 10, Wrapper.Data[0]: 10 \nBacking[1]: 20, Wrapper.Data[1]: 20 \nBacking[2]: 30, Wrapper.Data[2]: 30 \nBacking[3]: 0, Wrapper.Data[3]: 0 \nBacking[4]: 0, Wrapper.Data[4]: 0 \nBacking[5]: 0, Wrapper.Data[5]: 0 \nBacking[6]: 0, Wrapper.Data[6]: 0 \nBacking[7]: 0, Wrapper.Data[7]: 0 \nBacking[8]: 0, Wrapper.Data[8]: 0 \nBacking[9]: 0, Wrapper.Data[9]: 0 \n")
	})
	t.Run("Init int field with a symbol", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./structs/all/init-field-from-number.swa",
			"1000")
	})
}

// TODO test is not working on github
// Figure out why
// func TestStructPropertyBool(t *testing.T) {
// 	NewSuccessfulCompileRequest(t,
// 		"./structs/booleans.swa",
// 		"Door 1 is open ? : 0\nDoor 2 is open ? : 0\nDoor 3 is open ? : 1\nDoor 4 is open ? : 0\n",
// 	)
// }
