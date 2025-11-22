package lexer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWithMalinke(t *testing.T) {
	lex, dial, err := New("dialect:malinke;")
	assert.NoError(t, err)

	wf := Malinke{}
	assert.Equal(t, wf, dial)
	assert.Equal(t, len(wf.Patterns()), len(lex.Patterns()))
	assert.Equal(t, len(wf.Reserved()), len(lex.reservedWords))
}

func TestNewWithWolof(t *testing.T) {
	lex, dial, err := New("dialect:wolof;")
	assert.NoError(t, err)

	wf := Wolof{}

	assert.Equal(t, wf, dial)
	if len(wf.Patterns()) != len(lex.patterns) {
		t.Errorf("Expected patterns count to be %d", len(wf.Patterns()))
	}

	if len(wf.Reserved()) != len(lex.reservedWords) {
		t.Errorf("Expected reserved words count to be %d", len(wf.Patterns()))
	}

	assert.Equal(t, len(wf.Reserved()), len(lex.reservedWords))

	for key, val := range wf.Reserved() {
		if val != lex.reservedWords[key] {
			t.Errorf("Expected reserved word %s to be %s", key, val)
		}
	}
}

func TestNewWithFrench(t *testing.T) {
	lex, dial, err := New("dialecte:français;")
	assert.NoError(t, err)

	wf := French{}
	assert.Equal(t, len(wf.Patterns()), len(lex.Patterns()))
	assert.Equal(t, len(wf.Reserved()), len(lex.reservedWords))
	assert.Equal(t, wf, dial)
}

func TestNewWithEmptyString(t *testing.T) {
	_, _, err := New("")
	assert.Equal(t, errors.New("you must define your dialect"), err)
}

func TestNewWithUnknownDialect(t *testing.T) {
	_, _, err := New("dialect:japanese;")
	assert.Equal(t, errors.New("you must define your dialect"), err)
}
