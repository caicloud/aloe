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
	"github.com/caicloud/aloe/runtime"
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
	Variables() map[string]jsonutil.Variable
}

// ResponseMatcher defines a matcher to match http response
type ResponseMatcher struct {
	bodyMatcher gomegatypes.GomegaMatcher

	// emptyBody used to validate that body is empty
	emptyBody bool

	code int

	headers map[string]string

	defs []runtime.Definition

	parsed bool

	vars map[string]jsonutil.Variable

	failures []error
}

// MatchResponse returns a response matcher
func MatchResponse(rt *runtime.RoundTrip) (ResponseHandler, error) {
	if rt == nil {
		return nil, fmt.Errorf("empty roundtrip")
	}
	resp := rt.Response
	rm := &ResponseMatcher{
		code:    resp.StatusCode,
		headers: resp.Headers,
		defs:    rt.Definitions,
	}
	if resp.Body == nil {
		return rm, nil
	}

	if len(resp.Body) == 0 {
		rm.emptyBody = true
		return rm, nil
	}

	m, err := matcher.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse json error: %v", err)
	}
	rm.bodyMatcher = m
	return rm, nil
}

// Variables returns variable of matcher
func (m *ResponseMatcher) Variables() map[string]jsonutil.Variable {
	if !m.parsed {
		return nil
	}
	return m.vars
}

// Response used for match
type Response struct {
	Resp *http.Response
	Err  error
}

// Match implements gomegatypes.GomegaMatcher
func (m *ResponseMatcher) Match(actual interface{}) (bool, error) {
	var resp *http.Response
	switch obj := actual.(type) {
	case *http.Response:
		resp = obj
	case *Response:
		if obj.Err != nil {
			return false, obj.Err
		}
		resp = obj.Resp
	default:
		return false, fmt.Errorf("%v is type %T, expected *http.Response or *roundtrip.Response", actual, actual)
	}
	defer close.Close(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.failures = append(m.failures, fmt.Errorf("can't read body from response"))
		return false, nil
	}
	if m.code != 0 && resp.StatusCode != m.code {
		m.failures = append(m.failures, fmt.Errorf("status code is not matched, expected: %v, actual: %v", m.code, resp.StatusCode))
		m.failures = append(m.failures, fmt.Errorf("api status: %v", string(body)))
	}

	for k, v := range m.headers {
		hv := resp.Header.Get(k)
		if hv != v {
			m.failures = append(m.failures, fmt.Errorf("response header %v is not matched, expected: %v, actual: %v", k, v, hv))
		}
	}

	if m.emptyBody && len(body) != 0 {
		m.failures = append(m.failures, fmt.Errorf("body should be empty, actual: %v", string(body)))
	}

	if m.bodyMatcher != nil {
		var b interface{}
		if err := json.Unmarshal(body, &b); err != nil {
			m.failures = append(m.failures, fmt.Errorf("can't unmarshal body(%v) to json, NOW only json Content-Type is supported", string(body)))
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
		switch def.Type {
		case runtime.BodyType:
			v, err := jsonutil.GetVariable(body, def.Name, def.Selector...)
			if err != nil {
				m.failures = append(m.failures, err)
				isErr = true
			}
			m.vars[def.Name] = v
		case runtime.StatusType:
			m.vars[def.Name] = jsonutil.NewIntVariable(def.Name, int64(resp.StatusCode))
		case runtime.HeaderType:
			if len(def.Selector) != 1 {
				m.failures = append(m.failures, fmt.Errorf("header definition expected selector with len 1, actual is %v", def.Selector))
			}
			m.vars[def.Name] = jsonutil.NewStringVariable(def.Name, resp.Header.Get(def.Selector[0]))
		}
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
