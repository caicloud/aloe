package cleaner

import "github.com/caicloud/aloe/utils/jsonutil"

// Cleaner defines custom cleaners
type Cleaner interface {
	// Name defines cleaner name
	Name() string

	// Clean will clean the variables
	Clean(variables map[string]jsonutil.Variable) error

	// ForceClean will be called if function clean failed
	// If ForceClean returns false, test case will fail
	// ForceClean(variables map[string]jsonutil.Variable) bool
}
