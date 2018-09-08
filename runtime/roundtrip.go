package runtime

import (
	"time"
)

// RoundTripTemplate defines template of round trip
type RoundTripTemplate struct {
	// Request defines http request
	Request Request

	// Response defines http response validator
	Response Response
}

// RoundTrip defines a http round trip
type RoundTrip struct {
	// RoundTripTemplate defines round trip template
	RoundTripTemplate

	// Definitions defines variables from response
	Definitions []Definition
}

// Request defines http request
type Request struct {
	// Scheme defines scheme of http request
	// e.g. http or https
	Scheme string

	// Host defines host of http request
	Host string

	// Method defines http method, e.g. GET
	Method string

	// Path defines http request path
	Path string

	// Headers defines http request header
	Headers map[string]string

	// Body defines http request body
	Body []byte
}

// Response defines http response
type Response struct {
	// StatusCode checks response code
	StatusCode int

	// Headers defines http request header
	Headers map[string]string

	// Body defines http request body
	Body []byte

	// Async defines whether the task is a async task
	Async bool

	// Timeout defines deadline of checking
	Timeout time.Duration

	// Interval defines interval of polling and checking
	Interval time.Duration
}

// Definition defines variable definitions
type Definition struct {
	// Name defines variable name
	Name string

	// Selector defines variable selector
	Selector []string
}
