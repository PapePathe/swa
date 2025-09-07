package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionalsEnglish(t *testing.T) {
	t.Parallel()

	output, err := CompileSwaCode(t, "./examples/conditional.swa", "conditional.english")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertCodeGenerated(t, "conditional.english")
	assertFileContent(t, "./examples/conditional.english.ll", "./conditional.english.ll")
	cleanupSwaCode(t, "conditional.english")
}

func TestConditionalsFrench(t *testing.T) {
	t.Parallel()

	output, err := CompileSwaCode(t, "./examples/conditional.french.swa", "conditional.french")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertCodeGenerated(t, "conditional.french")
	assertFileContent(t, "./examples/conditional.french.ll", "./conditional.french.ll")
	cleanupSwaCode(t, "conditional.french")
}
