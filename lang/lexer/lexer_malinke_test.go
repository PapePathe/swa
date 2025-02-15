package lexer

import (
	"os"
	"testing"

	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenizeMalinke ...
func TestTokenizeMalinke(t *testing.T) {
	sourcePath := "../examples/malinke/age_calculator.swa"
	destPath := "../test/fixtures/tokenizer/malinke/age_calculator.ast"

	bytes, err := os.ReadFile(sourcePath)
	require.NoError(t, err)

	source := string(bytes)
	tokens := Tokenize(source)
	result := litter.Sdump(tokens)
	expected, err := os.ReadFile(destPath)
	if err != nil {
		// First run, write test data since it doesn't exist
		if !os.IsNotExist(err) {
			t.Error(err)
		}
		err := os.WriteFile(destPath, []byte(result), 0644)
		assert.NoError(t, err)
		result = string(expected)
	}

	assert.Equal(t, result, string(expected))
}
