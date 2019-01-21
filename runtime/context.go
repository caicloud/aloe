package runtime

import (
	"github.com/caicloud/aloe/utils/jsonutil"
)

// Context defines context of test cases
type Context struct {
	// Summary defines context summary
	Summary string

	// Parent defines parent context
	Parent *Context

	// Variables defines variables the context has
	Variables jsonutil.VariableMap

	// Presetters defines the presetters of context
	Presetters []Presetter

	// Cleaners defines the cleaners of context
	Cleaners []Cleaner

	// RoundTripTemplate defines template of roundtrip
	RoundTripTemplate *RoundTripTemplate
}

// Presetter defines presetter args
type Presetter struct {
	// Name defines presetter name
	Name string
	// Args defines args of presetter
	Args map[string]string
}

// Cleaner defines cleaner args
type Cleaner struct {
	// Name defines name of cleaner
	Name string
	// ForEach defines whether cleaner is called for each case
	ForEach bool
	// Args defines cleaner args
	Args map[string]string
}
