package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/onsi/gomega/format"
	errorsutil "github.com/onsi/gomega/gstruct/errors"
	"github.com/onsi/gomega/types"
)

//MatchMap succeeds if every field of a struct matches the field matcher associated with
//it, and every element matcher is matched.
//  Expect([]string{"a", "b"}).To(MatchAllFields(idFn, gstruct.Fields{
//      "a": BeEqual("a"),
//      "b": BeEqual("b"),
//  })
func MatchMap(fields Fields, exists map[string]bool) types.GomegaMatcher {
	m := &FieldsMatcher{
		Fields: fields,
		Exists: map[string]bool{},
	}
	for k, v := range exists {
		m.Exists[k] = v
	}
	for name := range fields {
		if _, ok := exists[name]; !ok {
			m.Exists[name] = true
		}
	}
	return m
}

// FieldsMatcher for map
type FieldsMatcher struct {
	// Matchers for each field.
	Fields Fields

	// Exists defines existence of fields
	// If field is not in it, it will not be
	// checked for existence
	Exists map[string]bool

	// State.
	failures []error
}

// Fields name to matcher.
type Fields map[string]types.GomegaMatcher

// Match implements types.GomegaMatcher
func (m *FieldsMatcher) Match(actual interface{}) (success bool, err error) {
	typ := reflect.TypeOf(actual)
	if typ.Kind() != reflect.Map {
		return false, fmt.Errorf("%v is type %T, expected map", actual, actual)
	}
	if typ.Key().Kind() != reflect.String {
		return false, fmt.Errorf("%v is type %T, expected key of map is only support string", actual, actual)
	}

	m.failures = m.matchFields(actual)
	if len(m.failures) > 0 {
		return false, nil
	}
	return true, nil
}

func (m *FieldsMatcher) matchFields(actual interface{}) (errs []error) {
	val := reflect.ValueOf(actual)
	elemType := val.Type().Elem()
	fields := map[string]struct{}{}

	keys := val.MapKeys()
	for _, k := range keys {
		fieldName := k.String()
		fields[fieldName] = struct{}{}

		err := func() (err error) {
			// This test relies heavily on reflect, which tends to panic.
			// Recover here to provide more useful error messages in that case.
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic checking %+v: %v\n%s", actual, r, debug.Stack())
				}
			}()

			matcher, expected := m.Fields[fieldName]
			if !expected {
				return nil
			}

			var field interface{}
			fieldValue := val.MapIndex(k)

			if fieldValue.IsValid() {
				field = fieldValue.Interface()
			} else {
				field = reflect.Zero(elemType)
			}

			matched, err := matcher.Match(field)
			if err != nil {
				return err
			} else if !matched {
				if nesting, ok := matcher.(errorsutil.NestingMatcher); ok {
					return errorsutil.AggregateError(nesting.Failures())
				}
				return errors.New(matcher.FailureMessage(field))
			}
			return nil
		}()
		if err != nil {
			errs = append(errs, errorsutil.Nest("."+fieldName, err))
		}
	}

	for fieldName, exist := range m.Exists {
		_, ok := fields[fieldName]
		if exist != ok {
			errs = append(errs, errorsutil.Nest("."+fieldName, fmt.Errorf("field existence err, expected: %v, actual: %v", exist, ok)))
		}
	}

	return errs
}

// FailureMessage implements types.GomegaMatcher
func (m *FieldsMatcher) FailureMessage(actual interface{}) (message string) {
	failures := make([]string, len(m.failures))
	for i := range m.failures {
		failures[i] = m.failures[i].Error()
	}
	return format.Message(actual,
		fmt.Sprintf("to match fields: {\n%v\n}\n", strings.Join(failures, "\n")))
}

// NegatedFailureMessage implements types.GomegaMatcher
func (m *FieldsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match fields")
}

// Failures returns all errors
func (m *FieldsMatcher) Failures() []error {
	return m.failures
}
