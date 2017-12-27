package jsonutil

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/caicloud/aloe/template"
	"github.com/caicloud/aloe/types"
)

// GetVariable returns a variable from raw json
func GetVariable(rawJSON []byte, def *types.Definition) (*template.Variable, error) {
	v, dt, _, err := jsonparser.Get(rawJSON, def.Selector...)
	if err != nil {
		return nil, fmt.Errorf("can't get variable %v from json with selector %v: %v", def.Name, def.Selector, err)
	}
	t := convert(dt)
	if t == "" {
		return nil, fmt.Errorf("can't get variable %v from json with selector %v: unknown type", def.Name, def.Selector)
	}
	return &template.Variable{
		Raw:  v,
		Name: def.Name,
		Type: convert(dt),
	}, nil
}

func convert(dt jsonparser.ValueType) template.JSONType {
	switch dt {
	case jsonparser.NotExist, jsonparser.Unknown:
		return ""
	case jsonparser.String:
		return template.StringType
	case jsonparser.Number:
		return template.NumberType
	case jsonparser.Object:
		return template.ObjectType
	case jsonparser.Array:
		return template.ArrayType
	case jsonparser.Boolean:
		return template.BooleanType
	case jsonparser.Null:
		return template.NullType
	default:
		return ""

	}
}
