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
			t.Run("LLVM", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/array-sum.english.swa",
					"sum: 15, sumf: 15.00",
				)
			})
			t.Run("SWA", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/array-sum.english.swa",
					"sum: 15, sumf: 15.00",
				)
			})
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

			t.Run("DeepSeek Functions", func(t *testing.T) {
				t.Run("Basics", func(t *testing.T) {
					t.Run("Test 1: Simple function", func(t *testing.T) {
						NewSuccessfulCompileRequest(t, "./functions/test1_simple_function_with_no_parameters_returning_int.english.swa", "one() = 1")
					})
					t.Run("Test 2: One parameter", func(t *testing.T) {
						NewSuccessfulCompileRequest(t, "./functions/test2_function_with_one_parameter_call_with_constant.english.swa", "double(5) = 10")
					})
					t.Run("Test 3: Two parameters", func(t *testing.T) {
						NewSuccessfulCompileRequest(t, "./functions/test3_function_with_two_parameters_call_with_variables.english.swa", "add(3,4) = 7")
					})
					t.Run("Test 4: Return float", func(t *testing.T) {
						NewSuccessfulCompileRequest(t, "./functions/test4_function_returning_float.english.swa", "pi() = 3.141590")
					})
					t.Run("Test 5: Return string", func(t *testing.T) {
						NewSuccessfulCompileRequest(t, "./functions/test5_function_returning_string.english.swa", "greet() = hello")
					})
					t.Run("Test 6: Global side effect", func(t *testing.T) {
						NewSuccessfulCompileRequest(t, "./functions/test6_function_with_side_effect_modifying_global.english.swa", "global = 42")
					})
					t.Run("Test 8: Array as reference", func(t *testing.T) {
						NewSuccessfulCompileRequest(t,
							"./functions/test8_function_that_takes_array_by_value_if_allowed_arrays_might_be_passed_by_reference_need_to_test.english.swa",
							"a[0] = 0")
					})
				})
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
			t.Run("Float.X", func(t *testing.T) {
				NewSuccessfulXCompileRequest(CompileRequest{
					T:         t,
					InputPath: "./functions/add.float.english.swa"})
			})

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

			t.Run("Integer.X", func(t *testing.T) {
				NewSuccessfulXCompileRequest(CompileRequest{
					T:         t,
					InputPath: "./functions/add.integer.english.swa"})
			})
		})

		t.Run("Divide", func(t *testing.T) {
			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/divide.integer.english.swa",
					"10 / 5 = 2",
				)
			})

			t.Run("Integer.X", func(t *testing.T) {
				NewSuccessfulXCompileRequest(CompileRequest{
					T:         t,
					InputPath: "./functions/divide.integer.english.swa"})
			})

			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/divide.float.english.swa",
					"10.0 / 5.0 = 2.000000",
				)
			})

			t.Run("Float.X", func(t *testing.T) {
				NewSuccessfulXCompileRequest(CompileRequest{
					T:         t,
					InputPath: "./functions/divide.float.english.swa"})
			})
		})

		t.Run("Multiply", func(t *testing.T) {
			t.Run("Float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/multiply.float.english.swa",
					"10.02, 5.05 = 50.601000",
				)
			})

			t.Run("Float.X", func(t *testing.T) {
				NewSuccessfulXCompileRequest(CompileRequest{
					T:         t,
					InputPath: "./functions/multiply.float.english.swa"})
			})

			t.Run("Integer", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/multiply.integer.english.swa",
					"10 * 5 = 50",
				)
			})

			t.Run("Integer.X", func(t *testing.T) {
				NewSuccessfulXCompileRequest(CompileRequest{
					T:         t,
					InputPath: "./functions/multiply.integer.english.swa"})
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

		t.Run("RLE Compression", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/rle.french.swa",
				"function named longueur does not exist in symbol table\n",
			)
		})

	})

	t.Run("DeepSeek Functions", func(t *testing.T) {
		t.Run("Basics", func(t *testing.T) {
			t.Run("Test 1: Simple function", func(t *testing.T) {
				NewSuccessfulCompileRequest(
					t,
					"./functions/test1_simple_function_with_no_parameters_returning_int.english.swa",
					"one() = 1")
			})
			t.Run("Test 2: One parameter", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/test2_function_with_one_parameter_call_with_constant.english.swa",
					"double(5) = 10")
			})
			t.Run("Test 3: Two parameters", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test3_function_with_two_parameters_call_with_variables.english.swa", "add(3,4) = 7")
			})
			t.Run("Test 4: Return float", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test4_function_returning_float.english.swa", "pi() = 3.141590")
			})
			t.Run("Test 5: Return string", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test5_function_returning_string.english.swa", "greet() = hello")
			})
		})

		t.Run("Recursion", func(t *testing.T) {
			t.Run("Test 11: Factorial", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/test11_recursion_factorial.english.swa", "5! = 120")
			})
			//	FIXME
			//  This is currently out of scope. We need to have a first
			//  pass to just parse function haaders and add them to
			//  the symbol table entry
			//
			//	t.Run("Test 12: Mutual recursion", func(t *testing.T) {
			//		NewSuccessfulCompileRequest(t,
			//			"./functions/test12_mutual_recursion_even_odd.english.swa",
			//			"isEven(4) = 1\nisOdd(4) = 0\n")
			//	})
			t.Run("Test 27: Direct recursion", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/test27_function_call_inside_its_own_definition_direct_recursion_already_tested_but_add_a_simple_one.english.swa",
					"5 4 3 2 1 0")
			})
		})

		t.Run("Composition and Arguments", func(t *testing.T) {
			t.Run("Test 13: Expression arguments", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test13_function_call_with_expression_arguments.english.swa", "add(2*2, 3+1) = 8")
			})
			t.Run("Test 14: Nested calls", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test14_nested_function_calls_composition.english.swa", "twice(square(3)) = 18")
			})
			t.Run("Test 17: Many parameters", func(t *testing.T) {
				NewSuccessfulCompileRequest(t,
					"./functions/test17_function_with_many_parameters_stress.english.swa",
					"sum10 = 55")
			})
		})

		t.Run("Structs and Arrays", func(t *testing.T) {
			t.Run("Test 15: Return struct", func(t *testing.T) {
				NewFailedCompileRequest(t,
					"./functions/test15_function_returning_a_struct.english.swa",
					"returning complex types (structs and arrays) from functions is currently not supported\n")
			})
			t.Run("Test 36: Return struct with array", func(t *testing.T) {
				NewFailedCompileRequest(t,
					"./functions/test36_function_that_returns_a_struct_containing_an_array.english.swa",
					"returning complex types (structs and arrays) from functions is currently not supported\n")
			})
		})

		t.Run("Control Flow and Side Effects", func(t *testing.T) {
			t.Run("Test 7: Call by value", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test7_function_that_modifies_parameter_call_by_value_should_not_affect_caller.english.swa", "a = 5")
			})
			t.Run("Test 39: Call in loop", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test39_function_call_inside_a_loop.english.swa", "i=0\ni=1\ni=2\ni=3\ni=4\n")
			})
			t.Run("Test 40: Dual side effects", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test40_function_call_with_side_effect_that_modifies_global_used_in_argument.english.swa", "r = 3, g = 2")
			})
		})

		t.Run("Initializers", func(t *testing.T) {
			t.Run("Test 49: Call in struct init", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test49_function_call_inside_struct_literal_initialization.english.swa", "w.val = 42")
			})
			t.Run("Test 50: Call in array init", func(t *testing.T) {
				NewSuccessfulCompileRequest(t, "./functions/test50_function_call_in_array_initializer.english.swa", "arr[0] = 0\narr[1] = 1\narr[2] = 4\narr[3] = 9\narr[4] = 16\n")
			})
		})

		t.Run("Expected Failures", func(t *testing.T) {
			t.Run("Test 21: Mismatched return type", func(t *testing.T) {
				NewFailedCompileRequest(t, "./functions/test21_function_call_with_mismatched_return_type_should_not_compile.english.swa", "expected String but got Number\n")
			})
			t.Run("Test 22: Too few arguments", func(t *testing.T) {
				NewFailedCompileRequest(t, "./functions/test22_function_call_with_too_few_arguments_should_not_compile.english.swa", "function two_args expect 2 arguments but was given 1\n")
			})
			t.Run("Test 23: Too many arguments", func(t *testing.T) {
				NewFailedCompileRequest(t, "./functions/test23_function_call_with_too_many_arguments_should_not_compile.english.swa", "function two_args expect 2 arguments but was given 3\n")
			})
			t.Run("Test 24: Wrong argument type", func(t *testing.T) {
				NewFailedCompileRequest(t, "./functions/test24_function_call_with_wrong_argument_type_should_not_compile.english.swa", "expected argument of type IntegerType(32 bits) but got PointerType(Reference)\n")
			})
			t.Run("Test 25: Missing return", func(t *testing.T) {
				NewFailedCompileRequest(t, "./functions/test25_function_with_missing_return_statement_should_not_compile.english.swa", "function has no return statement\n")
			})
			t.Run("Test 26: Duplicate function", func(t *testing.T) {
				NewFailedCompileRequest(t, "./functions/test26_function_redeclaration_should_not_compile.english.swa", "function named foo already exists in symbol table\n")
			})
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
}
