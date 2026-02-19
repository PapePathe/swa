package tests

import (
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	t.Run("01", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/01.swa", "5 + 3 = 8")
	})
	t.Run("02", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/02.swa", "10 - 4 = 6")
	})
	t.Run("03", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/03.swa", "7 * 6 = 42")
	})
	t.Run("04", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/04.swa", "15 / 4 = 3")
	})
	t.Run("05", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/05.swa", "2 + 3 * 4 = 14")
	})
	t.Run("06", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/06.swa", "(2 + 3) * 4 = 20")
	})
	t.Run("07", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/07.swa", "10 - 3 - 2 = 5")
	})
	t.Run("08", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/08.swa", "16 / 4 / 2 = 2")
	})
	t.Run("09", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/09.swa", "5 + 3 * 2 - 8 / 4 = 9")
	})
	t.Run("10", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/10.swa", "-5 = -5")
	})
	t.Run("11", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/11.swa", "- 5 + 3 = -2")
	})
	t.Run("12", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/12.swa", "-(3 + 4) = -7")
	})
	t.Run("13", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/13.swa", "5 - (-3) = 8")
	})
	t.Run("14", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/14.swa", "-4 * 6 = -24")
	})
	t.Run("15", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/15.swa", "-15 / 4 = -3")
	})
	t.Run("16", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/16.swa", "-15 % 4 = -3")
	})
	t.Run("17", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/17.swa", "15 % -4 = 3")
	})
	t.Run("19", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/19.swa", "-15 % -4 = -3")
	})
	t.Run("21", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/21.swa", "-(-5) = 5")
	})
	t.Run("22", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/22.swa", "-x (x=7) = -7")
	})
	t.Run("23", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/23.swa", "-17 % 5 = -2")
	})
	t.Run("24", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/24.swa", "7 / 2 = 3")
	})
	t.Run("25", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/25.swa", "-7 / 2 = -3")
	})
	t.Run("26", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/26.swa", "7 / -2 = -3")
	})
	t.Run("27", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/27.swa", "-7 / -2 = 3")
	})

	t.Run("28", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/28.swa", "2 + 3*4 - (6/2) + -5%3 = 9")
	})
	t.Run("29", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/29.swa", "12 * 3 / 4 = 9")
	})
	t.Run("30", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/30.swa", "5 - 3 + 2 = 4")
	})
	t.Run("31", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/31.swa", "-(3 * 4) = -12")
	})
	t.Run("32", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/32.swa", "- -5 = 5")
	})
	t.Run("33", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/33.swa", "(-3) * (-4) = 12")
	})
	t.Run("34", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/34.swa", "3 % 7 = 3")
	})
	t.Run("35", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/35.swa", "-3 % 7 = -3")
	})
	t.Run("36", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/36.swa", "-2 * 3 + 4 / -2 - (-5) % 3 = -6")
	})
	t.Run("37", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/37.swa", "-(2 + 3 * 4 - 5) = -9")
	})
	t.Run("38", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/integer-arithmetic/38.swa", "((2+3)*(4-1))/(5-2) = 5")
	})
}

func TestMixedNumbers(t *testing.T) {
	t.Run("Mixed", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/mixed.english.swa",
			"(x + y = 7.50, y + x = 7.50)(x * y = 12.50, y * x = 12.50)(x / y = 2.00, y / x = 0.50)(x - y = 2.50, y - x = -2.50)",
		)
	})
}

func TestNumberZeroValues(t *testing.T) {
	t.Run("Float", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})

		t.Run("English.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/float.swa",
			})
		})

		t.Run("Wolof", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.wolof.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})

		t.Run("Wolof.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/float.wolof.swa",
			})
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.french.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})

		t.Run("French.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/float.french.swa",
			})
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.soussou.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})

		t.Run("Soussou.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/float.soussou.swa",
			})
		})
	})

	t.Run("Int", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("English.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/int.swa",
			})
		})

		t.Run("Wolof", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.wolof.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("Wolof.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/int.wolof.swa",
			})
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.french.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("French.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/int.french.swa",
			})
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.soussou.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("Soussou.X", func(t *testing.T) {
			NewSuccessfulXCompileRequest(CompileRequest{
				T:         t,
				InputPath: "./numbers/zero-values/int.soussou.swa",
			})
		})
	})
}

func TestSignedNumbers(t *testing.T) {
	t.Run("English", func(t *testing.T) {
		t.Run("Signed", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/signed.english.swa",
				"x+y: 0, x*y: -4, x-y: -4, x/y: -1, x%y: 0",
			)
		})

		t.Run("integer arithmetic", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/integer-arithmetic/source.english.swa",
				"okokokokokok",
			)
		})

		t.Run("integer arithmetic french", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/integer-arithmetic/source.french.swa",
				"okokokokokok",
			)
		})

		t.Run("integer arithmetic soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/integer-arithmetic/source.soussou.swa",
				"okokokokokok",
			)
		})

		t.Run("float-arithmetic", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/float-arithmetic/source.english.swa",
				"okokokokokok",
			)
		})

		t.Run("integer-arrays", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/integer-arrays/source.english.swa",
				"okokok",
			)
		})

		t.Run("float-arrays", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/float-arrays/source.english.swa",
				"okokok",
			)
		})

		t.Run("integer-float-separation", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/integer-float-separation/source.english.swa",
				"",
			)
		})
	})

	t.Run("French", func(t *testing.T) {
		NewSuccessfulCompileRequest(t,
			"./numbers/signed.french.swa",
			"x+y: 0, x*y: -4, x-y: -4, x/y: -1, x%y: 0",
		)
	})
}
