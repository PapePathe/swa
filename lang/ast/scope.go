/*
* swahili/lang
* Copyright (C) 2025  Papa Pathe SENE
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

func (s *Scope) Exists(name string) bool {
	return false
}

func (s *Scope) Set(name string, val values.Value) {
	s.Variables[name] = val
}
