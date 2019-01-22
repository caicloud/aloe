package template

import (
	"strings"

	"github.com/caicloud/aloe/utils/jsonutil"
)

const (
	// Exist defines exist function
	// It checks whether argument exists
	Exist = "exist"
)

const (
	// True is true string
	True = "true"
	// False is false string
	False = "false"
)

func isExist(arg jsonutil.Variable) (string, error) {
	if arg == nil {
		return False, nil
	}
	return True, nil
}

func isExistWithSelector(arg jsonutil.Variable, selector string) (string, error) {
	if arg == nil {
		return False, nil
	}
	selectors := strings.Split(selector, ",")
	_, err := arg.Select(selectors...)
	if err != nil {
		return False, nil
	}
	return True, nil
}
