package tests

import (
	"strings"
	"testing"
)

func TestBinaryExpressionsModulo(t *testing.T) {

	t.Run("Strings in binary-expressions", func(t *testing.T) {
		t.Parallel()

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
					req := CompileRequest{
						InputPath:               "./binary-expressions/strings/" + tt.filename,
						T:                       t,
						ExpectedExecutionOutput: "Passed",
					}
					req.AssertCompileAndExecute()
				})
				t.Run("Soussou", func(t *testing.T) {
					req := CompileRequest{
						InputPath:               "./binary-expressions/strings/" + strings.Replace(tt.filename, ".english.swa", ".soussou.swa", 1),
						T:                       t,
						ExpectedExecutionOutput: "Passed",
					}
					req.AssertCompileAndExecute()
				})
			})
		}
	})

	t.Run("English", func(t *testing.T) {
		t.Parallel()

		t.Run("Modulo with symbol expressions", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/variable-declaration-2.english.swa",
				T:                       t,
				ExpectedExecutionOutput: "v1: 0, v2: 1",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Modulo in conditionals", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/conditionals.english.swa",
				T:                       t,
				ExpectedExecutionOutput: "(4 modulo 3 is equal to 1)(4 modulo 2 is equal to 0)",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Modulo with integer expressions", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/variable-declaration.english.swa",
				T:                       t,
				ExpectedExecutionOutput: "x: 1, y: 0",
			}

			req.AssertCompileAndExecute()
		})
	})

	t.Run("Soussou", func(t *testing.T) {
		t.Parallel()

		t.Run("Modulo with symbol expressions", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/variable-declaration-2.soussou.swa",
				T:                       t,
				ExpectedExecutionOutput: "v1: 0, v2: 1",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Modulo in conditionals", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/conditionals.soussou.swa",
				T:                       t,
				ExpectedExecutionOutput: "(4 modulo 3 is equal to 1)(4 modulo 2 is equal to 0)",
			}

			req.AssertCompileAndExecute()
		})

		t.Run("Modulo with integer expressions", func(t *testing.T) {
			req := CompileRequest{
				InputPath:               "./binary-expressions/modulo/variable-declaration.soussou.swa",
				T:                       t,
				ExpectedExecutionOutput: "x: 1, y: 0",
			}

			req.AssertCompileAndExecute()
		})
	})
}
