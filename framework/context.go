package framework

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/caicloud/aloe/roundtrip"
	"github.com/caicloud/aloe/runtime"
	"github.com/caicloud/aloe/types"
	"github.com/caicloud/aloe/utils/jsonutil"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	// IteratorName defines iterator variable name
	IteratorName = "iterator"
)

// Iterate set iterator variable in context variable
func Iterate(ctx *runtime.Context, iter int) {
	ctx.Variables.Set(IteratorName, jsonutil.NewStringVariable(IteratorName, "["+strconv.Itoa(iter)+"]"))
}

var (
	// defaultTimeout defines default timeout of an async task
	defaultTimeout = 1 * time.Second

	// defaultInterval defines default interval of an async task checking
	defaultInterval = 100 * time.Millisecond
)

var (
	evalFuncs = map[string]govaluate.ExpressionFunction{
		"int": func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("expected 1 arg, but received: %v", len(args))
			}
			l, err := strconv.ParseInt(args[0].(string), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("can't convert %v to int", args[0])
			}
			return float64(l), nil
		},
		"bool": func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("expected 1 arg, but received: %v", len(args))
			}
			b, err := strconv.ParseBool(args[0].(string))
			if err != nil {
				return nil, fmt.Errorf("can't convert %v to bool", args[0])
			}
			return b, nil
		},
	}
)

// iter: iter are variables defined in one flow
// all: all are variables produced by all flow
// patch: patch are variables produced by this context
// current: current are variables now
// parent: parent are variables of parent context
// iter should be merge into all
func mergeVariable(parent, current, all, iter jsonutil.VariableMap) error {
	// roundtrip should not redefine variables which have been defined in parent context
	conflict := jsonutil.IsConflict(parent, iter)
	if conflict {
		return fmt.Errorf("context define variables which have been defined in parents")
	}
	// roundtrip should not redefine variables which have been defined in previous roundtrips
	if _, err := jsonutil.Merge(all, jsonutil.ConflictOption, false, iter); err != nil {
		return err
	}
	if _, err := jsonutil.Merge(current, jsonutil.DeepOverwriteOption, false, iter); err != nil {
		return err
	}
	return nil
}

func (gf *genericFramework) constructFlow(ctx *runtime.Context, flow []types.RoundTrip) {
	flowVs := jsonutil.NewVariableMap("", nil)
	for _, rt := range flow {
		vs := gf.roundTrip(ctx, &rt)
		err := mergeVariable(ctx.Parent.Variables, ctx.Variables, flowVs, vs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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

func (gf *genericFramework) onceRoundTrip(ctx *runtime.Context, originRoundTrip *types.RoundTrip, iter int) jsonutil.VariableMap {
	if iter != -1 {
		Iterate(ctx, iter)
	}

	if originRoundTrip.When != nil {
		when, err := runtime.RenderWhen(ctx, originRoundTrip.When)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		result, err := eval(ctx, when)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		if !result {
			return nil
		}
	}

	rt, err := runtime.RenderRoundTrip(ctx, originRoundTrip)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By(fmt.Sprintf("%s: %s %s://%s%s",
		originRoundTrip.Description,
		rt.Request.Method,
		rt.Request.Scheme,
		rt.Request.Host,
		rt.Request.Path,
	))

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
	return jsonutil.NewVariableMap("", respMatcher.Variables())
}

func eval(ctx *runtime.Context, when *runtime.When) (bool, error) {
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(when.Expr, evalFuncs)
	if err != nil {
		return false, err
	}
	ps := make(map[string]interface{})
	vars := expression.Vars()
	for _, k := range vars {
		arg, ok := when.Args[k]
		if ok {
			ps[k] = arg
			continue
		}
		v, ok := ctx.Variables.Get(k)
		if ok {
			ps[k] = v.String()
			continue
		}
		// use empty as expr variable
		ps[k] = ""
	}
	result, err := expression.Evaluate(ps)
	if err != nil {
		ginkgo.By(fmt.Sprintf("failed by condition `%v`, args: %v, variables:\n%v", when.Expr, ps, ctx.Variables))
		return false, err
	}
	b, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("when condition MUST be eval as a bool")
	}
	if !b {
		ginkgo.By(fmt.Sprintf("skip by condition `%v`, args: %v", when.Expr, ps))
	}
	return b, nil

}

func (gf *genericFramework) roundTrip(ctx *runtime.Context, rt *types.RoundTrip) jsonutil.VariableMap {
	if rt.Loop == 0 {
		return gf.onceRoundTrip(ctx, rt, -1)
	}
	vars := jsonutil.NewVariableMap("", nil)
	empty := jsonutil.NewVariableMap("", nil)
	for _, def := range rt.Definitions {
		empty.Set(def.Name, nil)
	}

	for i := 0; i < rt.Loop; i++ {
		vs := gf.onceRoundTrip(ctx, rt, i)
		if vs == nil {
			vs = empty
		}
		newVars, err := jsonutil.Merge(vars, jsonutil.CombineOption, false, vs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		vars = newVars
	}
	return vars
}
