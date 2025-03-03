package ast

import "swahili/lang/values"

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Variables: make(map[string]values.Value),
		Parent:    parent,
	}
}

type Scope struct {
	Variables map[string]values.Value
	Parent    *Scope
}

func (s *Scope) Get(name string) (values.Value, bool) {
	val, ok := s.Variables[name]
	if !ok && s.Parent != nil {
		return s.Parent.Get(name)
	}

	return val, ok
}

func (s *Scope) Set(name string, val values.Value) {
	s.Variables[name] = val
}
