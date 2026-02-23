package tests

import (
	"testing"
)

// TestDeepseekPlaintext contains test cases extracted from deepseek_plaintext_20260215_c430ab.txt.
// Each test corresponds to one test program in that file.
func TestDeepseekPlaintext(t *testing.T) {

	t.Run("Arithmetic and Numerics", func(t *testing.T) {

		// Test 1: Integer addition overflow (2^31-1 + 1)
		// On 64-bit systems the result is -2147483648 (wraps); on 64-bit it may differ.
		t.Run("Test 1: Integer addition overflow", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test1_integer_addition_overflow.english.swa",
				"max+1 = -2147483648",
			)
		})

		// Test 2: Float division producing infinity
		// FIXME division by zero should error
		t.Run("Test 2: Float division infinity", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test2_float_division_infinity.english.swa",
				"1.0/0.0 = inf",
			)
		})

		// Test 3: Float NaN (0.0/0.0)
		t.Run("Test 3: Float NaN", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test3_float_nan.english.swa",
				"0.0/0.0 = nan",
			)
		})

		// Test 4: Modulo by zero – runtime error
		//	FIXME module by zero is undefined behaviour
		//	t.Run("Test 4: Modulo by zero", func(t *testing.T) {
		//		NewSuccessfulCompileRequest(t,
		//			"./functions/deepseek_test4_modulo_by_zero.english.swa",
		//			"",
		//		)
		//	})

		// Test 33: Operator precedence (2+3*4 should be 14)
		t.Run("Test 33: Operator precedence", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test33_operator_precedence.english.swa",
				"2+3*4 = 14",
			)
		})

		// Test 41: Large integer literal (64-bit)
		t.Run("Test 41: Large integer literal", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test41_large_integer_literal.english.swa",
				"1234567890123456789 is greater than max value for int32\n",
			)
		})

		// Test 42: Float literal with many decimals
		t.Run("Test 42: Float many decimals", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test42_float_many_decimals.english.swa",
				"pi = 3.14159265358979311600",
			)
		})

		// Test 43: Modulo with negative divisor (7 % -3 = 1 in C semantics)
		t.Run("Test 43: Modulo negative divisor", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test43_modulo_negative_divisor.english.swa",
				"7 % -3 = 1",
			)
		})
	})

	t.Run("Literals and Special Syntax", func(t *testing.T) {

		// Test 23: Hexadecimal integer literal (0xFF = 255)
		//	FIXME add support for Hexadecimal numbers
		//	t.Run("Test 23: Hexadecimal literal", func(t *testing.T) {
		//		NewSuccessfulCompileRequest(t,
		//			"./functions/deepseek_test23_hexadecimal_literal.english.swa",
		//			"\x1b[33mexpected SEMI_COLON, but got IDENTIFIER at line 6\x1b[0m\n\n\x1b[34m6\x1b[0m \x1b[32mhex\x1b[0m \x1b[32m:\x1b[0m \x1b[32mint\x1b[0m \x1b[32m=\x1b[0m \x1b[32m0\x1b[0m \x1b[31mxFF\x1b[0m \n",
		//		)
		//	})

		// Test 24: Octal integer literal (077 = 63)
		//	FIXME add support for octal numbers
		//	t.Run("Test 24: Octal literal", func(t *testing.T) {
		//		NewSuccessfulCompileRequest(t,
		//			"./functions/deepseek_test24_octal_literal.english.swa",
		//			"077 = 63",
		//		)
		//	})

		// Test 25: Float literal with exponent (1.23e4 = 12300.0)
		t.Run("Test 25: Float exponent literal", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test25_float_literal_with_exponent.english.swa",
				"12300.000000",
			)
		})

		// Test 19: Empty statement (lone semicolon)
		t.Run("Test 19: Empty statement", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test19_empty_statement.english.swa",
				"nud handler expected for token SEMI_COLON and binding power 0 \n 7",
			)
		})

		// Test 20: Escape sequences in strings
		t.Run("Test 20: Escape sequences in strings", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test20_escape_sequences_in_strings.english.swa",
				"Hello\nWorld\tTab\\Backslash\"Quote",
			)
		})
	})

	t.Run("Arrays", func(t *testing.T) {

		// Test 5: Array out-of-bounds with variable index – runtime error
		t.Run("Test 5: Array index out of bounds with variable", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test5_array_index_out_of_bounds_with_variable.english.swa",
				"SWA_ERROR: index 5 out of bounds of array with size 5",
			)
		})

		// Test 6: Array index with negative variable – runtime error
		t.Run("Test 6: Array index negative variable", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test6_array_index_negative_variable.english.swa",
				"SWA_ERROR: index -1 out of bounds of array with size 5",
			)
		})

		// Test 7: Array initialization with fewer elements (remaining zero-initialized)
		t.Run("Test 7: Array init fewer elements", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test7_array_init_fewer_elements.english.swa",
				"arr[3] = 0, arr[4] = 0",
			)
		})

		// Test 37: Array equality – should fail to compile
		t.Run("Test 37: Array equality not allowed", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test37_array_equality.english.swa",
				"Invalid operand types for ICmp instruction\n  %8 = icmp eq [3 x i32] %6, %7\n\n",
			)
		})
	})

	t.Run("Structs", func(t *testing.T) {

		// Test 9: Struct missing field in init
		t.Run("Test 9: Struct field missing in init", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test9_struct_missing_field_in_init.english.swa",
				"",
			)
		})

		// Test 10: Nested struct missing inner field
		t.Run("Test 10: Nested struct missing inner field", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test10_nested_struct_missing_inner_field.english.swa",
				"",
			)
		})

		// Test 28: Struct with zero fields
		t.Run("Test 28: Struct with zero fields", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test28_struct_with_zero_fields.english.swa",
				"Struct (Empty) must have at least one field\n",
			)
		})

		// Test 36: Struct equality – should fail to compile
		t.Run("Test 36: Struct equality not allowed", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test36_struct_equality.english.swa",
				"Invalid operand types for ICmp instruction\n  %2 = icmp eq %Point %0, %1\n\n",
			)
		})
	})

	t.Run("Strings", func(t *testing.T) {

		// Test 21: String concatenation using + – may not be supported
		t.Run("Test 21: String concatenation", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test21_string_concatenation.english.swa",
				"Operation PLUS not supported on strings\n",
			)
		})

		// Test 35: String comparison with accented characters
		t.Run("Test 35: String comparison with accented chars", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test35_string_comparison_accented_chars.english.swa",
				"Not equal",
			)
		})
	})

	t.Run("Booleans and Control Flow", func(t *testing.T) {

		// Test 26: Boolean operations
		t.Run("Test 26: Boolean operations", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test26_boolean_operations.english.swa",
				"true",
			)
		})

		// Test 34: Short-circuit evaluation
		t.Run("Test 34: Short-circuit evaluation", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test34_short_circuit_evaluation.english.swa",
				"short-circuit works",
			)
		})

		// Test 17: Break statement inside while loop
		// FIXME implement break statement
		t.Run("Test 17: Break in while loop", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test17_break_in_while_loop.english.swa",
				"variable named break does not exist in symbol table\n",
			)
		})

		// Test 18: Continue statement inside while loop
		// Prints i for all i from 1 to 10 except 5
		// FIXME implement continue statement
		t.Run("Test 18: Continue in while loop", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test18_continue_in_while_loop.english.swa",
				"variable named continue does not exist in symbol table\n",
			)
		})
	})

	t.Run("Functions", func(t *testing.T) {

		// Test 13: Function missing return statement – compile error
		t.Run("Test 13: Function missing return", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test13_function_missing_return.english.swa",
				"function has no return statement\n",
			)
		})

		// Test 30: Function modifies parameter (call by value)
		t.Run("Test 30: Call by value no side effect", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test30_function_modifies_param_call_by_value.english.swa",
				"a = 5",
			)
		})

		// Test 31: Deep recursion (100000 levels)
		t.Run("Test 31: Deep recursion", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test31_deep_recursion_stack_overflow.english.swa",
				"100000",
			)
		})

		// Test 39: Empty function body – should cause compile error (no return)
		t.Run("Test 39: Empty function body", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test39_empty_function_body.english.swa",
				"function has no return statement\n",
			)
		})

		// Test 40: Function with many parameters (10 params, sum = 55)
		t.Run("Test 40: Many parameters", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test40_function_many_parameters.english.swa",
				"sum = 55",
			)
		})
	})

	t.Run("Type Errors (Expected Failures)", func(t *testing.T) {
		// Test 44: Type mismatch: int vs float in comparison
		t.Run("Test 44: Type mismatch int vs float", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/deepseek_test44_type_mismatch_int_float.english.swa",
				"equal",
			)
		})

		// Test 45: Type mismatch: bool vs string in comparison – compile error
		t.Run("Test 45: Type mismatch bool vs string", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test45_type_mismatch_bool_string.english.swa",
				"cannot coerce IntegerTypeKind and PointerTypeKind\n",
			)
		})

		// Test 46: Variable redeclaration in same scope – compile error
		t.Run("Test 46: Variable redeclaration", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test46_variable_redeclaration.english.swa",
				"variable a is already defined\n",
			)
		})

		// Test 47: Use variable before declaration – compile error
		t.Run("Test 47: Use before declaration", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test47_use_before_declaration.english.swa",
				"variable named x does not exist in symbol table\n",
			)
		})

		// Test 48: Function call with wrong argument type – compile error
		t.Run("Test 48: Wrong argument type", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test48_wrong_argument_type.english.swa",
				"expected argument of type IntegerType(32 bits) but got PointerType(Reference)\n",
			)
		})

		// Test 49: Accessing field of non-struct – compile error
		t.Run("Test 49: Field access on non-struct", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/deepseek_test49_field_access_on_non_struct.english.swa",
				"variable a is not a struct instance\n",
			)
		})
	})
}
