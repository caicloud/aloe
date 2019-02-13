package jsonutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// JSONType defines type of JSON
type JSONType string

const (
	// ObjectType is json snippet warpped by {}
	ObjectType JSONType = "object"

	// ArrayType is json snippet warpped by []
	ArrayType JSONType = "array"

	// StringType is json snippet of string
	// Value of variable with this type is not warpped by quote
	StringType JSONType = "string"

	// NumberType is json snippet which is number
	NumberType JSONType = "number"

	// BooleanType is json snippet which is bool
	BooleanType JSONType = "boolean"

	// NullType is json snippet means null
	NullType JSONType = "null"
)

// Measurable defines interface to return object length
type Measurable interface {
	// Len returns variable length
	// if variable can't be measured, return -1
	Len() int
}

// Variable defines variable which can unmarshal to an object
type Variable interface {
	// Name returns variable name
	Name() string
	// Unmarshal will unmarshal the variable to an object
	Unmarshal(obj interface{}) error
	// String will returns the string value
	String() string
	// Select get subpath of variable
	Select(selector ...string) (Variable, error)

	// Type returns variable type
	Type() JSONType
}

// variable is a snippet of json string
// object => {}
// array  => []
// string => "1bb"
// number => 1.3
// null   => null
type variable struct {
	raw      []byte
	name     string
	jsonType JSONType
}

// Name returns variable name
func (v *variable) Name() string {
	return v.name
}

// Type returns variable type
func (v *variable) Type() JSONType {
	return v.jsonType
}

// String returns variable value
func (v *variable) String() string {
	switch v.jsonType {
	case StringType:
		return string(v.raw)
	}
	return string(v.raw)
}

func (v *variable) Unmarshal(obj interface{}) error {
	switch v.jsonType {
	case ObjectType, ArrayType:
		return json.Unmarshal(v.raw, obj)
	case StringType:
		elem := reflect.ValueOf(obj).Elem()
		if elem.Kind() != reflect.String {
			return fmt.Errorf("can't unmarshal string to %T type", obj)
		}
		elem.SetString(string(v.raw[1 : len(v.raw)-1]))
	case BooleanType:
		elem := reflect.ValueOf(obj).Elem()
		k := elem.Kind()
		switch k {
		case reflect.Bool:
			b, err := strconv.ParseBool(string(v.raw))
			if err != nil {
				return err
			}
			elem.SetBool(b)
		case reflect.String:
			elem.SetString(string(v.raw))
		}
	case NumberType:
		elem := reflect.ValueOf(obj).Elem()
		k := elem.Kind()
		switch k {
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(string(v.raw), bitSizes[k])
			if err != nil {
				return err
			}
			elem.SetFloat(f)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(string(v.raw), 10, bitSizes[k])
			if err != nil {
				return err
			}
			elem.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			u, err := strconv.ParseUint(string(v.raw), 10, bitSizes[k])
			if err != nil {
				return err
			}
			elem.SetUint(u)
		case reflect.String:
			elem.SetString(string(v.raw))
		default:
			return fmt.Errorf("can't convert number to %T", obj)
		}
	case NullType:
		elem := reflect.ValueOf(obj).Elem()
		elem.Set(reflect.Zero(elem.Type()))
	}
	return nil
}

func (v *variable) Select(selector ...string) (Variable, error) {
	if len(selector) == 0 {
		return v, nil
	}
	name := strings.Join(selector, ".")
	return GetVariable(v.raw, name, selector...)
}

// Len implements Measurable interface
func (v *variable) Len() int {
	switch v.jsonType {
	case ArrayType:
		l, err := countArray(v.raw)
		if err != nil {
			return -1
		}
		return l
	case ObjectType:
		l, err := countObject(v.raw)
		if err != nil {
			return -1
		}
		return l
	}
	return -1
}

var bitSizes = map[reflect.Kind]int{
	reflect.Float32: 32,
	reflect.Float64: 64,

	reflect.Int:   0,
	reflect.Int8:  8,
	reflect.Int16: 16,
	reflect.Int32: 32,
	reflect.Int64: 64,

	reflect.Uint:   0,
	reflect.Uint8:  8,
	reflect.Uint16: 16,
	reflect.Uint32: 32,
	reflect.Uint64: 64,

	reflect.Uintptr: 64,
}

type stringVar struct {
	name  string
	value string
}

// Name implements Variable interface
func (v *stringVar) Name() string {
	return v.name
}

// Type implements Variable interface
func (v *stringVar) Type() JSONType {
	return StringType
}

// String implements Variable interface
func (v *stringVar) String() string {
	return v.value
}

// Unmarshal implements Variable interface
func (v *stringVar) Unmarshal(obj interface{}) error {
	value := reflect.ValueOf(obj).Elem()
	switch value.Kind() {
	case reflect.String:
		value.SetString(v.value)
	default:
		return fmt.Errorf("can't unmarshal string to object with type %T", obj)
	}
	return nil
}

func (v *stringVar) Select(selector ...string) (Variable, error) {
	if len(selector) == 0 {
		return v, nil
	}
	return nil, fmt.Errorf("can't select from json(%v) by selector %v", v, selector)
}

// NewStringVariable returns a variable with value s
func NewStringVariable(name, value string) Variable {
	return &stringVar{
		name:  name,
		value: value,
	}
}

type intVar struct {
	name  string
	value int64
}

// Name implements Variable interface
func (v *intVar) Name() string {
	return v.name
}

// Type implements Variable interface
func (v *intVar) Type() JSONType {
	return NumberType
}

// String implements Variable interface
func (v *intVar) String() string {
	return strconv.FormatInt(v.value, 10)
}

// Unmarshal implements Variable interface
func (v *intVar) Unmarshal(obj interface{}) error {
	value := reflect.ValueOf(obj).Elem()
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(v.value)
	default:
		return fmt.Errorf("can't unmarshal int to object with type %T", obj)
	}
	return nil
}

// Select implements Variable interface
func (v *intVar) Select(selector ...string) (Variable, error) {
	if len(selector) == 0 {
		return v, nil
	}
	return nil, fmt.Errorf("can't select from json(%v) by selector %v", v, selector)
}

// NewIntVariable returns a variable with value s
func NewIntVariable(name string, value int64) Variable {
	return &intVar{
		name:  name,
		value: value,
	}
}
