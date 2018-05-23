package roundtrip

import "github.com/caicloud/aloe/types"

// MergeRoundTrip merge patch into origin and
// return a new round trip
// NOTE(liubog2008): Definitions will be just replaced but not
// merged
func MergeRoundTrip(origin *types.RoundTrip, patch *types.RoundTrip) *types.RoundTrip {
	new := CopyRoundTrip(origin)
	if new == nil {
		new = &types.RoundTrip{}
	}

	if patch.Description != "" {
		new.Description = patch.Description
	}
	if patch.Request.API != nil {
		new.Request.API = patch.Request.API
	}
	if new.Request.Headers == nil {
		new.Request.Headers = map[string]string{}
	}
	for k, v := range patch.Request.Headers {
		new.Request.Headers[k] = v
	}
	if patch.Request.Body != nil {
		new.Request.Body = patch.Request.Body
	}

	if patch.Response.StatusCode != 0 {
		new.Response.StatusCode = patch.Response.StatusCode
	}
	if new.Response.Headers == nil {
		new.Response.Headers = map[string]string{}
	}
	for k, v := range patch.Response.Headers {
		new.Response.Headers[k] = v
	}
	if patch.Response.Body != nil {
		new.Response.Body = patch.Response.Body
	}
	mergeEventually(new.Response.Eventually, patch.Response.Eventually)
	new.Definitions = patch.Definitions
	return new
}

// CopyRoundTrip will copy the round trip
func CopyRoundTrip(origin *types.RoundTrip) *types.RoundTrip {
	if origin == nil {
		return nil
	}
	// shallow copy
	new := *origin

	new.Request.Headers = map[string]string{}
	for k, v := range origin.Request.Headers {
		new.Request.Headers[k] = v
	}

	new.Response.Headers = map[string]string{}
	for k, v := range origin.Response.Headers {
		new.Response.Headers[k] = v
	}

	if origin.Response.Eventually != nil {
		new.Response.Eventually = &types.Eventually{
			Timeout:  origin.Response.Eventually.Timeout,
			Interval: origin.Response.Eventually.Interval,
		}
	}

	new.Definitions = append(new.Definitions, origin.Definitions...)
	return &new
}

func mergeEventually(origin *types.Eventually, patch *types.Eventually) {
	if origin == nil {
		origin = patch
		return
	}
	if patch == nil {
		return
	}
	if patch.Interval != nil {
		origin.Interval = patch.Interval
	}
	if patch.Timeout != nil {
		origin.Timeout = patch.Timeout
	}
}
