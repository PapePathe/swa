package lexer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWithMalinke(t *testing.T) {
	lex, err := New("dialect:malinke;")
	assert.NoError(t, err)

	wf := Malinke{}
	assert.Equal(t, len(wf.Patterns()), len(lex.Patterns()))
}

func TestNewWithWolof(t *testing.T) {
	lex, err := New("dialect:wolof;")
	assert.NoError(t, err)

	wf := Wolof{}
	assert.Equal(t, len(wf.Patterns()), len(lex.Patterns()))
}

func TestNewWithFrench(t *testing.T) {
	lex, err := New("dialect:french;")
	assert.NoError(t, err)

	wf := French{}
	assert.Equal(t, len(wf.Patterns()), len(lex.Patterns()))
}

func TestNewWithEmptyString(t *testing.T) {
	_, err := New("")
	assert.Equal(t, errors.New("You must define your dialect"), err)
}

func TestNewWithUnknownDialect(t *testing.T) {
	_, err := New("dialect:japanese;")
	assert.Equal(t, errors.New("dialect <japanese> is not supported"), err)
}
