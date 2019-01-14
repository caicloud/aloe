package template

const (
	Exist = "exist"
)

func isExist(arg Argument) (string, error) {
	if arg.IsNil() {
		return "false", nil
	}
	return "true", nil
}
