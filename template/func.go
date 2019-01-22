package template

import (
	"fmt"
	"strconv"

	"github.com/caicloud/aloe/utils/jsonutil"
)

// Call calls function of template
func Call(name string, args ...jsonutil.Variable) (string, error) {
	switch name {
	case Random:
		if len(args) == 0 {
			return "", fmt.Errorf("func random expected 1 or 2 args, but received: %v", len(args))
		}
		if args[0] == nil {
			return "", fmt.Errorf("first argument of random is nil")
		}
		regexp := args[0].String()
		if len(args) == 1 {
			return random(regexp)
		}
		if args[1] == nil {
			return "", fmt.Errorf("second argument of random is nil")
		}
		limitStr := args[1].String()
		limit, err := strconv.ParseInt(limitStr, 32, 10)
		if err != nil {
			return "", err
		}
		if len(args) == 2 {
			return randomWithLimit(regexp, int(limit))
		}
		return "", fmt.Errorf("func random expected 1 or 2 args, but received: %v", len(args))
	case Exist:
		if len(args) == 1 {
			return isExist(args[0])
		}
		if len(args) == 2 {
			if args[1] == nil {
				return "", fmt.Errorf("second argument of exist is nil")
			}
			return isExistWithSelector(args[0], args[1].String())
		}
		return "", fmt.Errorf("func exist expected 1 or 2 arg, but received: %v", len(args))

	case Select:
		if len(args) != 2 {
			return "", fmt.Errorf("func select expected 2 arg, but received: %v", len(args))
		}
		if args[1] == nil {
			return "", fmt.Errorf("second argument of select is nil")
		}
		return selectVar(args[0], args[1].String())
	default:
		return "", fmt.Errorf("unknown function named %v", name)
	}
}
