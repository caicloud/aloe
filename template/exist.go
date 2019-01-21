package template

const (
	// Exist defines exist function
	// It checks whether argument exists
	Exist = "exist"
)

func isExist(arg Argument) (string, error) {
	if arg.IsNil() {
		return "false", nil
	}
	return "true", nil
}
