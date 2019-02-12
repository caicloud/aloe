package data

import (
	"fmt"
	"testing"

	"github.com/caicloud/aloe/types"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	cases := []struct {
		description string
		errs        ErrorList
		expected    string
	}{
		{
			description: "empty error list",
			errs:        ErrorList{},
			expected:    "",
		},
		{
			description: "single error in list",
			errs: ErrorList{
				fmt.Errorf("aaa"),
			},
			expected: "aaa",
		},
		{
			description: "multiple errors in list",
			errs: ErrorList{
				fmt.Errorf("aaa"),
				fmt.Errorf("bbb"),
			},
			expected: "aaa\nbbb",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.errs.Error(), c.description)
	}
}

func TestValidateCase(t *testing.T) {
	cases := []struct {
		description string
		c           *types.Case
		expected    error
	}{
		{
			description: "empty case",
			c:           nil,
			expected:    fmt.Errorf("case is empty"),
		},
		{
			description: "not empty case",
			c:           &types.Case{},
			expected:    nil,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, ValidateCase(c.c), c.description)
	}
}

func TestValidateContext(t *testing.T) {
	cases := []struct {
		description string
		c           *types.Context
		expected    error
	}{
		{
			description: "empty context",
			c:           nil,
			expected:    fmt.Errorf("context is empty"),
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, ValidateContext(c.c), c.description)
	}
}
