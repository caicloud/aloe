package template

import (
	"fmt"
	"reflect"
	"strconv"
)

// Argument defines function argument
type Argument interface {
	// IsNil returns whether argument is nil
	IsNil() bool
	// Eval eval argument into obj
	Eval(obj interface{}) (bool, error)
}

// FromString defines object which value can be parsed by String
type FromString interface {
	// Parse parse string into object
	Parse(s string) error
}

type arg struct {
	value string
	isNil bool
}

func (a *arg) IsNil() bool {
	return a.isNil
}

func (a *arg) Eval(obj interface{}) (bool, error) {
	if a.isNil {
		return false, nil
	}
	fs, ok := obj.(FromString)
	if ok {
		if err := fs.Parse(a.value); err != nil {
			return false, err
		}
		return true, nil
	}
	v := reflect.ValueOf(obj).Elem()
	switch v.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(a.value)
		if err != nil {
			return false, err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(a.value, 10, 64)
		if err != nil {
			return false, err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(a.value, 10, 64)
		if err != nil {
			return false, err
		}
		v.SetUint(u)
	case reflect.String:
		v.SetString(a.value)
	default:
		return false, fmt.Errorf("unsupported kind %v, please implement FromString interface", v.Kind())
	}
	return true, nil
}

// NewArgument returns a generic argument
func NewArgument(s string, isNil bool) Argument {
	return &arg{
		value: s,
		isNil: isNil,
	}
}

// Call calls function of template
func Call(name string, args ...Argument) (string, error) {
	switch name {
	case Random:
		if len(args) == 0 {
			return "", fmt.Errorf("func random expected 1 or 2 args, but received: %v", len(args))
		}
		var regexp string
		exist, err := args[0].Eval(&regexp)
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("first argument of random is nil")
		}
		if len(args) == 1 {
			return random(regexp)
		}
		var limit int
		hasSecond, err := args[1].Eval(&limit)
		if err != nil {
			return "", err
		}
		if !hasSecond {
			return "", fmt.Errorf("second argument of random is nil")
		}
		if len(args) == 2 {
			return randomWithLimit(regexp, limit)
		}
		return "", fmt.Errorf("func random expected 1 or 2 args, but received: %v", len(args))
	case Exist:
		if len(args) != 1 {
			return "", fmt.Errorf("func exist expected 1 arg, but received: %v", len(args))
		}
		return isExist(args[0])
	default:
		return "", fmt.Errorf("unknown function named %v", name)
	}
}
