package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestStructsEnglish(t *testing.T) {
//	t.Parallel()
//
//	output, err := CompileSwaCode(t, "./examples/struct.english.swa", "struct.english")
//
//	assert.NoError(t, err)
//	assert.Equal(t, string(output), "")
//
//	assertCodeGenerated(t, "struct.english")
//	assertFileContent(t, "./examples/struct.english.ll", "./conditional.english.ll")
//	cleanupSwaCode(t, "struct.english")
//}

func TestStructsFrench(t *testing.T) {
	t.Parallel()

	output, err := CompileSwaCode(t, "./examples/struct.french.swa", "struct.french")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertCodeGenerated(t, "struct.french")
	assertFileContent(t, "./examples/struct.french.ll", "./struct.french.ll")
	cleanupSwaCode(t, "struct.french")
}
