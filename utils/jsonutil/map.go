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

	// Copy copies and returns a new VariableMap
	Copy() VariableMap
	// to returns the internal variable map struct
	// it is only used by internal helper function
	to() map[string]Variable
}

type varMap struct {
	name string
	vars map[string]Variable
}

// Name implements Variable interface
func (m *varMap) Name() string {
	return m.name
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
}

// Delete implements VariableMap interface
func (m *varMap) Delete(s string) {
	delete(m.vars, s)
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
	return &varMap{
		name: name,
		vars: vs,
	}
}
