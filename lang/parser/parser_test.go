package parser

import (
	"testing"

	"swahili/lang/lexer"

	"github.com/stretchr/testify/assert"
)

func TestLSPResilience(t *testing.T) {
	// Test case 1: Missing semicolon
	source1 := `
dialect:english;
main {
	let x: int = 10
	print(x);
}
`
	tokens1, _, errs1 := lexer.Tokenize(source1)
	assert.Empty(t, errs1)

	ast1, errs1 := Parse(tokens1)
	assert.NotEmpty(t, errs1) // Should have error about missing semicolon
	assert.NotNil(t, ast1)

	// Test case 2: Invalid type
	source2 := `
dialect:english;
main {
	let x: invalid = 10;
}
`
	tokens2, _, errs2 := lexer.Tokenize(source2)
	assert.Empty(t, errs2)

	ast2, errs2 := Parse(tokens2)
	assert.NotEmpty(t, errs2) // Should have error about invalid type
	assert.NotNil(t, ast2)

	// Test case 3: Lexer error
	source3 := `
dialect:english;
main {
	let x = @;
}
`
	tokens3, _, errs3 := lexer.Tokenize(source3)
	assert.NotEmpty(t, errs3)   // Should have lexer error
	assert.NotEmpty(t, tokens3) // Should still have tokens
}
