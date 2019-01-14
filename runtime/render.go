package runtime

import (
	"fmt"
	"strings"

	"github.com/caicloud/aloe/types"
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

func splitMethodAndPath(api string) (string, string) {
	s := strings.SplitN(api, " ", 2)
	return strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
}

func renderRequest(ctx *Context, runtimereq *Request, req *types.Request) error {
	if req.Host != "" {
		runtimereq.Host = req.Host
	}
	if req.Scheme != "" {
		runtimereq.Scheme = req.Scheme
	}
	if req.API != nil {
		api, err := req.API.Render(ctx.Variables)
		if err != nil {
			return err
		}
		method, path := splitMethodAndPath(api)
		runtimereq.Method = method
		runtimereq.Path = path
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

func renderDefinition(ctx *Context, dcs []types.Definition) ([]Definition, error) {
	if len(dcs) == 0 {
		return nil, nil
	}
	ds := make([]Definition, 0, len(dcs))

	for _, dc := range dcs {
		d := Definition{
			Name: dc.Name,
		}
		switch dc.Type {
		case "body", "header", "status":
			d.Type = DefinitionType(dc.Type)
		case "":
			d.Type = BodyType
		default:
			return nil, fmt.Errorf("can't understand definition type %v: only [body, header, status] is allowed", dc.Type)
		}
		for _, st := range dc.Selector {
			s, err := st.Render(ctx.Variables)
			if err != nil {
				return nil, err
			}
			d.Selector = append(d.Selector, s)
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
