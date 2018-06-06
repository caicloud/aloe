package types

const (
	// ContextFile defines default filename of spec
	ContextFile = "context.yaml"
)

// Context defines some configs for ginkgo.Describe
// or ginkgo.Context
type Context struct {
	// Summary used to display message in
	// ginkgo.Describe or ginkgo.Context
	Summary string `json:"summary"`

	// Definitions defines variable in this context
	// Definitions map[string]string `json:"definitions,omitempty"`

	// Presetters preset some common fields of round-trip in context
	Presetters []PresetConfig `json:"presetters,omitempty"`

	// Flow will be called to construct context
	Flow []RoundTrip `json:"flow,omitempty"`

	// ValidatedFlow defines flow with validator
	ValidatedFlow []RoundTripTuple `json:"validatedFlow,omitempty"`

	// Cleaners defines cleaner of the context
	Cleaners []CleanerConfig `json:"cleaners,omitempty"`
}

// RoundTripTuple defines a tuple of round trips
// It used to construct context and validate it
// Each case will try to validate the context, if false
// constructor will be called
type RoundTripTuple struct {
	// Constructor defines constructor roundtrip of context
	// Normally constructor will construct the context
	Constructor []RoundTrip `json:"constructor,omitempty"`

	// Validator defines validator of context
	// If constructor is failed, validator can be used to
	// ignore the error
	// Normally it used to assert that context is "clean"
	Validator []RoundTrip `json:"validator,omitempty"`
}
