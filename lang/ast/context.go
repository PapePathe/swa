package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
)

type Var struct {
	def value.Value
	cst constant.Constant
}

type Context struct {
	*ir.Block
	mod    *ir.Module
	parent *Context
	vars   map[string]Var
}

func NewContext(b *ir.Block, mod *ir.Module) *Context {
	return &Context{
		Block:  b,
		mod:    mod,
		parent: nil,
		vars:   make(map[string]Var),
	}
}

func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b, c.mod)
	ctx.parent = c

	return ctx
}

func (c Context) LookupVariable(name string) Var {
	if v, ok := c.vars[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.LookupVariable(name)
	} else {
		fmt.Printf("variable: `%s`\n", name)
		panic("no such variable")
	}
}
