package jsonutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

// Variable defines variable which can unmarshal to an object
type Variable interface {
	// Name returns variable name
	Name() string
	// Unmarshal will unmarshal the variable to an object
	Unmarshal(obj interface{}) error
	// String will returns the string value
	String() string
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

// String implements Variable interface
func (v *stringVar) String() string {
	return v.value
}

// Unmarshal implements Variable interface
func (v *stringVar) Unmarshal(obj interface{}) error {
	return fmt.Errorf("Not Supported")
}

// NewVariable returns a variable with value s
func NewVariable(name, value string) Variable {
	return &stringVar{
		name:  name,
		value: value,
	}
}
