package scopes

import (
	"errors"
)

// Obvious implementation of Scope
type basicScope struct {
	names map[string]Object
}

// TODO: instead of name, we'll need to look up by Signature (which is a name with optional argument types)
//   Variable Signature: "x" (no type data, just a name)
//   Function Signature: "f(T...)" (name and argument types, no return type)
//   Function return types and clauses are not part of the signature
//     e.g. two separate functions cannot be distinguished by just their return types and clauses
// Separate name from the things that can be named. Object should be a signature including name and type.
// Things that can be named should have a way to return their type.
// When assigning, we should call value.Type().AssignableTo(signature.Type())
func (s *basicScope) Lookup(name string) (Object, error) {
	if n, ok := s.names[name]; ok {
		return n, nil
	}

	return nil, errors.New("name not found: " + name)
}

func (s *basicScope) DeclareValue(name string, value Value) error {
	if _, ok := s.names[name]; ok {
		return errors.New("name already declared in scope: " + name)
	}

	s.names[name] = value

	return nil
}

func (s *basicScope) DeclareFunction(name string, function Function) error {
	if _, ok := s.names[name]; ok {
		return errors.New("name already declared in scope: " + name)
	}

	s.names[name] = function

	return nil
}

func (s *basicScope) AssignValue(name string, value Value) error {
	if _, ok := s.names[name]; ok {
		s.names[name] = value

		return nil
	}

	return errors.New("name not found: " + name)
}

// Obvious implementation of a ScopeStack
type basicScopeStack struct {
	scopes []Scope
}

// Push creates a new basicScope and pushes it onto scopes
func (s *basicScopeStack) Push() {
	s.scopes = append(s.scopes, &basicScope{names: map[string]Object{}})
}

func (s *basicScopeStack) Pop() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *basicScopeStack) Lookup(name string) (Object, error) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if v, err := s.scopes[i].Lookup(name); err != nil {
			return v, nil
		}
	}

	return nil, errors.New("name not found: " + name)
}

func (s *basicScopeStack) DeclareValue(name string, value Value) error {
	return s.scopes[len(s.scopes)].DeclareValue(name, value)
}

func (s *basicScopeStack) DeclareFunction(name string, function Function) error {
	return s.scopes[len(s.scopes)].DeclareFunction(name, function)
}

func (s *basicScopeStack) AssignValue(name string, value Value) error {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if err := s.scopes[i].AssignValue(name, value); err != nil {
			return nil
		}
	}

	return errors.New("name not found: " + name)
}

func NewScopeStack() *basicScopeStack {
	return &basicScopeStack{}
}
