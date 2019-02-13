package jsonutil

import (
	"fmt"
	"strconv"
)

// VariableArray defines array variable
type VariableArray interface {
	// Variable is Variable interface
	Variable

	// Measurable defines Len function
	Measurable

	// Append append Variable into array
	Append(v Variable)
	// Get gets Variable from array
	Get(i int) Variable
	// Set sets Variable into array
	Set(i int, v Variable)
	// Slice returns slice from array
	Slice(i, k int) VariableArray

	// to returns the internal variable array
	// it is only used by internal helper function
	to() []Variable
}

type varArray struct {
	name string
	vars []Variable
}

// Name implements Variable interface
func (arr *varArray) Name() string {
	return arr.name
}

// Type implements Variable interface
func (arr *varArray) Type() JSONType {
	return ArrayType
}

// String implements Variable interface
func (arr *varArray) String() string {
	return fmt.Sprint(arr.vars)
}

// Unmarshal implements Variable interface
func (arr *varArray) Unmarshal(obj interface{}) error {
	return fmt.Errorf("Not Supported")
}

// Select implements Variable interface
func (arr *varArray) Select(selector ...string) (Variable, error) {
	if len(selector) == 0 {
		return arr, nil
	}

	// 3 is length of [0]
	if len(selector[0]) < 3 {
		return nil, fmt.Errorf("can't select from json(%v) with selector %v: missing []", arr, selector)
	}

	index, err := strconv.Atoi(selector[0][1 : len(selector[0])-1])
	if err != nil {
		return nil, fmt.Errorf("can't select from json(%v) with selector %v: %v", arr, selector, err)
	}

	if index >= arr.Len() {
		return nil, fmt.Errorf("index out of bounds: expected %v, len %v", index, arr.Len())
	}
	v := arr.Get(index)
	if v == nil {
		if len(selector) == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("can't select from json(%v) with selector %v", v, selector[1:])
	}
	return v.Select(selector[1:]...)
}

func (arr *varArray) Append(v Variable) {
	arr.vars = append(arr.vars, v)
}

func (arr *varArray) Len() int {
	return len(arr.vars)
}

func (arr *varArray) Get(i int) Variable {
	return arr.vars[i]
}

func (arr *varArray) Set(i int, v Variable) {
	arr.vars[i] = v
}

func (arr *varArray) Slice(i, k int) VariableArray {
	arr.vars = arr.vars[i:k]
	return arr
}

// to implements VariableArray interface
func (arr *varArray) to() []Variable {
	return arr.vars
}

// NewVariableArray returns array variable
func NewVariableArray(name string, vs []Variable) VariableArray {
	return &varArray{
		name: name,
		vars: vs,
	}
}
