package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingFile(t *testing.T) {
	t.Parallel()

	_, err := compile(t, "./examples/not-found.swa", "not-found")

	assert.Equal(t, err.Error(), "exit status 1")
}

func TestMissingDialect(t *testing.T) {
	t.Parallel()

	output, err := compile(t, "./examples/missing_dialect.swa", "missing_dialect")

	assert.Error(t, err)
	assert.Equal(t, string(output), "you must define your dialect\n")
}

func TestConditionals(t *testing.T) {
	t.Parallel()

	output, err := compile(t, "./examples/conditional.swa", "conditional.swa")

	assert.NoError(t, err)
	assert.Equal(t, string(output), "")

	assertCodeGenerated(t)
	assertFileContent(t, "./examples/conditional.ll", "./start.ll")
}

func TestMissingSemiColon(t *testing.T) {
	t.Parallel()

	expected := "expected SEMI_COLON, but got CLOSE_CURLY  current: {} CLOSE_CURLY CLOSE_CURLY}\n"
	output, err := compile(
		t,
		"./examples/missing_semi_colon.swa",
		"missing_semi_colon.swa",
	)

	assert.Error(t, err)
	assert.Equal(t, string(output), expected)
}

func assertFileContent(t *testing.T, actual string, expected string) {
	actualStr, err := ioutil.ReadFile(actual)
	assert.NoError(t, err)

	expectedStr, err := ioutil.ReadFile(expected)
	assert.NoError(t, err)

	assert.Equal(t, string(actualStr), string(expectedStr))
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)

	assert.True(t, !os.IsNotExist(err))
}

func assertCodeGenerated(t *testing.T) {
	t.Helper()

	assertFileExists(t, "./start.ll")
	assertFileExists(t, "./start.s")
	assertFileExists(t, "./start.o")
	assertFileExists(t, "./start.exe")
}

func compile(t *testing.T, src string, dest string) ([]byte, error) {
	t.Helper()

	cmd := exec.Command("./swa", "compile", "-s", src)
	return cmd.CombinedOutput()
}
