package template

import (
	"errors"
	"fmt"
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

// Variable is a snippet of json string
// object => {}
// array  => []
// string => "1bb"
// number => 1.3
// null   => null
type Variable struct {
	Raw  []byte
	Name string
	Type JSONType
}

// String returns variable value
func (v *Variable) String() string {
	switch v.Type {
	case StringType:
		return string(v.Raw)
	}
	return string(v.Raw)
}

// Template is a simple template support variable
// Golang template is too complex to use in this case
type Template interface {
	Render(vs map[string]Variable) (string, error)
}

// Template defines template of request
type template struct {
	varNames []string
	snippts  []string
}

func (t *template) fromRaw(raw string) error {
	t.varNames = nil
	t.snippts = nil
	snippt, varName := "", ""
	isVariable, isOpen := false, false

	for _, r := range raw {
		switch r {
		case '%':
			if isOpen || isVariable {
				snippt += "%"
			}
			isVariable = !isVariable
		case '{':
			if !isVariable {
				snippt += "{"
			} else {
				isVariable = false
				isOpen = true
				t.snippts = append(t.snippts, snippt)
				snippt = ""
			}
		case '}':
			if isOpen {
				isOpen = false
				if varName == "" {
					return errors.New("Param name should not be empty")
				}
				t.varNames = append(t.varNames, varName)
				varName = ""
			} else {
				snippt += "}"
			}
		default:
			if isVariable {
				return errors.New("Only %% or %{} is allowed")
			}
			if isOpen {
				varName += string(r)
			} else {
				snippt += string(r)
			}
		}
	}
	if isVariable || isOpen {
		return errors.New("Single '%' or unclosed '%{'")
	}
	t.snippts = append(t.snippts, snippt)
	return nil
}

// New returns raw string to template
func New(raw string) (Template, error) {
	t := template{}
	if err := t.fromRaw(raw); err != nil {
		return nil, err
	}
	return &t, nil
}

// Render renders template by variables and returns result
// Examples of rendering:
// Variables: {
//   "string": {
//     "raw": "xxx",
//     "type": "string",
//   },
//   "number": {
//     "raw": "1.5"
//     "type": "array",
//   }
// }
// template => rendered:
// %{string} => xxx
// "%{string}" => "xxx"
// %{number} => 1.5
// "%{number}" => "1.5"
// %% => %
// %%{string} => %{string}
func (t *template) Render(vs map[string]Variable) (string, error) {
	out := ""
	for i, varName := range t.varNames {
		out += t.snippts[i]
		v, ok := vs[varName]
		if !ok {
			return "", fmt.Errorf("can't find varibale %v", varName)
		}
		out += v.String()
	}
	out += t.snippts[len(t.snippts)-1]
	return out, nil
}
