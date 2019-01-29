package aloe

import (
	"flag"
	"testing"

	"github.com/caicloud/aloe/cleaner"
	"github.com/caicloud/aloe/config"
	"github.com/caicloud/aloe/framework"
	"github.com/caicloud/aloe/preset"
	glogutil "github.com/caicloud/aloe/utils/glog"
)

var f framework.Framework

func assertAloeInit() {
	if f == nil {
		panic("please call aloe.Init(default *config.Config) at first")
	}
}

// Init inits framework
func Init(defaultConfig *config.Config) {
	glogutil.ChangeGlogFlag()
	c := config.Config{}
	if err := c.ParseFlags(flag.CommandLine, "aloe", defaultConfig); err != nil {
		panic(err)
	}
	f = framework.NewFramework(&c)
}

// Run will run the default framework
func Run(t *testing.T) {
	assertAloeInit()
	f.Run(t)
}

// Env sets the env of the default framework
func Env(key, value string) error {
	assertAloeInit()
	return f.Env(key, value)
}

// AppendDataDirs adds data dirs to the default framework
func AppendDataDirs(dataDirs ...string) {
	assertAloeInit()
	f.AppendDataDirs(dataDirs...)
}

// RegisterPresetter registers prestter to the default framework
func RegisterPresetter(ps ...preset.Presetter) error {
	assertAloeInit()
	return f.RegisterPresetter(ps...)
}

// RegisterCleaner registers cleaner to the default framework
func RegisterCleaner(cs ...cleaner.Cleaner) error {
	assertAloeInit()
	return f.RegisterCleaner(cs...)
}
