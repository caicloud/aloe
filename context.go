package aloe

import (
	"fmt"

	"github.com/caicloud/aloe/roundtrip"
	"github.com/caicloud/aloe/types"
	"github.com/caicloud/aloe/utils/jsonutil"
)

func (gf *genericFramework) constructContext(ctx *types.Context, ctxConfig *types.ContextConfig) error {
	if ctx.Variables == nil {
		ctx.Variables = map[string]jsonutil.Variable{}
	}
	if err := gf.preset(ctx, ctxConfig.Presetter); err != nil {
		return err
	}
	for _, rt := range ctxConfig.Flow {

		_, vs, err := gf.roundTrip(ctx, &rt)
		if err != nil {
			return err
		}
		for k, v := range vs {
			if _, ok := ctx.Variables[k]; ok {
				return fmt.Errorf("variable %v has been defined", k)
			}
			ctx.Variables[k] = v
		}
	}
	for _, rts := range ctxConfig.ValidatedFlow {
		for _, rt := range rts.Constructor {
			_, vs, err := gf.roundTrip(ctx, &rt)
			if err != nil {
				return err
			}
			for k, v := range vs {
				if _, ok := ctx.Variables[k]; ok {
					return fmt.Errorf("variable %v has been defined", k)
				}
				ctx.Variables[k] = v
			}
		}
		for _, rt := range rts.Validator {
			_, vs, err := gf.roundTrip(ctx, &rt)
			if err != nil {
				return err
			}
			for k, v := range vs {
				if _, ok := ctx.Variables[k]; ok {
					return fmt.Errorf("variable %v has been defined", k)
				}
				ctx.Variables[k] = v
			}
		}
	}

	ctx.CleanerName = ctxConfig.Cleaner
	return nil
}

func (gf *genericFramework) reconstructContext(ctx *types.Context, ctxConfig *types.ContextConfig) error {
	if ctx.Variables == nil {
		ctx.Variables = map[string]jsonutil.Variable{}
	}
	if err := gf.preset(ctx, ctxConfig.Presetter); err != nil {
		return err
	}
	for _, rts := range ctxConfig.ValidatedFlow {
		shouldReconstruct := false
		for _, rt := range rts.Validator {
			notMatched, _, err := gf.roundTrip(ctx, &rt)
			if err != nil && !notMatched {
				return err
			}
			if notMatched {
				shouldReconstruct = true
				break
			}
		}
		if shouldReconstruct {
			for _, rt := range rts.Constructor {
				_, vs, err := gf.roundTrip(ctx, &rt)
				if err != nil {
					return err
				}
				for k, v := range vs {
					ctx.Variables[k] = v
				}
			}
		}
	}
	ctx.CleanerName = ctxConfig.Cleaner
	return nil
}

func (gf *genericFramework) preset(ctx *types.Context, pcs []types.PresetConfig) error {
	for _, pc := range pcs {
		p, ok := gf.presetters[pc.Name]
		if !ok {
			return fmt.Errorf("can't get presetter called %v", pc.Name)
		}
		rt, err := p.Preset(ctx.RoundTripTemplate, pc.Args)
		if err != nil {
			return err
		}
		ctx.RoundTripTemplate = rt
	}
	return nil
}

func (gf *genericFramework) roundTrip(ctx *types.Context, originRoundTrip *types.RoundTrip) (bool, map[string]jsonutil.Variable, error) {
	rt := roundtrip.MergeRoundTrip(ctx.RoundTripTemplate, originRoundTrip)

	respMatcher, err := roundtrip.MatchResponse(ctx, rt)
	if err != nil {
		return false, nil, err
	}
	resp, err := gf.client.DoRequest(ctx, rt)
	if err != nil {
		return false, nil, err
	}
	matched, err := respMatcher.Match(resp)
	if err != nil {
		return false, nil, err
	}
	if !matched {
		return true, nil, fmt.Errorf(respMatcher.FailureMessage(resp))
	}
	vs, err := respMatcher.Variables()
	if err != nil {
		return false, nil, err
	}
	return false, vs, nil
}

func saveContext(origin *types.Context) *types.Context {
	copy := types.Context{}
	copy.Variables = map[string]jsonutil.Variable{}
	for k, v := range origin.Variables {
		copy.Variables[k] = v
	}
	copy.RoundTripTemplate = roundtrip.CopyRoundTrip(origin.RoundTripTemplate)
	copy.CleanerName = origin.CleanerName
	return &copy
}

func restoreContext(origin *types.Context, copy *types.Context) {
	origin.Variables = copy.Variables
	origin.RoundTripTemplate = copy.RoundTripTemplate
	origin.CleanerName = copy.CleanerName
}
