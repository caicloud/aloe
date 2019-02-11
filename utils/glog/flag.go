package glog

import (
	"flag"
	"os"
	"strings"
)

// ChangeGlogFlag do three things for glog flag
// 1. add 'glog.' prefix for all glog flag
// 2. change '_' to '-'
// 3. change default value of logtostderr to 'true'
func ChangeGlogFlag() {
	newFlagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.VisitAll(func(f *flag.Flag) {
		newFlag := f
		fn, ok := flagMap[f.Name]
		if ok {
			newFlag = &flag.Flag{
				Name:     f.Name,
				DefValue: f.DefValue,
				Usage:    f.Usage,
				Value:    f.Value,
			}
			fn(newFlag)
		}
		newFlagSet.Var(newFlag.Value, newFlag.Name, newFlag.Usage)
	})
	flag.CommandLine = newFlagSet
}

var flagMap = map[string]func(*flag.Flag){
	"logtostderr":      combine(changeDefaultValue("true"), normalizeName),
	"alsologtostderr":  normalizeName,
	"log_dir":          normalizeName,
	"stderrthreshold":  normalizeName,
	"log_backtrace_at": normalizeName,
	"v":                normalizeName,
	"vmodule":          normalizeName,
}

func normalizeName(f *flag.Flag) {
	name := "glog." + strings.Replace(f.Name, "_", "-", -1)
	f.Name = name
}

func changeDefaultValue(defValue string) func(*flag.Flag) {
	return func(f *flag.Flag) {
		f.DefValue = defValue
		if err := f.Value.Set(defValue); err != nil {
			panic(err)
		}
	}
}

func combine(fns ...func(*flag.Flag)) func(*flag.Flag) {
	return func(f *flag.Flag) {
		for _, fn := range fns {
			fn(f)
		}
	}
}
