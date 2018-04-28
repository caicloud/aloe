package preset

import "github.com/caicloud/aloe/types"

// Presetter defines presetter
type Presetter interface {
	// Name defines name of presetter
	Name() string

	// Preset parse args and set
	// TODO(liubog2008): define a round trip template type
	Preset(rt *types.RoundTrip, args map[string]string) (*types.RoundTrip, error)
}
