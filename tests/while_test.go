package tests

import (
	"testing"
)

func TestWhileStatement(t *testing.T) {
	t.Run("French", func(t *testing.T) {
		t.Run("Iterate from ten to zero", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/from-ten-to-zero.french.swa",
				"10 9 8 7 6 5 4 3 2 1 ")
		})

		t.Run("Iterate from zero to ten", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/from-zero-to-ten.french.swa",
				"0 1 2 3 4 5 6 7 8 9 ")
		})

		t.Run("Iterate from zero to ten included", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/from-zero-to-ten-included.french.swa",
				"0 1 2 3 4 5 6 7 8 9 10 ")
		})
	})

	t.Run("English", func(t *testing.T) {
		t.Run("Iterate from ten to zero", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/from-ten-to-zero.english.swa",
				"10 9 8 7 6 5 4 3 2 1 ")
		})

		t.Run("Iterate from zero to ten", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/from-zero-to-ten.english.swa",
				"0 1 2 3 4 5 6 7 8 9 ")
		})

		t.Run("Iterate array", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/array.english.swa",
				"5 15 25 35 45 55 65 75 85 95 ")
		})

		t.Run("Display odd numbers from zero to ten", func(t *testing.T) {
			tests := []MultiDialectTest{
				{
					name:                    "english",
					inputPath:               "./while/odd-numbers.swa",
					expectedExecutionOutput: "1 3 5 7 9 ",
				},
				{
					name:                    "french",
					inputPath:               "./while/odd-numbers.french.swa",
					expectedExecutionOutput: "1 3 5 7 9 ",
				},
				{
					name:                    "soussou",
					inputPath:               "./while/odd-numbers.soussou.swa",
					expectedExecutionOutput: "1 3 5 7 9 ",
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					NewSuccessfulCompileRequest(t,
						test.inputPath,
						test.expectedExecutionOutput)
				})
			}
		})
		t.Run("Display even numbers from zero to ten", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/even-numbers.swa",
				"0 2 4 6 8 ")
		})

		t.Run("Iterate from zero to ten included", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./while/from-zero-to-ten-included.english.swa",
				"0 1 2 3 4 5 6 7 8 9 10 ")
		})
	})

	//	t.Run("Odd Numbers func", func(t *testing.T) {
	//		tests := []MultiDialectTest{
	//			{
	//				name:                    "english",
	//				inputPath:               "./while/odd-numbers-bool-func.swa",
	//				expectedExecutionOutput: "1 3 5 7 9 ",
	//			},
	//			{
	//				name:                    "soussou",
	//				inputPath:               "./while/odd-numbers-bool-func.soussou.swa",
	//				expectedExecutionOutput: "1 3 5 7 9 ",
	//			},
	//			{
	//				name:                    "french",
	//				inputPath:               "./while/odd-numbers-bool-func.french.swa",
	//				expectedExecutionOutput: "1 3 5 7 9 ",
	//			},
	//		}
	//
	//		for _, test := range tests {
	//			t.Run(test.name, func(t *testing.T) {
	//				NewSuccessfulCompileRequest(t,
	//					test.inputPath,
	//					test.expectedExecutionOutput)
	//			})
	//		}
	//	})
}
