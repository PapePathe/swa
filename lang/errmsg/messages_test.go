package errmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstError_Error(t *testing.T) {
	err := NewAstError("%s %d", "hello", 1)
	assert.Equal(t, "hello 1", err.Error())

	err = NewAstError("hello %s", "pathe")
	assert.Equal(t, "hello pathe", err.Error())
}
