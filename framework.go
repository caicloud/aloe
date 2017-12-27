package framework

import (
	"net/http"
	"time"

	"github.com/caicloud/aloe/data"
	"github.com/caicloud/aloe/roundtrip"
	"github.com/caicloud/aloe/template"
	"github.com/caicloud/aloe/types"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// Framework defines an API test framework
type Framework interface {
	Run() error
}

// ClearFn defines function to clear context
type ClearFn func()

// NewFramework returns an API test framework
func NewFramework(host string, clearFn ClearFn, dataDirs ...string) Framework {
	return &genericFramework{
		dataDirs,
		roundtrip.NewClient(host),
		clearFn,
	}
}

type genericFramework struct {
	dataDirs []string

	client *roundtrip.Client

	clearFn ClearFn
}

func (gf *genericFramework) Run() error {
	for _, r := range gf.dataDirs {
		dir, err := data.Walk(r)
		if err != nil {
			return err
		}
		ctx := &types.Context{}
		f := gf.walk(ctx, dir, true)
		ginkgo.Describe(dir.Context.Summary, f)
	}

	return nil
}

func (gf *genericFramework) walk(ctx *types.Context, dir *data.Dir, isTop bool) func() {
	dirs, files := dir.Dirs, dir.Files
	ctxConfig := dir.Context

	return func() {
		var contextVs map[string]template.Variable

		ginkgo.BeforeEach(func() {
			contextVs = ctx.Variables
			ctx.Variables, ctx.Error = gf.constructContext(ctx, &ctxConfig)
		})

		ginkgo.AfterEach(func() {
			if isTop {
				gf.clearFn()
				ctx.Variables = nil
			} else {
				ctx.Variables = contextVs
			}
			ctx.Error = nil
		})

		for name, d := range dirs {
			f := gf.walk(ctx, &d, false)
			summary := genSummary(name, d.Context.Summary)
			ginkgo.Context(summary, f)
		}
		for name, c := range files {
			summary := genSummary(name, c.Case.Description)
			f := gf.itFunc(ctx, &c)
			ginkgo.It(summary, f)
		}
	}
}

func genSummary(name, summary string) string {
	return name + ": " + summary
}

var (
	defaultTimeout  = 1 * time.Second
	defaultInterval = 100 * time.Millisecond
)

func (gf *genericFramework) itFunc(ctx *types.Context, file *data.File) func() {
	c := file.Case
	return func() {
		ginkgo.By("Context should be constructed successfully")
		gomega.Expect(ctx.Error).NotTo(gomega.HaveOccurred())

		for _, rt := range c.Flow {
			ginkgo.By(rt.Description)

			respMatcher, err := roundtrip.MatchResponse(ctx, &rt)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			if ev := rt.Response.Eventually; ev != nil {
				timeout := ev.Timeout
				if timeout == nil {
					timeout = &types.Duration{defaultTimeout}
				}
				interval := ev.Interval
				if interval == nil {
					interval = &types.Duration{defaultInterval}
				}
				gomega.Eventually(func() *http.Response {
					resp, err := gf.client.DoRequest(ctx, &rt)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					return resp
				}, ev.Timeout.Duration, ev.Interval.Duration).Should(respMatcher)

			} else {
				resp, err := gf.client.DoRequest(ctx, &rt)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).To(respMatcher)
			}
			vs, err := respMatcher.Variables()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			for k, v := range vs {
				ctx.Variables[k] = v
			}
		}
	}

}
