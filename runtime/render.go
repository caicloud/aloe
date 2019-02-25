package runtime

import (
	"fmt"
	"strings"

	"github.com/caicloud/aloe/types"
	"github.com/caicloud/aloe/utils/jsonutil"
)

// RenderPresetters render preset config with current context
// and save into context
func RenderPresetters(ctx *Context, pcs []types.PresetConfig) error {
	ps := make([]Presetter, 0, len(pcs))

	for _, pc := range pcs {
		p := Presetter{
			Name: pc.Name,
		}
		if len(pc.Args) != 0 {
			p.Args = map[string]string{}
		}

		for k, v := range pc.Args {
			rendered, err := v.Render(ctx.Variables)
			if err != nil {
				return err
			}
			p.Args[k] = rendered
		}
		ps = append(ps, p)
	}
	ctx.Presetters = ps
	return nil
}

// RenderCleaners render cleaner config with current context
// and save into context
func RenderCleaners(ctx *Context, ccs []types.CleanerConfig) error {
	cs := make([]Cleaner, 0, len(ccs))

	for _, cc := range ccs {
		c := Cleaner{
			Name:    cc.Name,
			ForEach: cc.ForEach,
		}
		if len(cc.Args) != 0 {
			c.Args = map[string]string{}
		}

		for k, v := range cc.Args {
			rendered, err := v.Render(ctx.Variables)
			if err != nil {
				return err
			}
			c.Args[k] = rendered
		}
		cs = append(cs, c)
	}
	ctx.Cleaners = cs
	return nil
}

// RenderRoundTrip render template in round trip config with current context
func RenderRoundTrip(ctx *Context, rtc *types.RoundTrip) (*RoundTrip, error) {
	cp := CopyRoundTripTemplate(ctx.RoundTripTemplate)
	if cp == nil {
		cp = &RoundTripTemplate{}
	}
	rt := &RoundTrip{
		RoundTripTemplate: *cp,
	}
	if rtc.Client != "" {
		rt.Client = rtc.Client
	}
	if err := renderRequest(ctx, &rt.Request, &rtc.Request); err != nil {
		return nil, err
	}
	if err := renderResponse(ctx, &rt.Response, &rtc.Response); err != nil {
		return nil, err
	}
	ds, err := renderDefinition(ctx, rtc.Definitions)
	if err != nil {
		return nil, err
	}
	rt.Definitions = ds

	return rt, nil
}

// RenderExports render export variables
func RenderExports(ctx *Context, exports []types.Var) error {
	vs := jsonutil.NewVariableMap("", nil)
	for _, exportConf := range exports {
		export := Var{}
		if err := renderVar(ctx, &exportConf, &export); err != nil {
			return err
		}
		v, err := ctx.Variables.Select(export.Selector...)
		if err != nil {
			return fmt.Errorf("can't export var %v: %v", export.Name, err)
		}
		if _, ok := vs.Get(export.Name); ok {
			return fmt.Errorf("can't export var %v twice", export.Name)
		}
		vs.Set(export.Name, v)
	}

	newExports, err := jsonutil.Merge(ctx.Exports, jsonutil.OverwriteOption, false, vs)
	if err != nil {
		return err
	}
	ctx.Exports = newExports
	// reconstruct ctx variable
	newVs, err := jsonutil.Merge(ctx.Parent.Variables, jsonutil.ConflictOption, true, ctx.Exports)
	if err != nil {
		return err
	}
	ctx.Variables = newVs
	return nil
}

func splitMethodAndPath(api string) (string, string) {
	s := strings.SplitN(api, " ", 2)
	return strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
}

func isAbs(path string) bool {
	if len(path) < 2 {
		return false
	}
	if path[0] == '`' && path[len(path)-1] == '`' {
		return true
	}
	return false
}

func renderRequest(ctx *Context, runtimereq *Request, req *types.Request) error {
	if req.Host != nil {
		host, err := req.Host.Render(ctx.Variables)
		if err != nil {
			return err
		}
		runtimereq.Host = host
	}
	if req.Scheme != nil {
		scheme, err := req.Scheme.Render(ctx.Variables)
		if err != nil {
			return err
		}
		runtimereq.Scheme = scheme
	}
	if req.API != nil {
		api, err := req.API.Render(ctx.Variables)
		if err != nil {
			return err
		}
		method, path := splitMethodAndPath(api)
		runtimereq.Method = method
		if isAbs(path) {
			runtimereq.Path = path[1 : len(path)-1]
		} else if runtimereq.PathTemplate != "" {
			runtimereq.Path = fmt.Sprintf(runtimereq.PathTemplate, path)
		} else {
			runtimereq.Path = path
		}
	}
	if err := renderHeader(ctx, runtimereq.Headers, req.Headers); err != nil {
		return err
	}
	if req.Body != nil {
		body, err := req.Body.Render(ctx.Variables)
		if err != nil {
			return err
		}
		runtimereq.Body = []byte(body)
	}
	return nil
}

func renderResponse(ctx *Context, resp *Response, respConf *types.Response) error {
	if respConf.StatusCode != 0 {
		resp.StatusCode = respConf.StatusCode
	}

	if err := renderHeader(ctx, resp.Headers, respConf.Headers); err != nil {
		return err
	}

	if respConf.Body != nil {
		body, err := respConf.Body.Render(ctx.Variables)
		if err != nil {
			return err
		}
		resp.Body = []byte(body)
	}

	if respConf.Eventually != nil {
		resp.Async = true
		if respConf.Eventually.Interval != nil {
			resp.Interval = respConf.Eventually.Interval.Duration
		}
		if respConf.Eventually.Timeout != nil {
			resp.Timeout = respConf.Eventually.Timeout.Duration
		}
	}

	return nil
}

func renderVar(ctx *Context, vc *types.Var, v *Var) error {
	v.Name = vc.Name
	for _, st := range vc.Selector {
		s, err := st.Render(ctx.Variables)
		if err != nil {
			return err
		}
		v.Selector = append(v.Selector, s)
	}
	return nil

}

func renderDefinition(ctx *Context, dcs []types.Definition) ([]Definition, error) {
	if len(dcs) == 0 {
		return nil, nil
	}
	ds := make([]Definition, 0, len(dcs))

	for _, dc := range dcs {
		d := Definition{}
		switch dc.Type {
		case "body", "header", "status":
			d.Type = DefinitionType(dc.Type)
		case "":
			d.Type = BodyType
		default:
			return nil, fmt.Errorf("can't understand definition type %v: only [body, header, status] is allowed", dc.Type)
		}
		if err := renderVar(ctx, &dc.Var, &d.Var); err != nil {
			return nil, fmt.Errorf("can't render var in definition: %v", err)
		}

		ds = append(ds, d)
	}
	return ds, nil
}

func renderHeader(ctx *Context, current map[string]string, headers map[string]types.Template) error {
	if current == nil {
		current = map[string]string{}
	}
	for k, v := range headers {
		value, err := v.Render(ctx.Variables)
		if err != nil {
			return err
		}
		current[k] = value
	}
	return nil
}

// RenderWhen will render condition config into runtime condition
func RenderWhen(ctx *Context, whenConf *types.When) (*When, error) {
	if whenConf == nil {
		return nil, nil
	}
	when := When{
		Expr: whenConf.Expr,
		Args: map[string]string{},
	}
	for k, v := range whenConf.Args {
		arg, err := v.Render(ctx.Variables)
		if err != nil {
			return nil, err
		}
		when.Args[k] = arg
	}
	return &when, nil
}
