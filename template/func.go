package template

import "fmt"

func Call(name string, args ...string) (string, error) {
	switch name {
	case Random:
		if len(args) == 1 {
			return random(args[0])
		} else if len(args) == 2 {
			return randomWithLimit(args[0], args[1])
		} else {
			return "", fmt.Errorf("func random expected 1 or 2 args, but received: %v", len(args))
		}
	default:
		return "", fmt.Errorf("unknown function named %v", name)
	}
}
