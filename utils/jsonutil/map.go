package jsonutil

import (
	"fmt"
)

// VariableMap defines variable map
type VariableMap interface {
	// Variable is Variable interface
	Variable

	// Get gets Variable from map
	Get(s string) (Variable, bool)
	// Set sets Variable into map
	Set(s string, v Variable)
	// Delete delete Variable from map
	Delete(s string)
	// Keys returns variable map keys
	Keys() []string

	// Copy copies and returns a new VariableMap
	Copy() VariableMap
	// to returns the internal variable map struct
	// it is only used by internal helper function
	to() map[string]Variable
}

type varMap struct {
	name string
	vars map[string]Variable
	keys []string
}

// Name implements Variable interface
func (m *varMap) Name() string {
	return m.name
}

// Type implements Variable interface
func (m *varMap) Type() JSONType {
	return ObjectType
}

// String implements Variable interface
func (m *varMap) String() string {
	return fmt.Sprint(m.vars)
}

// Unmarshal implements Variable interface
func (m *varMap) Unmarshal(obj interface{}) error {
	return fmt.Errorf("Not Supported")
}

// Select implements Variable interface
func (m *varMap) Select(selector ...string) (Variable, error) {
	if len(selector) == 0 {
		return m, nil
	}
	v, ok := m.Get(selector[0])
	if !ok {
		return nil, fmt.Errorf("can't select from json(%v) with selector %v", m, selector)
	}
	if v == nil {
		if len(selector) == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("can't select from json(null) with selector %v", selector[1:])
	}
	return v.Select(selector[1:]...)
}

// Get implements VariableMap interface
func (m *varMap) Get(s string) (Variable, bool) {
	v, ok := m.vars[s]
	if !ok {
		return nil, false
	}
	return v, true
}

// Set implements VariableMap interface
func (m *varMap) Set(s string, v Variable) {
	m.vars[s] = v
	m.keys = append(m.keys, s)
}

// Delete implements VariableMap interface
func (m *varMap) Delete(s string) {
	delete(m.vars, s)
	index := -1
	for i, k := range m.keys {
		if k == s {
			index = i
			break
		}
	}
	if index != -1 {
		m.keys = append(m.keys[:index], m.keys[index+1:]...)
	}
}

// Keys implements VariableMap interface
func (m *varMap) Keys() []string {
	return m.keys
}

// Copy implements VariableMap interface
func (m *varMap) Copy() VariableMap {
	vs := make(map[string]Variable)
	for k, v := range m.vars {
		vs[k] = v
	}
	return NewVariableMap(m.name, vs)

}

// to implements VariableMap interface
func (m *varMap) to() map[string]Variable {
	return m.vars
}

// NewVariableMap returns a new VariableMap
func NewVariableMap(name string, vs map[string]Variable) VariableMap {
	if vs == nil {
		vs = map[string]Variable{}
	}
	keys := []string{}
	for k := range vs {
		keys = append(keys, k)
	}
	return &varMap{
		name: name,
		vars: vs,
		keys: keys,
	}
}
