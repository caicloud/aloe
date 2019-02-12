package runtime

import "github.com/caicloud/aloe/utils/jsonutil"

// CopyContext copy content of context from src context
func CopyContext(dest, src *Context) {
	dest.Summary = src.Summary
	if src.Variables == nil {
		dest.Variables = jsonutil.NewVariableMap("", nil)
	} else {
		dest.Variables = src.Variables.Copy()
	}
	dest.RoundTripTemplate = CopyRoundTripTemplate(src.RoundTripTemplate)
	dest.Presetters = nil
	dest.Cleaners = nil
}

// ReconstructContext will reconstruct context
// from parent and previous variable patch
func ReconstructContext(ctx *Context) error {
	vs, err := jsonutil.Merge(ctx.Parent.Variables, jsonutil.ConflictOption, true, ctx.Exports)
	if err != nil {
		return err
	}
	ctx.Variables = vs
	ctx.RoundTripTemplate = CopyRoundTripTemplate(ctx.Parent.RoundTripTemplate)
	ctx.Presetters = nil
	ctx.Cleaners = nil
	return nil
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
