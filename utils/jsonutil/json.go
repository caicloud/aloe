package jsonutil

import (
	"fmt"

	"github.com/buger/jsonparser"
)

const (
	// LenSelector used to select length of array or object
	LenSelector = "#"
)

// GetVariable returns a variable from raw json
func GetVariable(rawJSON []byte, name string, selector ...string) (Variable, error) {
	vSelector, selectLen := selector, false

	if len(selector) != 0 && selector[len(selector)-1] == LenSelector {
		vSelector = selector[:len(selector)-1]
		selectLen = true
	}

	v, dt, _, err := jsonparser.Get(rawJSON, vSelector...)
	if err != nil {
		return nil, getVariableErrorf(name, string(rawJSON), selector, err)
	}
	t := convert(dt)
	switch t {
	case ArrayType:
		if !selectLen {
			break
		}
		l, err := countArray(v)
		if err != nil {
			return nil, getVariableErrorf(name, string(rawJSON), selector, err)
		}
		return NewIntVariable(name, int64(l)), nil
	case ObjectType:
		if !selectLen {
			break
		}
		l, err := countObject(v)
		if err != nil {
			return nil, getVariableErrorf(name, string(rawJSON), selector, err)
		}
		return NewIntVariable(name, int64(l)), nil
	case "":
		return nil, getVariableErrorf(name, string(rawJSON), selector, fmt.Errorf("unknown type"))
	default:
		if selectLen {
			return nil, getVariableErrorf(name, string(rawJSON), selector, fmt.Errorf("can't select len(#) for %v", t))
		}
	}
	return &variable{
		raw:      v,
		name:     name,
		jsonType: convert(dt),
	}, nil
}

func getVariableErrorf(name string, json string, selector []string, err error) error {
	return fmt.Errorf("can't get variable %s from json(%s) with selector %v: %v", name, json, selector, err)
}

func countArray(arrayJSON []byte) (int, error) {
	count := 0
	if _, err := jsonparser.ArrayEach(arrayJSON, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		count++
	}); err != nil {
		return 0, err
	}
	return count, nil
}

func countObject(objJSON []byte) (int, error) {
	count := 0
	if err := jsonparser.ObjectEach(objJSON, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
		count++
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}

func convert(dt jsonparser.ValueType) JSONType {
	switch dt {
	case jsonparser.NotExist, jsonparser.Unknown:
		return ""
	case jsonparser.String:
		return StringType
	case jsonparser.Number:
		return NumberType
	case jsonparser.Object:
		return ObjectType
	case jsonparser.Array:
		return ArrayType
	case jsonparser.Boolean:
		return BooleanType
	case jsonparser.Null:
		return NullType
	default:
		return ""
	}
}
