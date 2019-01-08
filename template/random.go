package template

import (
	"strconv"

	"github.com/lucasjones/reggen"
)

const (
	Random = "random"
)

const defaultLimit = 10

func random(regexp string) (string, error) {
	return reggen.Generate(regexp, defaultLimit)
}

func randomWithLimit(regexp string, limitStr string) (string, error) {
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return "", err
	}
	return reggen.Generate(regexp, limit)
}
