package runtime

import (
	"time"
)

// DefinitionType defines where definition is from
type DefinitionType string

const (
	// BodyType means definition from body
	BodyType DefinitionType = "body"
	// HeaderType means definition from header
	HeaderType DefinitionType = "header"
	// StatusType means definition from status
	StatusType DefinitionType = "status"
)

// RoundTripTemplate defines template of round trip
type RoundTripTemplate struct {
	// Client defines http client used by this roundtrip
	Client string

	// Request defines http request
	Request Request

	// Response defines http response validator
	Response Response
}

// RoundTrip defines a http round trip
type RoundTrip struct {
	// RoundTripTemplate defines round trip template
	RoundTripTemplate

	// When defines round trip condition
	When *When

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

	// PathTemplate defines template of path
	// call fmt.Sprintf(PathTemplate, path) to generate real API path
	PathTemplate string

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

// Var defines variable
type Var struct {
	// Name defines variable name
	Name string

	// Selector defines variable selector
	Selector []string
}

// Definition defines variable definitions
type Definition struct {
	// Var defines Definition variable
	Var
	// Type defines variable from
	// enum ["body", "status", "header"]
	// default is body
	Type DefinitionType
}

// When defines roundtrip condition
type When struct {
	// Expr defines condition expression
	Expr string

	// Args defines additional args
	Args map[string]string
}
