package cleaner

import "github.com/caicloud/aloe/runtime"

// Cleaner defines custom cleaner which can clean up context
// It is a callback hook called when all cases in context are finished
// Cleaner can be used to delete objects which are created in the context
// or even drop database directly
type Cleaner interface {
	// Name defines cleaner name
	Name() string

	// Clean will be called after all of the cases in the context are
	// finished
	Clean(template *runtime.RoundTripTemplate, args map[string]string) error
}
