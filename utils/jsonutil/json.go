package jsonutil

import (
	"fmt"

	"github.com/buger/jsonparser"
)

// GetVariable returns a variable from raw json
func GetVariable(rawJSON []byte, name string, selector ...string) (Variable, error) {
	v, dt, _, err := jsonparser.Get(rawJSON, selector...)
	if err != nil {
		return nil, getVariableErrorf(name, string(rawJSON), selector, err)
	}
	t := convert(dt)
	if t == "" {
		return nil, getVariableErrorf(name, string(rawJSON), selector, fmt.Errorf("unknown type"))
	}
	return &variable{
		raw:      v,
		name:     name,
		jsonType: t,
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
