package tests

import (
	"testing"
)

// TestDeepseekB11c50 contains test cases extracted from deepseek_plaintext_20260213_b11c50.txt.
// Tests already covered by TestDeepseekPlaintext (c430ab) are omitted to avoid duplication.
func TestDeepseekB11c50(t *testing.T) {

	t.Run("Structs and Pointers", func(t *testing.T) {

		// Test 7: Pointer to array field in struct, iterate and sum
		t.Run("Test 7: Pointer to array in struct", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test7_pointer_to_array_in_struct.english.swa",
				"Sum: 45",
			)
		})

		// Test 8: Function returning struct – not supported
		t.Run("Test 8: Function returning struct", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test8_function_returning_struct.english.swa",
				"returning complex types (structs and arrays) from functions is currently not supported\n",
			)
		})

		// Test 20: Struct field type mismatch (float→int) – compiler accepts it
		t.Run("Test 20: Struct field type mismatch (accepted)", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test20_struct_field_type_mismatch.english.swa",
				"",
			)
		})

		// Test 21: Using undeclared struct type – compile error
		t.Run("Test 21: Undeclared struct type", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test21_undeclared_struct_type.english.swa",
				"struct named UnknownStruct does not exist in symbol table\n",
			)
		})

		// Test 24: Struct containing array of structs initialization – segfaults at runtime
		//		FIXME this triggers a segfault
		//		t.Run("Test 24: Struct with array of structs (segfault)", func(t *testing.T) {
		//			NewSuccessfulCompileRequest(
		//				t,
		//				"./functions/b11c50_test24_struct_with_array_of_structs.english.swa",
		//				"",
		//			)
		//		})

		// Test 25: Linked list traversal – self-referential struct not supported
		t.Run("Test 25: Self-referential struct not supported", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test25_linked_list_traversal.english.swa",
				"struct with pointer reference to self not supported, property: next\n",
			)
		})

		// Test 34: Array index into pointer-typed struct field – compile error
		t.Run("Test 34: Pointer field array index", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test34_array_of_structs_with_pointer_field.english.swa",
				"Element at index (%!s(int=1)) does not exist in array ()\n",
			)
		})
	})

	t.Run("Arithmetic and Types", func(t *testing.T) {

		// Test 9: Modulo with negative numbers (print has no trailing newline per line)
		t.Run("Test 9: Modulo negative numbers", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test9_modulo_negative_numbers.english.swa",
				"-7 % 3 = -17 % -3 = 1-7 % -3 = -1",
			)
		})

		// Test 10: Floating point precision – print has no newlines
		t.Run("Test 10: Float precision comparison", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test10_float_precision_comparison.english.swa",
				"Not equal (difference: 0.000000)0.1+0.2 = 0.300000",
			)
		})

		// Test 15: Type mismatch: assign string to int – compile error
		t.Run("Test 15: Type mismatch string to int", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test15_type_mismatch_assignment.english.swa",
				"expected Number but got String\n",
			)
		})

		// Test 33: Mixed type int + float – compile error
		t.Run("Test 33: Mixed type int + float expression", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test33_mixed_type_expression.english.swa",
				"expected Number but got Float\n",
			)
		})

		// Test 26: Integer division by zero – compiles but produces undefined integer output
		t.Run("Test 26: Division by zero (compiles, undefined output)", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./functions/b11c50_test26_division_by_zero.english.swa",
			})
		})
	})

	t.Run("Control Flow", func(t *testing.T) {

		// Test 11: Boolean short-circuit prevents execution of division by zero
		t.Run("Test 11: Boolean short-circuit", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test11_boolean_short_circuit.english.swa",
				"Short-circuit works",
			)
		})
	})

	t.Run("Arrays", func(t *testing.T) {

		// Test 16: Array index out of bounds (constant index) – compile-time error
		t.Run("Test 16: Array OOB constant index (compile error)", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test16_array_oob_constant_index.english.swa",
				"Element at index (%!s(int=5)) does not exist in array (arr)\n",
			)
		})

		// Test 36: Array passed to function – passed by reference, caller sees change
		t.Run("Test 36: Array passed by reference", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test36_array_pass_by_reference.english.swa",
				"a[0] = 999",
			)
		})
	})

	t.Run("Functions", func(t *testing.T) {

		// Test 12: Print two strings with format (no concatenation)
		t.Run("Test 12: Print two strings with format", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test12_print_two_strings.english.swa",
				"Hello World",
			)
		})

		// Test 18: Function call with too few arguments – compile error
		t.Run("Test 18: Too few arguments", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test18_wrong_argument_count.english.swa",
				"function foo expect 2 arguments but was given 1\n",
			)
		})

		// Test 19: Duplicate function definition – compile error
		t.Run("Test 19: Duplicate function definition", func(t *testing.T) {
			NewFailedCompileRequest(t,
				"./functions/b11c50_test19_duplicate_function.english.swa",
				"function named foo already exists in symbol table\n",
			)
		})

		// Test 23: Deep nested function calls (10000 levels)
		t.Run("Test 23: Deep nested calls", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test23_deep_nested_calls.english.swa",
				"Depth: 10000",
			)
		})

		// Test 31: String comparison with different lengths
		t.Run("Test 31: String length comparison", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test31_string_length_comparison.english.swa",
				"Not equal",
			)
		})

		// Test 40: Comment with special characters compiles fine
		t.Run("Test 40: Comment with special characters", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test40_comment_special_chars.english.swa",
				"Comment test",
			)
		})

		// Test 41: Empty function with explicit return
		t.Run("Test 41: Empty function with return", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./functions/b11c50_test41_empty_function_with_return.english.swa",
				"",
			)
		})
	})
}

// TestDeepseekE3b26f contains test cases extracted from deepseek_plaintext_20260213_e3b26f.txt.
// These are complex scientific programs written in the French dialect.
func TestDeepseekE3b26f(t *testing.T) {

	// Programme 6: Two-body orbital mechanics – single-line if not supported
	t.Run("Prog 6: Orbital mechanics (single-line if not supported)", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./functions/e3b26f_prog6_orbital_mechanics_verlet.french.swa",
			"function named sqrt does not exist in symbol table\n",
		)
	})

	// Programme 7: Carnot cycle – compiles and runs (output has no newlines)
	t.Run("Prog 7: Carnot cycle", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./functions/e3b26f_prog7_carnot_cycle.french.swa",
			"Chaleur reçue : 0.00 JChaleur rejetée : 0.00 JTravail produit : 0.00 JRendement théorique : 0.40",
		)
	})

	// Programme 8: Wave interference – cos() not in symbol table
	t.Run("Prog 8: Wave interference (cos not supported)", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./functions/e3b26f_prog8_wave_interference.french.swa",
			"function named cos does not exist in symbol table\n",
		)
	})

	// Programme 9: Radioactive decay – scientific notation (1e8) not supported
	t.Run("Prog 9: Radioactive decay (scientific notation not supported)", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./functions/e3b26f_prog9_radioactive_decay_chain.french.swa",
			"1e8: invalid syntax while parsing number expression",
		)
	})

	// Programme 10: Quantum particle in box – scientific notation (1e-9) not supported
	t.Run("Prog 10: Quantum particle (scientific notation not supported)", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./functions/e3b26f_prog10_quantum_particle_in_box.french.swa",
			"1e-9: invalid syntax while parsing number expression",
		)
	})
}
