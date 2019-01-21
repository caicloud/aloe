package config

import (
	"flag"
	"fmt"
)

// Config defines config of test
type Config struct {
	Focus string
	Skip  string
}

func withPrefix(prefix, flagName string) string {
	if prefix == "" {
		return flagName
	}
	return prefix + "." + flagName
}

// ParseFlags init flags
func (c *Config) ParseFlags(flagSet *flag.FlagSet, prefix string, defaults *Config) error {
	if flagSet.Parsed() {
		return fmt.Errorf("flags have been parsed")
	}
	if defaults == nil {
		defaults = &Config{}
	}
	flagSet.StringVar(&c.Focus,
		withPrefix(prefix, "focus"),
		defaults.Focus,
		`only run cases with specified labels. Labels should be splited by comma. e.g. "aaa,bbb" means only cases with "aaa" or "bbb" labels will be run`)

	flagSet.StringVar(&c.Skip,
		withPrefix(prefix, "skip"),
		defaults.Skip,
		`skip cases with specified labels. Labels should be splited by comma. e.g. "aaa,bbb" means skip cases with "aaa" or "bbb" labels`)
	return nil
}
