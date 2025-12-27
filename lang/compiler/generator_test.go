package compiler

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
	"tinygo.org/x/go-llvm"
)

func TestNumberExpression(t *testing.T) {
	context := llvm.GlobalContext()
	module := context.NewModule("swa-main")
	builder := context.NewBuilder()
	ctx := ast.NewCompilerContext(
		&context,
		&builder,
		&module,
		lexer.English{},
		nil,
	)

	t.Run("French", func(t *testing.T) {
		ctx.Dialect = lexer.French{}
		g := NewLLVMGenerator(ctx)
		t.Run("unsigned number", func(t *testing.T) {
			t.Parallel()

			err := g.VisitNumberExpression(&ast.NumberExpression{Value: 1})
			assert.NoError(t, err)
		})

		t.Run("signed number", func(t *testing.T) {
			t.Parallel()

			err := g.VisitNumberExpression(&ast.NumberExpression{Value: -1})
			assert.NoError(t, err)
		})

		t.Run("number too small", func(t *testing.T) {
			err := g.VisitNumberExpression(&ast.NumberExpression{Value: -19999999999999})
			assert.Error(t, err)
			assert.Equal(t, "-19999999999999 est plus petit que la valeur minimale pour un entier de 32 bits", err.Error())
		})

		t.Run("number too large", func(t *testing.T) {
			err := g.VisitNumberExpression(&ast.NumberExpression{Value: 19999999999999})
			assert.Error(t, err)
			assert.Equal(t, "19999999999999 est plus grand que la valeur maximale pour un entier de 32 bits", err.Error())
		})
	})
	t.Run("English", func(t *testing.T) {
		ctx.Dialect = lexer.English{}
		g := NewLLVMGenerator(ctx)
		t.Run("unsigned number", func(t *testing.T) {
			t.Parallel()

			err := g.VisitNumberExpression(&ast.NumberExpression{Value: 1})
			assert.NoError(t, err)
		})
		t.Run("signed number", func(t *testing.T) {
			t.Parallel()

			err := g.VisitNumberExpression(&ast.NumberExpression{Value: -1})
			assert.NoError(t, err)
		})

		t.Run("number too small", func(t *testing.T) {
			err := g.VisitNumberExpression(&ast.NumberExpression{Value: -19999999999999})
			assert.Error(t, err)
			assert.Equal(t, "-19999999999999 is smaller than min value for int32", err.Error())
		})

		t.Run("number too large", func(t *testing.T) {
			err := g.VisitNumberExpression(&ast.NumberExpression{Value: 19999999999999})
			assert.Error(t, err)
			assert.Equal(t, "19999999999999 is greater than max value for int32", err.Error())
		})
	})
}
