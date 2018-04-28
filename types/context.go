package types

import "github.com/caicloud/aloe/utils/jsonutil"

const (
	// ContextFile defines default filename of spec
	ContextFile = "_context.yaml"
)

// ContextConfig defines some configs for ginkgo.Describe
// or ginkgo.Context
type ContextConfig struct {
	// Summary used to display message in
	// ginkgo.Describe or ginkgo.Context
	Summary string `json:"summary"`

	// Description used to describe context of all cases
	Description string `json:"description,omitempty"`

	// Definitions defines variable in this context
	// Definitions map[string]string `json:"definitions,omitempty"`

	// Presetter preset some common fields of round-trip in context
	Presetter []PresetConfig `json:"presetter,omitempty"`

	// Flow will be called to construct context
	Flow []RoundTrip `json:"flow,omitempty"`

	// ValidatedFlow defines flow with validator
	ValidatedFlow []RoundTripTuple `json:"validatedFlow,omitempty"`

	// Cleaner defines cleaner of the context
	Cleaner string `json:"cleaner,omitempty"`
}

// Context defines context of test cases
type Context struct {
	// Variables defines variables the context has
	Variables map[string]jsonutil.Variable

	// RoundTripTemplate defines template of roundtrip
	RoundTripTemplate *RoundTrip

	// CleanerName defines the cleaner name of context
	CleanerName string
}

// RoundTripTuple defines a tuple of round trips
// It used to construct context and validate it
// Each case will try to validate the context, if false
// constructor will be called
type RoundTripTuple struct {
	// Constructor defines constructor roundtrip of context
	Constructor []RoundTrip `json:"constructor,omitempty"`
	// Validator defines validator of context
	// Normally validator is just used for trigger reconstruction
	// of context
	Validator []RoundTrip `json:"validator,omitempty"`
}
