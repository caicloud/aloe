package framework

import (
	"fmt"

	"github.com/caicloud/aloe/roundtrip"
	"github.com/caicloud/aloe/template"
	"github.com/caicloud/aloe/types"
)

func (gf *genericFramework) constructContext(ctx *types.Context, ctxConfig *types.ContextConfig) (map[string]template.Variable, error) {
	if ctx.Error != nil {
		return nil, ctx.Error
	}

	newVs := map[string]template.Variable{}
	for k, v := range ctx.Variables {
		newVs[k] = v
	}
	newCtx := types.Context{
		Variables: newVs,
	}
	for _, rt := range ctxConfig.Flow {
		respMatcher, err := roundtrip.MatchResponse(&newCtx, &rt)
		if err != nil {
			return nil, err
		}
		resp, err := gf.client.DoRequest(&newCtx, &rt)
		if err != nil {
			return nil, err
		}
		matched, err := respMatcher.Match(resp)
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, fmt.Errorf(respMatcher.FailureMessage(resp))
		}
		vs, err := respMatcher.Variables()
		if err != nil {
			return nil, err
		}
		for k, v := range vs {
			newCtx.Variables[k] = v
		}
	}

	return newCtx.Variables, nil
}
