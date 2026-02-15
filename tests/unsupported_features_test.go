package tests

import (
	"testing"
)

func TestUnsuppotedFeatures(t *testing.T) {
	t.Run("Self referencing pointer to struct", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./unsupported-features/self-referencing-pointer-to-struct.swa",
			"struct with pointer reference to self not supported, property: parent\n")
	})

	t.Run("Self referencing struct", func(t *testing.T) {
		NewFailedCompileRequest(t,
			"./unsupported-features/self-referencing-struct.swa",
			"struct with reference to self not supported, property: parent\n")
	})
}
