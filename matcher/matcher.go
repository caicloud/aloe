package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
)

func generateMapMatcher(matcher map[string]interface{}) (gomegatypes.GomegaMatcher, error) {
	fields := Fields{}
	exists := map[string]bool{}
	for k, expr := range matcher {
		if expr == nil {
			fields[k] = gomega.BeNil()
			continue
		}
		m, ok := convertToMap(expr)
		if ok {
			matcher, exist, err := generateSpMatcher(m)
			if err != nil {
				return nil, err
			}
			switch exist {
			case fieldNotExist:
				exists[k] = false
			case fieldExist:
				exists[k] = true
			}
			fields[k] = matcher
			continue
		}

		matcher, err := generateMatcher(expr)
		if err != nil {
			return nil, err
		}
		fields[k] = matcher
	}
	return MatchMap(fields, exists), nil
}

const (
	// matcher all keys begin with $
	specialMatcher = 0
	// matcher no key begins with $
	normalMatcher = 1
	// matcher mixed above
	mixedMatcher = -1
	// matcher without key
	unknownMatcher = -2
)

//  0: sp
//  1: normal
// -1: mixed
func checkSpMatcher(m map[string]interface{}) int {
	res := unknownMatcher
	for k := range m {
		if strings.HasPrefix(k, "$") {
			switch res {
			case unknownMatcher:
				res = specialMatcher
			case normalMatcher:
				res = mixedMatcher
			}
		} else {
			switch res {
			case unknownMatcher:
				res = normalMatcher
			case specialMatcher:
				res = mixedMatcher
			}
		}
	}
	if res == unknownMatcher {
		res = normalMatcher
	}
	return res
}

const (
	fieldUnkown   = -1
	fieldNotExist = 0
	fieldExist    = 1
)

func generateSpMatcher(m map[string]interface{}) (gomegatypes.GomegaMatcher, int, error) {
	res := checkSpMatcher(m)
	if res == mixedMatcher {
		return nil, fieldUnkown, fmt.Errorf("mixed special matcher and fields")
	}
	if res == normalMatcher {
		m, err := generateMapMatcher(m)
		return m, fieldUnkown, err
	}
	ms := map[string]gomegatypes.GomegaMatcher{}
	exist := fieldUnkown
	for key, expr := range m {
		switch key {
		case MatchMatcher:
			ma, err := generateMatcher(expr)
			if err != nil {
				return nil, fieldUnkown, err
			}
			ms[MatchMatcher] = ma
		case LenMatcher:
			ma, err := generateLenMatcher(expr)
			if err != nil {
				return nil, fieldUnkown, err
			}
			ms[LenMatcher] = ma

		case RegexpMatcher:
			ma, err := generateRegexpMatcher(expr)
			if err != nil {
				return nil, fieldUnkown, err
			}
			ms[RegexpMatcher] = ma
		case ExistsMatcher:
			b, ok := expr.(bool)
			if !ok {
				return nil, fieldUnkown, fmt.Errorf("value of $exists MUST be bool, actual: %T", expr)
			}

			if !b {
				if len(m) != 1 {
					return nil, fieldUnkown, fmt.Errorf("if $exists is false, all other matchers will be ignored")
				}
				exist = fieldNotExist
			} else {
				exist = fieldExist
			}
		default:
			return nil, fieldUnkown, fmt.Errorf("unknown special matcher: %v", key)
		}
	}
	return MatchSpecial(ms), exist, nil
}

func generateLenMatcher(expr interface{}) (gomegatypes.GomegaMatcher, error) {
	f, ok := expr.(float64)
	if !ok {
		return nil, fmt.Errorf("value of $len MUST be an int, actual: %T", expr)
	}
	n := int(f)
	if float64(n) != f {
		return nil, fmt.Errorf("value of $len MUST be an int, actual: %v", expr)
	}
	return gomega.HaveLen(n), nil
}

func generateRegexpMatcher(expr interface{}) (gomegatypes.GomegaMatcher, error) {
	s, ok := expr.(string)
	if !ok {
		return nil, fmt.Errorf("value of $regexp MUST be a string, actual: %T", expr)
	}
	return gomega.MatchRegexp(s), nil
}

func convertToMap(expr interface{}) (map[string]interface{}, bool) {
	m, ok := expr.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return m, true
}

const (
	// ExistsMatcher defines special matcher to match non-existent key
	// e.g.
	// matcher:
	// {
	//   "string": {
	//     "$exists": true
	//   }
	// }
	// data:
	// {
	//   "string": "string"
	// }
	ExistsMatcher = "$exists"

	// RegexpMatcher defines matcher to match regexp
	// e.g.
	// matcher:
	// {
	//   "string": {
	//     "$regexp": "[a-z]*"
	//   }
	// }
	// data:
	// {
	//   "string": "string"
	// }
	RegexpMatcher = "$regexp"

	// MatchMatcher defines matcher matches original data
	// e.g.
	// matcher:
	// {
	//   "string": {
	//     "$match": "string"
	//   }
	// }
	// data:
	// {
	//   "string": "string"
	// }
	MatchMatcher = "$match"

	// LenMatcher defines matcher matches length of data
	// support string, object, array
	// e.g.
	// matcher:
	// {
	//   "string": {
	//     "$len": 4
	//   }
	// }
	// data:
	// {
	//   "string": "1234"
	// }
	LenMatcher = "$len"
)

func generateSliceMatcher(matcher []interface{}) (gomegatypes.GomegaMatcher, error) {
	elems := Elements{}
	for _, expr := range matcher {
		elem, err := generateMatcher(expr)
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}
	return MatchSlice(elems), nil
}

func generateMatcher(expr interface{}) (gomegatypes.GomegaMatcher, error) {
	t := reflect.TypeOf(expr)
	switch t.Kind() {
	case reflect.String, reflect.Bool:
		return gomega.Equal(expr), nil
	case reflect.Float64:
		return gomega.BeNumerically("==", expr), nil
	case reflect.Array, reflect.Slice:
		s, ok := expr.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expr type %T is a slice(array) but can't be []interface{}", expr)
		}
		return generateSliceMatcher(s)
	case reflect.Map:
		m, ok := expr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expr type %T is a map but can't be map[string]interface{}", expr)
		}

		ma, _, err := generateSpMatcher(m)
		if err != nil {
			return nil, err
		}
		return ma, nil

	}
	return nil, fmt.Errorf("unexpected type %T: all kinds are from json.Unmarshal", expr)
}

// Parse parse matcher of response and returns GomegaMatcher
func Parse(matcher []byte) (gomegatypes.GomegaMatcher, error) {
	var m interface{}
	if err := json.Unmarshal(matcher, &m); err != nil {
		return nil, err
	}
	return generateMatcher(m)
}
