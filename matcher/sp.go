package matcher

import (
	"errors"
	"fmt"

	"github.com/onsi/gomega/format"
	errorsutil "github.com/onsi/gomega/gstruct/errors"
	"github.com/onsi/gomega/types"
)

func MatchSpecial(ms map[string]types.GomegaMatcher) types.GomegaMatcher {
	return &SpecialMatcher{
		ms: ms,
	}
}

// SpecialMatcher match one by one
type SpecialMatcher struct {
	ms map[string]types.GomegaMatcher

	failures []error
}

// Match implements types.GomegaMatcher
func (sp *SpecialMatcher) Match(actual interface{}) (bool, error) {
	for k, m := range sp.ms {
		match, err := m.Match(actual)
		if err != nil {
			return false, err
		}
		if match {
			continue
		}
		if nesting, ok := m.(errorsutil.NestingMatcher); ok {
			err = errorsutil.AggregateError(nesting.Failures())
		} else {
			err = errors.New(m.FailureMessage(actual))
		}
		sp.failures = append(sp.failures, errorsutil.Nest(k, err))
	}
	if len(sp.failures) > 0 {
		return false, nil
	}
	return true, nil
}

// FailureMessage implements types.GomegaMatcher
func (sp *SpecialMatcher) FailureMessage(actual interface{}) (message string) {
	failure := errorsutil.AggregateError(sp.failures)
	return format.Message(actual, fmt.Sprintf("to match by special matchers: %v", failure))
}

// NegatedFailureMessage implements types.GomegaMatcher
func (sp *SpecialMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match special matchers")
}

// Failures returns failures of matcher
func (sp *SpecialMatcher) Failures() []error {
	return sp.failures
}
