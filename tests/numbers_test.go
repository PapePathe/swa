package tests

import (
	"testing"
)

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

		t.Run("Wolof", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.wolof.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.french.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/float.soussou.swa",
				"zero value: 0.000000 assigned value: 100.000000",
			)
		})
	})

	t.Run("Int", func(t *testing.T) {
		t.Run("English", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("Wolof", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.wolof.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("French", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.french.swa",
				"zero value: 0 assigned value: 100",
			)
		})

		t.Run("Soussou", func(t *testing.T) {
			NewSuccessfulCompileRequest(t,
				"./numbers/zero-values/int.soussou.swa",
				"zero value: 0 assigned value: 100",
			)
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
