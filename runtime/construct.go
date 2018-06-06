package runtime

import (
	"github.com/caicloud/aloe/utils/jsonutil"
)

// CopyContext copy content of context from src context
func CopyContext(dest, src *Context) {
	dest.Variables = copyVariables(src.Variables)
	dest.RoundTripTemplate = CopyRoundTripTemplate(src.RoundTripTemplate)
	dest.Presetters = nil
	dest.Cleaners = nil
}

func copyVariables(variables map[string]jsonutil.Variable) map[string]jsonutil.Variable {
	if variables == nil {
		return nil
	}
	vs := make(map[string]jsonutil.Variable)
	for k, v := range variables {
		vs[k] = v
	}
	return vs
}

// CopyRoundTripTemplate will return a copy of round trip
func CopyRoundTripTemplate(rt *RoundTripTemplate) *RoundTripTemplate {
	if rt == nil {
		return nil
	}
	// shallow copy roundtrip
	nrt := *rt

	nrt.Request.Headers = copyHeader(rt.Request.Headers)
	nrt.Response.Headers = copyHeader(rt.Response.Headers)

	return &nrt
}

func copyHeader(header map[string]string) map[string]string {
	if header == nil {
		return nil
	}
	nh := map[string]string{}
	for k, v := range header {
		nh[k] = v
	}
	return nh
}
