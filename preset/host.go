package preset

import (
	"fmt"

	"github.com/caicloud/aloe/runtime"
)

type hostPresetter struct {
}

// NewHostPresetter returns host presetter
func NewHostPresetter() Presetter {
	return &hostPresetter{}
}

// Name implements preset.Presetter
func (p *hostPresetter) Name() string {
	return "host"
}

// Preset implements preset.Presetter
func (p *hostPresetter) Preset(rt *runtime.RoundTripTemplate, args map[string]string) (*runtime.RoundTripTemplate, error) {
	if rt == nil {
		rt = &runtime.RoundTripTemplate{}
	}

	host, ok := args["host"]
	if !ok {
		return nil, fmt.Errorf("host is not defined")
	}
	scheme, ok := args["scheme"]
	if !ok {
		scheme = "http"
	}

	rt.Request.Host = host
	rt.Request.Scheme = scheme

	return rt, nil
}
