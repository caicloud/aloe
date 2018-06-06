package framework

import (
	"fmt"
	"time"

	"github.com/caicloud/aloe/roundtrip"
	"github.com/caicloud/aloe/runtime"
	"github.com/caicloud/aloe/types"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var (
	// defaultTimeout defines default timeout of an async task
	defaultTimeout = 1 * time.Second

	// defaultInterval defines default interval of an async task checking
	defaultInterval = 100 * time.Millisecond
)

func (gf *genericFramework) constructFlow(ctx *runtime.Context, flow []types.RoundTrip) {
	for _, rt := range flow {
		gf.roundTrip(ctx, &rt, false)
	}
}

func (gf *genericFramework) constructValidatedFlow(ctx *runtime.Context, flow []types.RoundTripTuple) {
	for _, rts := range flow {
		failed := false
		for _, rt := range rts.Constructor {
			notMatched, err := gf.tryRoundTrip(ctx, &rt, false)
			// ignore round trip error
			if notMatched {
				failed = true
				break
			}
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}
		if !failed {
			continue
		}
		for _, rt := range rts.Validator {
			gf.roundTrip(ctx, &rt, false)
		}
	}
}

func (gf *genericFramework) constructRoundTripTemplate(ctx *runtime.Context) error {
	for _, pc := range ctx.Presetters {
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

func (gf *genericFramework) tryRoundTrip(ctx *runtime.Context, origin *types.RoundTrip, overwrite bool) (bool, error) {
	rt, err := runtime.RenderRoundTrip(ctx, origin)
	if err != nil {
		return false, err
	}

	ginkgo.By(origin.Description)

	respMatcher, err := roundtrip.MatchResponse(rt)
	if err != nil {
		return false, err
	}
	// TODO(liubog2008): support async tasks
	if rt.Response.Async {
		return false, fmt.Errorf("async round trip is not supported for constructor now")
	}
	resp, err := gf.client.DoRequest(rt)
	if err != nil {
		return false, err
	}
	matched, err := respMatcher.Match(resp)
	if err != nil {
		return false, err
	}
	if !matched {
		return true, fmt.Errorf(respMatcher.FailureMessage(resp))
	}
	vs := respMatcher.Variables()
	for k, v := range vs {
		if _, ok := ctx.Variables[k]; ok && !overwrite {
			return false, fmt.Errorf("variable %v has been defined", k)
		}
		ctx.Variables[k] = v
	}
	return false, nil
}

func (gf *genericFramework) roundTrip(ctx *runtime.Context, originRoundTrip *types.RoundTrip, overwrite bool) {
	rt, err := runtime.RenderRoundTrip(ctx, originRoundTrip)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By(originRoundTrip.Description)

	respMatcher, err := roundtrip.MatchResponse(rt)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	if rt.Response.Async {
		timeout := rt.Response.Timeout
		interval := rt.Response.Interval

		if timeout == time.Duration(0) {
			timeout = defaultTimeout
		}
		if interval == time.Duration(0) {
			interval = defaultInterval
		}

		gomega.Eventually(func() *roundtrip.Response {
			resp, err := gf.client.DoRequest(rt)
			return &roundtrip.Response{
				Resp: resp,
				Err:  err,
			}
		}, timeout, interval).Should(respMatcher)
	} else {
		resp, err := gf.client.DoRequest(rt)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(resp).To(respMatcher)
	}
	vs := respMatcher.Variables()
	for k, v := range vs {
		_, ok := ctx.Variables[k]
		gomega.Expect(ok && !overwrite).To(gomega.BeFalse(), "variable %v has been defined", k)
		ctx.Variables[k] = v
	}
}
