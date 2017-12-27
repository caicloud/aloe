package types

import "github.com/caicloud/aloe/template"

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

	// Preset defines some common fields for each round-trip in context
	Preset RoundTrip `json:"preset,omitempty"`

	// Flow will be called to construct context
	Flow []RoundTrip `json:"flow,omitempty"`
}

// Context defines context of test cases
type Context struct {
	Variables map[string]template.Variable

	Error error
}
