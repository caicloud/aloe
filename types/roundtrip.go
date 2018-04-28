package types

import (
	"strconv"
	"time"

	"github.com/caicloud/aloe/template"
)

// RoundTrip defines a test case
// It usually means one http request and response
type RoundTrip struct {
	// Description describe the round trip
	Description string `json:"description,omitempty"`

	// Request defines a http request template
	Request Request `json:"request,omitempty"`

	// Response defines a http response checker
	Response Response `json:"response,omitempty"`

	// Definitions defines new variables from response
	Definitions []Definition `json:"definitions,omitempty"`
}

// Request defines a part template of http request
type Request struct {
	// API is a http verb + http path
	// e.g GET /api/v1/users
	API *Template `json:"api"`

	// Headers defines http header of request
	// NOTE(liubog2008): whether to use map[string][]string
	Headers map[string]string `json:"headers,omitempty"`

	// Body defines a template with variable
	Body *Template `json:"body,omitempty"`
}

// Response defines a http response checker
type Response struct {
	// StatusCode checks response code
	StatusCode int `json:"statusCode"`

	// Headers defines http header of request
	// NOTE(liubog2008): whether to use map[string][]string
	Headers map[string]string `json:"headers,omitempty"`

	// Body is also a template like request body
	// It can be used to generate a matcher which
	// can test response body
	Body *Template `json:"body,omitempty"`

	// Eventually defines an async checker for response
	// It means response will eventually be matched
	Eventually *Eventually `json:"eventually,omitempty"`
}

// Definition defines new variable from response
type Definition struct {
	// Name defines variable name
	Name string `json:"name"`

	// Selector select variable value from response
	Selector []string `json:"selector"`
}

// Template is used to get template from json
type Template struct {
	template.Template

	raw []byte
}

// Eventually defines config for eventually
type Eventually struct {
	// Timeout defines deadline of checking
	Timeout *Duration `json:"timeout,omitempty"`

	// Interval defines interval of polling and checking
	// Default interval is 1 second
	Interval *Duration `json:"interval,omitempty"`
}

// Duration defines duration can be unmarshal from json
type Duration struct {
	time.Duration
}

// MarshalJSON implements json.Marshaler
func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

// UnmarshalJSON implements json.Marshaler
func (d *Duration) UnmarshalJSON(body []byte) error {
	dur, err := strconv.Unquote(string(body))
	if err != nil {
		return err
	}
	nd, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}
	d.Duration = nd
	return nil
}

// MarshalJSON implements json.Marshaler
func (t *Template) MarshalJSON() ([]byte, error) {
	return t.raw, nil
}

// UnmarshalJSON implements json.Marshaler
func (t *Template) UnmarshalJSON(body []byte) error {
	s, err := strconv.Unquote(string(body))
	if err != nil {
		return err
	}
	templ, err := template.New(s)
	if err != nil {
		return err
	}
	t.Template = templ
	t.raw = []byte(s)
	return nil
}
