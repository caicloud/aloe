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
		{
			description: "context has a tuple with a validator and a constructor",
			c: &types.Context{
				ValidatedFlow: []types.RoundTripTuple{
					{
						Validator: []types.RoundTrip{
							types.RoundTrip{},
						},
						Constructor: []types.RoundTrip{
							types.RoundTrip{},
						},
					},
				},
			},
			expected: nil,
		},
		{
			description: "context has a tuple without a validator",
			c: &types.Context{
				ValidatedFlow: []types.RoundTripTuple{
					{
						Constructor: []types.RoundTrip{
							types.RoundTrip{},
						},
					},
				},
			},
			expected: ErrorList{
				fmt.Errorf("constructor and validator should not be empty"),
			},
		},
		{
			description: "context has a tuple without a constructor",
			c: &types.Context{
				ValidatedFlow: []types.RoundTripTuple{
					{
						Validator: []types.RoundTrip{
							types.RoundTrip{},
						},
					},
				},
			},
			expected: ErrorList{
				fmt.Errorf("constructor and validator should not be empty"),
			},
		},
		{
			description: "context has no validated flow tuple",
			c: &types.Context{
				ValidatedFlow: []types.RoundTripTuple{},
			},
			expected: nil,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, ValidateContext(c.c), c.description)
	}
}
