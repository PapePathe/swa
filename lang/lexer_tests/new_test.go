package lexertests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/lexer"
)

func TestNewWithMalinke(t *testing.T) {
	_, err := lexer.New("dialect:malinke;")
	assert.NoError(t, err, "New lexer should not error")
}

func TestNewWithWolof(t *testing.T) {
	_, err := lexer.New("dialect:wolof;")
	assert.NoError(t, err, "New lexer should not error")
}

func TestNewWithEmptyString(t *testing.T) {
	_, err := lexer.New("")
	assert.Equal(t, errors.New("You must define your dialect"), err)
}

func TestNewWithUnknownDialect(t *testing.T) {
	_, err := lexer.New("dialect:english;")
	assert.Equal(t, errors.New("dialect <english> is not supported"), err)
}
