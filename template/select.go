package template

import (
	"strings"

	"github.com/caicloud/aloe/utils/jsonutil"
)

const (
	// Select defines select function
	Select = "select"
)

// selectVar select particial variable from jsonvariable
// if field is not exists, return empty string
func selectVar(v jsonutil.Variable, selector string) (string, error) {
	if v == nil {
		return "", nil
	}
	selectors := strings.Split(selector, ",")
	res, err := v.Select(selectors...)
	if err != nil {
		return "", nil
	}
	if res == nil {
		return "", nil
	}
	return res.String(), nil
}
