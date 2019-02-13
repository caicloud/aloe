package template

import (
	"fmt"
	"strconv"

	"github.com/caicloud/aloe/utils/jsonutil"
)

const (
	// Length defines len function
	Length = "len"
)

func length(v jsonutil.Variable) (string, error) {
	m, ok := v.(jsonutil.Measurable)
	if !ok {
		return "", fmt.Errorf("%v is not measurable", v)
	}
	l := m.Len()
	if l == -1 {
		return "", fmt.Errorf("can't get %v length", v)
	}
	return strconv.Itoa(l), nil
}
