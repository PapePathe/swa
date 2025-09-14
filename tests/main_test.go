package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingFile(t *testing.T) {
	t.Parallel()

	_, err := CompileSwaCode(t, "./examples/not-found.swa", "not-found")

	assert.Equal(t, err.Error(), "exit status 1")
}

func assertFileContent(t *testing.T, actual string, expected string) {
	actualStr, err := os.ReadFile(actual)
	assert.NoError(t, err)

	assert.Equal(t, string(actualStr), expected)
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)

	assert.True(t, !os.IsNotExist(err))
}

func assertCodeGenerated(t *testing.T, output string) {
	t.Helper()

	assertFileExists(t, fmt.Sprintf("%s.ll", output))
	assertFileExists(t, fmt.Sprintf("%s.s", output))
	assertFileExists(t, fmt.Sprintf("%s.o", output))
	assertFileExists(t, fmt.Sprintf("%s.exe", output))
}

func cleanupSwaCode(t *testing.T, output string) {
	assert.NoError(t, os.Remove(fmt.Sprintf("%s.s", output)))
	assert.NoError(t, os.Remove(fmt.Sprintf("%s.o", output)))
	assert.NoError(t, os.Remove(fmt.Sprintf("%s.ll", output)))
	assert.NoError(t, os.Remove(fmt.Sprintf("%s.exe", output)))
}

func CompileSwaCode(t *testing.T, src string, dest string) ([]byte, error) {
	t.Helper()

	cmd := exec.Command("./swa", "compile", "-s", src, "-o", dest)
	return cmd.CombinedOutput()
}
