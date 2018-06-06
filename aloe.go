package aloe

import (
	"testing"

	"github.com/caicloud/aloe/cleaner"
	"github.com/caicloud/aloe/framework"
	"github.com/caicloud/aloe/preset"
)

var f = framework.NewFramework()

// Run will run the default framework
func Run(t *testing.T) {
	f.Run(t)
}

// Env sets the env of the default framework
func Env(key, value string) error {
	return f.Env(key, value)
}

// AppendDataDirs adds data dirs to the default framework
func AppendDataDirs(dataDirs ...string) {
	f.AppendDataDirs(dataDirs...)
}

// RegisterPresetter registers prestter to the default framework
func RegisterPresetter(ps ...preset.Presetter) error {
	return f.RegisterPresetter(ps...)
}

// RegisterCleaner registers cleaner to the default framework
func RegisterCleaner(cs ...cleaner.Cleaner) error {
	return f.RegisterCleaner(cs...)
}
