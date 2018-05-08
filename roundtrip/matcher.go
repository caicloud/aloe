package roundtrip

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/caicloud/aloe/matcher"
	"github.com/caicloud/aloe/types"
	"github.com/caicloud/aloe/utils/close"
	"github.com/caicloud/aloe/utils/indent"
	"github.com/caicloud/aloe/utils/jsonutil"
	"github.com/onsi/gomega/format"
	gomegatypes "github.com/onsi/gomega/types"
)

// ResponseHandler defines handler for response
// It implements gomegatypes.GomegaMatcher
type ResponseHandler interface {
	gomegatypes.GomegaMatcher

	// Variables returns variables defined in the round trip
	Variables() (map[string]jsonutil.Variable, error)
}

// ResponseMatcher defines a matcher to match http response
type ResponseMatcher struct {
	bodyMatcher gomegatypes.GomegaMatcher

	// emptyBody used to validate that body is empty
	emptyBody bool

	code int

	defs []types.Definition

	parsed bool

	vars map[string]jsonutil.Variable

	failures []error
}

// MatchResponse returns a response matcher
func MatchResponse(ctx *types.Context, rt *types.RoundTrip) (ResponseHandler, error) {
	respConf := rt.Response
	rm := &ResponseMatcher{
		code: respConf.StatusCode,
		defs: rt.Definitions,
	}
	if respConf.Body == nil {
		return rm, nil

	}
	matcherConf, err := respConf.Body.Render(ctx.Variables)
	if err != nil {
		return nil, err
	}

	if len(matcherConf) == 0 {
		rm.emptyBody = true
		return rm, nil
	}

	m, err := matcher.Parse(matcherConf)
	if err != nil {
		return nil, fmt.Errorf("parse json error: %v", err)
	}
	rm.bodyMatcher = m
	return rm, nil
}

// Variables returns variable of matcher
func (m *ResponseMatcher) Variables() (map[string]jsonutil.Variable, error) {
	if !m.parsed {
		return nil, fmt.Errorf("response should be matched before get variables")
	}
	return m.vars, nil
}

// Match implements gomegatypes.GomegaMatcher
func (m *ResponseMatcher) Match(actual interface{}) (bool, error) {
	resp, ok := actual.(*http.Response)
	if !ok {
		return false, fmt.Errorf("%v is type %T, expected response", actual, actual)
	}
	defer close.Close(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.failures = append(m.failures, fmt.Errorf("can't read body from response"))
		return false, nil
	}
	if resp.StatusCode != m.code {
		m.failures = append(m.failures, fmt.Errorf("status code is not matched, expected: %v, actual: %v", m.code, resp.StatusCode))
		m.failures = append(m.failures, fmt.Errorf("api status: %v", string(body)))
	}

	if m.emptyBody && len(body) != 0 {
		m.failures = append(m.failures, fmt.Errorf("body should be empty, actual: %v", string(body)))
	}

	if m.bodyMatcher != nil {
		b := map[string]interface{}{}
		if err := json.Unmarshal(body, &b); err != nil {
			m.failures = append(m.failures, fmt.Errorf("can't unmarshal body to json, NOW only json Content-Type is supported"))
			return false, nil

		}
		if err := func() error {
			matched, err := m.bodyMatcher.Match(b)
			if err != nil {
				return err
			} else if !matched {
				return errors.New(
					indent.Indent(m.bodyMatcher.FailureMessage(b), "\t"),
				)
			}
			return nil
		}(); err != nil {
			// TODO(zjj2wry): print got data when 'want' not match 'got'
			m.failures = append(m.failures, fmt.Errorf("can't match response body: \n%v", err))
		}
	}
	if len(m.failures) > 0 {
		return false, nil
	}

	m.vars = map[string]jsonutil.Variable{}
	isErr := false
	for _, def := range m.defs {
		v, err := jsonutil.GetVariable(body, def.Name, def.Selector...)
		if err != nil {
			m.failures = append(m.failures, err)
			isErr = true
		}
		m.vars[def.Name] = v
	}
	if isErr {
		return false, nil
	}
	m.parsed = true
	return true, nil
}

// FailureMessage implements gomegatypes.GomegaMatcher
func (m *ResponseMatcher) FailureMessage(actual interface{}) (message string) {
	failures := make([]string, len(m.failures))
	for i := range m.failures {
		failures[i] = m.failures[i].Error()
	}

	// can't add response directly, ignored now
	// TODO(liubog2008): Fix it or use more readable response format
	return format.Message("",
		fmt.Sprintf("to match response: {\n%v\n}\n", strings.Join(failures, "\n")))

}

// NegatedFailureMessage implements gomegatypes.GomegaMatcher
func (m *ResponseMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(reflect.TypeOf(actual).Name(), "not to match response")
}
