package tests

import (
	"strings"
	"testing"
)

func TestBinaryExpressionsStrings(t *testing.T) {
	t.Run("Strings in binary-expressions", func(t *testing.T) {
		tests := []struct {
			name     string
			filename string
		}{
			{"equality_literals", "equality_literals.english.swa"},
			{"inequality_literals", "inequality_literals.english.swa"},
			{"equality_variables", "equality_variables.english.swa"},
			{"inequality_variables", "inequality_variables.english.swa"},
			{"empty_strings", "empty_strings.english.swa"},
			{"empty_string_inequality", "empty_string_inequality.english.swa"},
			{"empty_literals", "empty_literals.english.swa"},
			{"case_sensitivity", "case_sensitivity.english.swa"},
			{"prefix_matching", "prefix_matching.english.swa"},
			{"spaces", "spaces.english.swa"},
			{"variable_literal_equality", "variable_literal_equality.english.swa"},
			{"variable_literal_inequality", "variable_literal_inequality.english.swa"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Run("English", func(t *testing.T) {
					NewSuccessfulCompileRequest(t, "./binary-expressions/strings/"+tt.filename, "Passed")
				})
				t.Run("Soussou", func(t *testing.T) {
					NewSuccessfulCompileRequest(
						t,
						"./binary-expressions/strings/"+strings.Replace(tt.filename, ".english.swa", ".soussou.swa", 1),
						"Passed",
					)
				})
			})
		}
	})
}

func TestBinaryExpressionsModulo(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Run("Modulo with symbol expressions", func(t *testing.T) {
			NewSuccessfulCompileRequest(t, "./binary-expressions/modulo/variable-declaration-2.english.swa", "v1: 0, v2: 1")
		})

		t.Run("Modulo in conditionals", func(t *testing.T) {
			NewSuccessfulCompileRequest(t, "./binary-expressions/modulo/conditionals.english.swa", "(4 modulo 3 is equal to 1)(4 modulo 2 is equal to 0)")
		})

		t.Run("Modulo with integer expressions", func(t *testing.T) {
			NewSuccessfulCompileRequest(t, "./binary-expressions/modulo/variable-declaration.english.swa", "x: 1, y: 0")
		})
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Run("Modulo with symbol expressions", func(t *testing.T) {
			NewSuccessfulCompileRequest(t, "./binary-expressions/modulo/variable-declaration-2.soussou.swa", "v1: 0, v2: 1")
		})

		t.Run("Modulo in conditionals", func(t *testing.T) {
			NewSuccessfulCompileRequest(t, "./binary-expressions/modulo/conditionals.soussou.swa", "(4 modulo 3 is equal to 1)(4 modulo 2 is equal to 0)")
		})

		t.Run("Modulo with integer expressions", func(t *testing.T) {
			NewSuccessfulCompileRequest(t, "./binary-expressions/modulo/variable-declaration.soussou.swa", "x: 1, y: 0")
		})
	})
}
