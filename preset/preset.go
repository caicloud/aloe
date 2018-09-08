package preset

import "github.com/caicloud/aloe/runtime"

// Presetter defines presetter
type Presetter interface {
	// Name defines name of presetter
	Name() string

	// Preset parse args and set roundtrip template
	Preset(rt *runtime.RoundTripTemplate, args map[string]string) (*runtime.RoundTripTemplate, error)
}
