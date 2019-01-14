package jsonutil

import (
	"fmt"
	"strconv"
)

// VariableArray defines array variable
type VariableArray interface {
	// Variable is Variable interface
	Variable

	// Append append Variable into array
	Append(v Variable)
	// Len returns array length
	Len() int
	// Get gets Variable from array
	Get(i int) (Variable, bool)
	// Slice returns slice from array
	Slice(i, k int) (VariableArray, error)
}

type varArray struct {
	name string
	vars []Variable
}

// Name implements Variable interface
func (arr *varArray) Name() string {
	return arr.name
}

// String implements Variable interface
func (arr *varArray) String() string {
	return fmt.Sprint(arr.vars)
}

// Unmarshal implements Variable interface
func (arr *varArray) Unmarshal(obj interface{}) error {
	return fmt.Errorf("Not Supported")
}

func (arr *varArray) Select(selector ...string) (Variable, error) {
	if len(selector) == 0 {
		return arr, nil
	}
	index, err := strconv.Atoi(selector[0])
	if err != nil {
		return nil, fmt.Errorf("can't select from json(%v) with selector %v: %v", arr, selector, err)
	}

	v, ok := arr.Get(index)
	if !ok {
		return nil, fmt.Errorf("index out of bounds: expected %v, len %v", index, arr.Len())
	}
	return v.Select(selector[1:]...)
}

func (arr *varArray) Append(v Variable) {
	arr.vars = append(arr.vars, v)
}

func (arr *varArray) Len() int {
	return len(arr.vars)
}

func (arr *varArray) Get(i int) (Variable, bool) {
	if i >= arr.Len() {
		return nil, false
	}
	return arr.vars[i], true
}

func (arr *varArray) Slice(i, k int) (VariableArray, error) {
	if i < 0 || k < 0 {
		return nil, fmt.Errorf("invalid slice index (index must be non-negative)")
	}
	if i > k {
		return nil, fmt.Errorf("invalid slice index: %d > %d", i, k)
	}
	if k >= arr.Len() {
		return nil, fmt.Errorf("index out of bounds: %d > len(%d)", k, arr.Len())
	}
	arr.vars = arr.vars[i:k]
	return arr, nil
}

// NewVariableArray returns array variable
func NewVariableArray(name string, vs []Variable) VariableArray {
	return &varArray{
		name: name,
		vars: vs,
	}
}
