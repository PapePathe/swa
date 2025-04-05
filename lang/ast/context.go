package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
)

type GlobalVar struct {
	def value.Value
	cst constant.Constant
}

type LocalVariable struct {
	Value *ir.InstAlloca
}

type Context struct {
	*ir.Block
	mod        *ir.Module
	parent     *Context
	globalVars map[string]GlobalVar
	localVars  map[string]LocalVariable
}

func NewContext(b *ir.Block, mod *ir.Module) *Context {
	return &Context{
		Block:      b,
		mod:        mod,
		parent:     nil,
		globalVars: make(map[string]GlobalVar),
		localVars:  make(map[string]LocalVariable),
	}
}

func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b, c.mod)
	ctx.parent = c

	return ctx
}

func (c Context) AddLocal(name string, value LocalVariable) error {
	c.localVars[name] = value

	return nil
}

func (c Context) LookupLocalVariable(name string) (*LocalVariable, error) {
	if v, ok := c.localVars[name]; ok {
		return &v, nil
	} else if c.parent != nil {
		return c.parent.LookupLocalVariable(name)
	} else {
		return nil, fmt.Errorf("no such local variable: `%s`\n", name)
	}
}

func (c Context) LookupVariable(name string) (*GlobalVar, error) {
	if v, ok := c.globalVars[name]; ok {
		return &v, nil
	} else if c.parent != nil {
		return c.parent.LookupVariable(name)
	} else {
		return nil, fmt.Errorf("undefined variable: `%s`\n", name)
	}
}
