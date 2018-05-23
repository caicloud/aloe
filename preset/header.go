package preset

import (
	"net/http"

	"github.com/caicloud/aloe/types"
)

const (
	// RequestType defines header presetter type
	RequestType string = "request"
	// ResponseType defines header presetter type
	ResponseType string = "response"
)

type headerPresetter struct {
	typ string
}

// NewHeaderPresetter returns header presetter
func NewHeaderPresetter(typ string) Presetter {
	return &headerPresetter{
		typ: typ,
	}
}

// Name implements preset.Presetter
func (p *headerPresetter) Name() string {
	return p.typ + "Header"
}

// Preset implements preset.Presetter
func (p *headerPresetter) Preset(rt *types.RoundTrip, args map[string]string) (*types.RoundTrip, error) {
	if rt == nil {
		rt = &types.RoundTrip{}
	}

	if p.typ == RequestType {
		if rt.Request.Headers == nil {
			rt.Request.Headers = map[string]string{}
		}

		for k, v := range args {
			rt.Request.Headers[http.CanonicalHeaderKey(k)] = v
		}
	}
	if p.typ == ResponseType {
		if rt.Response.Headers == nil {
			rt.Response.Headers = map[string]string{}
		}

		for k, v := range args {
			rt.Response.Headers[http.CanonicalHeaderKey(k)] = v
		}
	}

	return rt, nil
}
