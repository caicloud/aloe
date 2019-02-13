package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/onsi/gomega/format"
	errorsutil "github.com/onsi/gomega/gstruct/errors"
	"github.com/onsi/gomega/types"
)

// Simplify element matcher
// See https://github.com/onsi/gomega/blob/master/gstruct/elements.go

// MatchSlice succeeds if every element of a slice matches the element matcher it maps to through the id function, and every element matcher is matched.
func MatchSlice(elements Elements) types.GomegaMatcher {
	m := &SliceMatcher{
		Elements: elements,
	}
	return m
}

// SliceMatcher is a NestingMatcher that applies custom matchers to each element of a slice mapped
// by the Identifier function.
// TODO: Extend this to work with arrays & maps (map the key) as well.
type SliceMatcher struct {
	// Matchers for each element.
	Elements Elements

	// State.
	failures []error
}

// Elements ID to matcher.
type Elements []types.GomegaMatcher

// Match implements gomega.Matcher
func (m *SliceMatcher) Match(actual interface{}) (success bool, err error) {
	if reflect.TypeOf(actual).Kind() != reflect.Slice {
		return false, fmt.Errorf("%v is type %T, expected slice", actual, actual)
	}

	m.failures = m.matchElements(actual)
	if len(m.failures) > 0 {
		return false, nil
	}
	return true, nil
}

func (m *SliceMatcher) matchElements(actual interface{}) (errs []error) {
	// Provide more useful error messages in the case of a panic.
	defer func() {
		if err := recover(); err != nil {
			errs = append(errs, fmt.Errorf("panic checking %+v: %v\n%s", actual, err, debug.Stack()))
		}
	}()

	val := reflect.ValueOf(actual)
	length := val.Len()
	if len(m.Elements) != length {
		errs = append(errs, fmt.Errorf("unexpected slice length, expected: %v, actual: %v", len(m.Elements), length))
		return errs
	}
	for i := 0; i < length; i++ {
		element := val.Index(i).Interface()

		matcher := m.Elements[i]

		match, err := matcher.Match(element)
		if match {
			continue
		}

		if err == nil {
			if nesting, ok := matcher.(errorsutil.NestingMatcher); ok {
				err = errorsutil.AggregateError(nesting.Failures())
			} else {
				err = errors.New(matcher.FailureMessage(element))
			}
		}
		errs = append(errs, errorsutil.Nest(fmt.Sprintf("[%v]", i), err))
	}

	return errs
}

// FailureMessage implements types.GomegaMatcher
func (m *SliceMatcher) FailureMessage(actual interface{}) (message string) {
	failure := errorsutil.AggregateError(m.failures)
	return format.Message(actual, fmt.Sprintf("to match elements: %v", failure))
}

// NegatedFailureMessage implements types.GomegaMatcher
func (m *SliceMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match elements")
}

// Failures returns failures of matcher
func (m *SliceMatcher) Failures() []error {
	return m.failures
}
