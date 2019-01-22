package template

import (
	"fmt"
	"strings"

	"github.com/caicloud/aloe/utils/jsonutil"
)

// Template is a simple template support variable
// Golang template is too complex to use in this case
type Template interface {
	Render(vs jsonutil.VariableMap) (string, error)
}

// Template defines template of request
type template struct {
	identitors map[int]identitor
	snippets   []string
	args       map[int][]identitor
}

type identitor struct {
	name  string
	isVar bool
}

func (t *template) fromRaw(raw string) error {
	lexer := NewLexer([]rune(raw))

	var args []identitor
	var funcName string
	for !lexer.IsEnd() {
		token, tokenType, err := lexer.NextToken()
		if err != nil {
			return err
		}
		switch tokenType {
		case TextToken:
			if funcName != "" {
				t.identitors[len(t.snippets)] = identitor{
					name:  funcName,
					isVar: false,
				}
				t.args[len(t.snippets)] = args
				funcName = ""
				args = nil
			}
			t.snippets = append(t.snippets, string(token))
		case VariableNameToken:
			t.identitors[len(t.snippets)] = identitor{
				name:  string(token),
				isVar: true,
			}
		case FuncNameToken:
			funcName = string(token)
		case ArgToken:
			args = append(args, identitor{
				name:  string(token),
				isVar: false,
			})
		case ArgVariableToken:
			args = append(args, identitor{
				name:  string(token),
				isVar: true,
			})
		}
	}
	if funcName != "" {
		t.identitors[len(t.snippets)] = identitor{
			name:  funcName,
			isVar: false,
		}
		t.args[len(t.snippets)] = args
	}
	return nil
}

// New returns raw string to template
func New(raw string) (Template, error) {
	t := template{
		identitors: map[int]identitor{},
		args:       map[int][]identitor{},
	}
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
func (t *template) Render(vs jsonutil.VariableMap) (string, error) {
	out := ""
	for i, snippet := range t.snippets {
		identitor, ok := t.identitors[i]
		if !ok {
			out += snippet
			continue
		}
		str, err := t.renderScript(&identitor, i, vs)
		if err != nil {
			return "", err
		}
		out += str
		out += snippet
	}
	index := len(t.snippets)
	identitor, ok := t.identitors[index]
	if !ok {
		return out, nil
	}
	str, err := t.renderScript(&identitor, index, vs)
	if err != nil {
		return "", err
	}
	out += str

	return out, nil
}

func (t *template) renderScript(ident *identitor, index int, vs jsonutil.VariableMap) (string, error) {
	if ident.isVar {
		names := strings.Split(ident.name, ".")
		v, err := vs.Select(names...)
		if err != nil {
			return "", fmt.Errorf("render %v err: %v", ident.name, err)
		}
		return v.String(), nil
	}
	args := t.args[index]
	funcArgs := []jsonutil.Variable{}
	for _, arg := range args {
		var funcArg jsonutil.Variable
		if arg.isVar {
			argNames := strings.Split(arg.name, ".")
			v, err := vs.Select(argNames...)
			if err != nil {
				funcArg = nil
			} else {
				funcArg = v
			}
		} else {
			funcArg = jsonutil.NewStringVariable("", arg.name)
		}
		funcArgs = append(funcArgs, funcArg)

	}
	s, err := Call(ident.name, funcArgs...)
	if err != nil {
		return "", fmt.Errorf("render %v(%v) err: %v", ident.name, join(funcArgs, ", "), err)
	}
	return s, nil
}

func join(vs []jsonutil.Variable, sep string) string {
	ss := make([]string, 0, len(vs))
	for _, v := range vs {
		ss = append(ss, v.String())
	}
	return strings.Join(ss, sep)
}
