package template

import (
	"github.com/lucasjones/reggen"
)

const (
	// Random defines random function
	Random = "random"
)

const defaultLimit = 10

func random(regexp string) (string, error) {
	return reggen.Generate(regexp, defaultLimit)
}

func randomWithLimit(regexp string, limit int) (string, error) {
	return reggen.Generate(regexp, limit)
}
