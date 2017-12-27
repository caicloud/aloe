package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"

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
		isSpMatcher := false
		typ := reflect.TypeOf(expr)
		if typ.Kind() == reflect.Map {
			m, err := convertToMap(expr)
			if err != nil {
				return nil, err
			}
			for childKey, childExpr := range m {
				switch childKey {
				case RegexpMatcher:
					ma, err := generateRegexpMatcher(childExpr)
					if err != nil {
						return nil, err
					}
					fields[k] = ma
					isSpMatcher = true
				case ExistsMatcher:
					b, ok := childExpr.(bool)
					if !ok {
						return nil, fmt.Errorf("value of $exists MUST be bool, actual: %T", expr)
					}

					if len(m) != 1 {
						if !b {
							return nil, fmt.Errorf("if $exists is false, all other matchers will be ignored")
						}
						isSpMatcher = false
					} else {
						exists[k] = b
						isSpMatcher = true
					}
				}
			}
		}
		if !isSpMatcher {
			ma, err := generateMatcher(expr)
			if err != nil {
				return nil, err
			}
			fields[k] = ma
		}
	}
	return MatchMap(fields, exists), nil
}

func generateRegexpMatcher(expr interface{}) (gomegatypes.GomegaMatcher, error) {
	s, ok := expr.(string)
	if !ok {
		return nil, fmt.Errorf("value of $regexp MUST be a string, actual: %T", expr)
	}
	return gomega.MatchRegexp(s), nil
}

func convertToMap(expr interface{}) (map[string]interface{}, error) {
	m, ok := expr.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expr type %T is a map but can't be map[string]interface{}", expr)

	}
	return m, nil
}

const (
	// ExistsMatcher defines special matcher to match non-existant key
	ExistsMatcher = "$exists"

	// RegexpMatcher defines matcher to match regexp
	RegexpMatcher = "$regexp"
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

		return generateMapMatcher(m)
	}
	return nil, fmt.Errorf("unexpected type %T: all kinds are from json.Unmarshal", expr)
}

// Parse parse matcher of response and returns GomegaMatcher
func Parse(matcher string) (gomegatypes.GomegaMatcher, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(matcher), &m); err != nil {
		return nil, err
	}
	return generateMatcher(m)
}
