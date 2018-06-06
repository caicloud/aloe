package framework

import (
	"fmt"
	"testing"

	"github.com/caicloud/aloe/cleaner"
	"github.com/caicloud/aloe/data"
	"github.com/caicloud/aloe/preset"
	"github.com/caicloud/aloe/roundtrip"
	"github.com/caicloud/aloe/runtime"
	"github.com/caicloud/aloe/utils/jsonutil"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// Framework defines an API test framework
type Framework interface {
	// Env sets the envirment context
	Env(key, value string) error

	// AppendDataDirs add data into framework
	AppendDataDirs(dataDirs ...string)

	// RegisterCleaner registers cleaner of framework
	RegisterCleaner(cs ...cleaner.Cleaner) error

	// RegisterPresetter registers presetter of framework
	RegisterPresetter(ps ...preset.Presetter) error

	// Run will run the framework
	Run(t *testing.T)
}

// NewFramework returns an API test framework
func NewFramework(dataDirs ...string) Framework {
	reqHeader := preset.NewHeaderPresetter(preset.RequestType)
	respHeader := preset.NewHeaderPresetter(preset.ResponseType)
	host := preset.NewHostPresetter()
	return &genericFramework{
		dataDirs: dataDirs,
		client:   roundtrip.NewClient(),
		cleaners: map[string]cleaner.Cleaner{},
		presetters: map[string]preset.Presetter{
			reqHeader.Name():  reqHeader,
			respHeader.Name(): respHeader,
			host.Name():       host,
		},
		adam: &runtime.Context{},
	}
}

type genericFramework struct {
	dataDirs []string

	client *roundtrip.Client

	cleaners map[string]cleaner.Cleaner

	presetters map[string]preset.Presetter

	adam *runtime.Context
}

// Env implements Framework interface
func (gf *genericFramework) Env(key, value string) error {
	if gf.adam.Variables == nil {
		gf.adam.Variables = map[string]jsonutil.Variable{}
	}
	if _, ok := gf.adam.Variables[key]; ok {
		return fmt.Errorf("%v has been defined", key)
	}
	gf.adam.Variables[key] = jsonutil.NewVariable(key, value)
	return nil
}

// AppendDataDirs implements Framework interface
func (gf *genericFramework) AppendDataDirs(ds ...string) {
	gf.dataDirs = append(gf.dataDirs, ds...)
}

// RegisterCleaner implements Framework interface
func (gf *genericFramework) RegisterCleaner(cs ...cleaner.Cleaner) error {
	for _, c := range cs {
		if _, ok := gf.cleaners[c.Name()]; ok {
			return fmt.Errorf("can't register cleaner %v: already exists", c.Name())
		}
		gf.cleaners[c.Name()] = c
	}
	return nil
}

// RegisterPresetter implements Framework interface
func (gf *genericFramework) RegisterPresetter(ps ...preset.Presetter) error {
	for _, p := range ps {
		if _, ok := gf.presetters[p.Name()]; ok {
			return fmt.Errorf("can't register presetter %v: already exists", p.Name())
		}
		gf.presetters[p.Name()] = p
	}
	return nil
}

func (gf *genericFramework) Run(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	for _, r := range gf.dataDirs {
		dir, err := data.Walk(r)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		f := gf.walk(gf.adam, dir)
		ginkgo.Describe(dir.Context.Summary, f)
	}
	ginkgo.RunSpecs(t, "Test Suit")
}

func (gf *genericFramework) walk(parent *runtime.Context, dir *data.Dir) func() {
	dirs, files := dir.Dirs, dir.Files
	ctxConfig := dir.Context
	total := dir.CaseNum

	return func() {
		// TODO(liubog2008): need to support concurrency
		count := 0

		var ctx runtime.Context

		for name, d := range dirs {
			f := gf.walk(&ctx, &d)
			summary := genSummary(name, d.Context.Summary)
			ginkgo.Context(summary, f)
		}
		for name, c := range files {
			summary := genSummary(name, c.Case.Summary)
			f := gf.itFunc(&ctx, &c)
			ginkgo.It(summary, f)
		}

		ginkgo.BeforeEach(func() {
			// inherit parent context
			runtime.CopyContext(&ctx, parent)

			// render preset config
			gomega.Expect(runtime.RenderPresetters(&ctx, ctxConfig.Presetters)).
				NotTo(gomega.HaveOccurred())

			gomega.Expect(gf.constructRoundTripTemplate(&ctx)).
				NotTo(gomega.HaveOccurred())

			if count == 0 {
				gf.constructFlow(&ctx, ctxConfig.Flow)
			}
			gf.constructValidatedFlow(&ctx, ctxConfig.ValidatedFlow)

		})

		ginkgo.AfterEach(func() {
			// render cleaner config
			gomega.Expect(runtime.RenderCleaners(&ctx, ctxConfig.Cleaners)).
				NotTo(gomega.HaveOccurred())

			count++
			for _, c := range ctx.Cleaners {
				if !c.ForEach && count != total {
					continue
				}
				cleaner, ok := gf.cleaners[c.Name]
				gomega.Expect(ok).To(gomega.BeTrue())
				gomega.Expect(cleaner.Clean(ctx.RoundTripTemplate, c.Args)).
					NotTo(gomega.HaveOccurred())
			}
		})
	}
}

func genSummary(name, summary string) string {
	return name + ": " + summary
}

func (gf *genericFramework) itFunc(ctx *runtime.Context, file *data.File) func() {
	c := file.Case
	return func() {
		for _, rt := range c.Flow {
			gf.roundTrip(ctx, &rt, true)
		}
	}
}
