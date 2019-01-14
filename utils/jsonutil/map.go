package jsonutil

import (
	"fmt"
)

type VariableMap interface {
	Variable

	Get(s string) (Variable, bool)
	Set(s string, v Variable)
	Delete(s string)

	Copy() VariableMap
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

func (m *varMap) Get(s string) (Variable, bool) {
	v, ok := m.vars[s]
	if !ok {
		return nil, false
	}
	return v, true
}

func (m *varMap) Set(s string, v Variable) {
	m.vars[s] = v
}

func (m *varMap) Len() int {
	return len(m.vars)
}

func (m *varMap) Delete(s string) {
	delete(m.vars, s)
}

func (m *varMap) Copy() VariableMap {
	vs := make(map[string]Variable)
	for k, v := range m.vars {
		vs[k] = v
	}
	return NewVariableMap(m.name, vs)

}

func (m *varMap) to() map[string]Variable {
	return m.vars
}

func NewVariableMap(name string, vs map[string]Variable) VariableMap {
	if vs == nil {
		vs = map[string]Variable{}
	}
	return &varMap{
		name: name,
		vars: vs,
	}
}
