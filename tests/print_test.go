package tests

import (
	"testing"
)

func TestPrint(t *testing.T) {
	t.Run("static string", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./print/binexpr.swa",
			"2",
		)
	})

	t.Run("table display", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./print/table.swa",
			`   id|                dept| salary
   20|     human resources|1000000
   21|          accounting|1000000
   22|             finance|1000000
   23|         engineering|1000000
`,
		)
	})

	t.Run("French", func(t *testing.T) {
		t.Run("static string", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t,
				"./print/variable.string.french.swa",
				"contenu de la variable: french",
			)
		})

		t.Run("float expression", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t,
				"./print/expression.float.french.swa",
				"a: 10.005000, b: -10.005000",
			)
		})

		t.Run("float variable", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t,
				"./print/variable.float.french.swa",
				"a: 10.005000, b: -10.005000",
			)
		})

		t.Run("Number variable", func(t *testing.T) {
			NewSuccessfulCompileRequest(
				t,
				"./print/variable.number.french.swa",
				"contenu de la variable: 10",
			)
		})
	})

	t.Run("modern print", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./print/modern_print.swa",
			"Modern print: 42 3.140000 hello\n423.140000hello\n",
		)
	})

	t.Run("expanded print", func(t *testing.T) {
		NewSuccessfulCompileRequest(
			t,
			"./print/expanded_print.swa",
			"Struct member: 50\nNested member: 50\nArray of structs: 2\nArray access: 20\nFunction call: 100\nBinary expr: 30 1\nError type: runtime error\nPrefix expr: -10\nZero value: 0\n",
		)
	})
}
