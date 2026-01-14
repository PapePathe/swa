package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsuppotedFeatures(t *testing.T) {
	t.Parallel()

	t.Run("Self refencing pointer to struct", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./unsupported-features/self-referencing-pointer-to-struct.swa",
			ExpectedOutput: "struct with pointer reference to self not supported, property: parent\n",
			T:              t,
		}

		assert.Error(t, req.Compile())

	})

	t.Run("Self refencing struct", func(t *testing.T) {
		req := CompileRequest{
			InputPath:      "./unsupported-features/self-referencing-struct.swa",
			ExpectedOutput: "struct with reference to self not supported, property: parent\n",
			T:              t,
		}

		assert.Error(t, req.Compile())

	})
}
