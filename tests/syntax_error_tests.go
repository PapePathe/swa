package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingDialect(t *testing.T) {
	t.Parallel()

	output, err := CompileSwaCode(t, "./examples/missing_dialect.swa", "missing_dialect")

	assert.Error(t, err)
	assert.Equal(t, string(output), "you must define your dialect\n")
}

func TestMissingSemiColon(t *testing.T) {
	t.Parallel()

	expected := "expected SEMI_COLON, but got CLOSE_CURLY  current: {} CLOSE_CURLY CLOSE_CURLY}\n"
	output, err := CompileSwaCode(
		t,
		"./examples/missing_semi_colon.swa",
		"missing_semi_colon.swa",
	)

	assert.Error(t, err)
	assert.Equal(t, string(output), expected)
}
