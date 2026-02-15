package tests

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	t.Run("Function that returns only error", func(t *testing.T) {
		t.Run("english", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/error.english.swa",
				"Division by zero error (dividend is zero)\nNo division error",
			)
		})
		t.Run("soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/error.soussou.swa",
				"Division by zero error (dividend is zero)\nNo division error",
			)
		})

		t.Run("french", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/error.french.swa",
				"Division by zero erreur (dividend is zero)\nNo division erreur",
			)
		})

		t.Run("soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/error.soussou.swa",
				"Division by zero error (dividend is zero)\nNo division error",
			)
		})
	})

	t.Run("Declare external function", func(t *testing.T) {
		t.Run("english", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/declaration.english.swa",
				"",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/declaration.french.swa",
				"",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/declaration.soussou.swa",
				"",
			)
		})
	})

	t.Run("English", func(t *testing.T) {

		t.Run("static arrays as function parameter", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/array-sum.english.swa",
				"sum: 15, sumf: 15.00",
			)
		})

		t.Run("Function taking struct as argument by reference", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/struct.source.english.swa",
				"age: 40, height: 1.80, name: Pathe SENE",
			)
		})

		t.Run("Function taking struct as argument by value", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/struct-as-value.source.english.swa",
				"age: 0, age of copy: 40\nheight: 0.00, height of copy: 1.80\nname: Pathe, name of copy: Pathe SENE\n",
			)
		})

		t.Run("Substract", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/substract.integer.english.swa",
					"10 - 5 = 5",
				)
			})

			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/substract.float.english.swa",
					"10.5 - 5.25 = 5.25",
				)
			})
		})

		t.Run("Function call params", func(t *testing.T) {
			t.Run("Member expression", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/pass-member-expression-as-param.swa",
					"10.25 + 4.75 = 15.00",
				)
			})

			t.Run("Array access expression", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/pass-array-access-expression-as-param.swa",
					"10.25 + 4.75 = 15.00",
				)
			})

			t.Run("Incomplete Array access expression", func(t *testing.T) {
				NewFailedCompileRequest(t,
					"./functions/pass-incomplete-array-of-structs-access-expression-as-param.swa",
					"Struct property value type should be set -- Also cannot index array item that is a struct\n",
				)
			})

			t.Run("Array of structs access expression", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/pass-array-of-structs-access-expression-as-param.swa",
					"10.25 + 4.75 = 15.00",
				)
			})

			t.Run("Function call expression", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/pass-function-call-as-param.swa",
					"10.25 + 4.75 = 15.00",
				)
			})

			t.Run("Array initialization expression", func(t *testing.T) {
				// TODO
			})

			t.Run("Struct initialization expression", func(t *testing.T) {
				// TODO
			})
		})

		t.Run("Add", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/add.float.english.swa",
					"10.25 + 4.75 = 15.00",
				)
			})

			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/add.integer.english.swa",
					"10 + 5 = 15",
				)
			})
		})

		t.Run("Divide", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/divide.integer.english.swa",
					"10 / 5 = 2",
				)
			})

			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/divide.float.english.swa",
					"10.0 / 5.0 = 2.000000",
				)
			})

		})

		t.Run("Multiply", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/multiply.float.english.swa",
					"10.02, 5.05 = 50.601000",
				)
			})

			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/multiply.integer.english.swa",
					"10 * 5 = 50",
				)
			})
		})

	})

	t.Run("French", func(t *testing.T) {
		t.Run("Substract", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/substract.integer.french.swa",
					"10 - 5 = 5",
				)
			})

			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/substract.float.french.swa",
					"10.5 - 5.25 = 5.25",
				)
			})
		})

		t.Run("Add", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/add.float.french.swa",
					"10 + 5 = 15",
				)
			})
			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/add.integer.french.swa",
					"10 + 5 = 15",
				)
			})
		})

		t.Run("Divide", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/divide.float.french.swa",
					"10.0 / 5.0 = 2.000000",
				)
			})
			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/divide.integer.french.swa",
					"10 / 5 = 2",
				)
			})
		})

		t.Run("Multiply", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/multiply.float.french.swa",
					"10.02, 5.05 = 50.601000",
				)
			})

			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/multiply.integer.french.swa",
					"10 * 5 = 50",
				)
			})
		})

		t.Run("Soussou", func(t *testing.T) {
			t.Run("Add", func(t *testing.T) {
				t.Run("Integer", func(t *testing.T) {
					NewSuccessfulCompileRequest(t,
						"./functions/add.integer.soussou.swa",
						"10 + 5 = 15",
					)
				})
			})
			t.Run("Subtract", func(t *testing.T) {
				t.Run("Integer", func(t *testing.T) {
					NewSuccessfulCompileRequest(t,
						"./functions/substract.integer.soussou.swa",
						"10 - 5 = 5",
					)
				})
			})
		})
	})
}
