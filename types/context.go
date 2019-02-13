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

	// Exports defines variables which can be access by children
	Exports []Var `json:"exports,omitempty"`

	// Cleaners defines cleaner of the context
	Cleaners []CleanerConfig `json:"cleaners,omitempty"`
}
