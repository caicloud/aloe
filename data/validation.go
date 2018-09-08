package data

import (
	"fmt"
	"strings"

	"github.com/caicloud/aloe/types"
)

// ErrorList defines list of error
type ErrorList []error

// Error implements error interface
func (el ErrorList) Error() string {
	ss := make([]string, 0, len(el))
	for _, e := range el {
		ss = append(ss, e.Error())
	}
	return strings.Join(ss, "\n")
}

// ValidateCase will validate case before case is running
func ValidateCase(c *types.Case) error {
	if c == nil {
		return fmt.Errorf("case is empty")
	}
	return nil
}

// ValidateContext will validate context before case is running
func ValidateContext(c *types.Context) error {
	if c == nil {
		return fmt.Errorf("context is empty")
	}
	errList := ErrorList{}
	for _, rt := range c.ValidatedFlow {
		if len(rt.Constructor) == 0 || len(rt.Validator) == 0 {
			errList = append(errList, fmt.Errorf("constructor and validator should not be empty"))
		}
	}
	if len(errList) != 0 {
		return errList
	}
	return nil
}
