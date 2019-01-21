package jsonutil

import (
	"fmt"

	"github.com/buger/jsonparser"
)

// GetVariable returns a variable from raw json
func GetVariable(rawJSON []byte, name string, selector ...string) (Variable, error) {
	v, dt, _, err := jsonparser.Get(rawJSON, selector...)
	if err != nil {
		return nil, fmt.Errorf("can't get variable %v from json(%v) with selector %v: %v", name, string(rawJSON), selector, err)
	}
	t := convert(dt)
	if t == "" {
		return nil, fmt.Errorf("can't get variable %v from json(%v) with selector %v: unknown type", name, string(rawJSON), selector)
	}
	return &variable{
		raw:      v,
		name:     name,
		jsonType: convert(dt),
	}, nil
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
