package matcher

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		desc    string
		matcher []byte
		body    []byte
		res     bool
	}{
		{
			"normal case",
			[]byte(`{
				"string": "string",
				"int": 1,
				"float": 1.1,
				"null": null,
				"array": ["xx"],
				"obj": {"yy": "yy"}
			}`),
			[]byte(`{
				"string": "string",
				"int": 1,
				"float": 1.1,
				"null": null,
				"array": ["xx"],
				"obj": {"yy": "yy"}
			}`),
			true,
		},
		{
			"normal case -- match partial data",
			[]byte(`{
				"string": "string",
				"int": 1,
				"float": 1.1
			}`),
			[]byte(`{
				"string": "string",
				"int": 1,
				"float": 1.1,
				"null": null,
				"array": ["xx"],
				"obj": {"yy": "yy"}
			}`),
			true,
		},
		{
			"normal case -- match data which doesn't exist",
			[]byte(`{
				"xxx": "string"
			}`),
			[]byte(`{
				"string": "string",
				"int": 1,
				"float": 1.1,
				"null": null,
				"array": ["xx"],
				"obj": {"yy": "yy"}
			}`),
			false,
		},
		{
			"$regexp case",
			[]byte(`{
				"string": {
					"$regexp": "[a-z]*"
				}
			}`),
			[]byte(`{
				"string": "string"
			}`),
			true,
		},
		{
			"array case -- with $regexp",
			[]byte(`{
				"array": [
					"aaa",
					{
						"$regexp": "[a-z]*"
					},
					"ccc"
				]
			}`),
			[]byte(`{
				"array": [
					"aaa",
					"bbb",
					"ccc"
				]
			}`),
			true,
		},
		{
			"$exists case -- $exists is true",
			[]byte(`{
				"string": {
					"$exists": true
				}
			}`),
			[]byte(`{
			}`),
			false,
		},
		{
			"$exists case -- $exists is false",
			[]byte(`{
				"string": {
					"$exists": false
				}
			}`),
			[]byte(`{
				"string": "string"
			}`),
			false,
		},
		{
			"$exists case -- $exists with $regexp",
			[]byte(`{
				"string": {
					"$exists": true,
					"$regexp": "[0-9]*"
				}
			}`),
			[]byte(`{
				"string": "123"
			}`),
			true,
		},
		{
			"$exists case -- $exists with $regexp and missing field",
			[]byte(`{
				"string": {
					"$exists": true,
					"$regexp": "[0-9]*"
				}
			}`),
			[]byte(`{
			}`),
			false,
		},
		{
			"$match case -- $match with $exists",
			[]byte(`{
				"string": {
					"$exists": true,
					"$match": "123"
				}
			}`),
			[]byte(`{
				"string": "123"
			}`),
			true,
		},
		{
			"$match case -- $match with $exists(another)",
			[]byte(`{
				"array": {
					"$exists": true,
					"$match": [
						"aaa"
					]
				}
			}`),
			[]byte(`{
				"array": [
					"aaa"
				]
			}`),
			true,
		},
		{
			"$len case -- string",
			[]byte(`{
				"string": {
					"$len": 4
				}
			}`),
			[]byte(`{
				"string": "1234"
			}`),
			true,
		},
		{
			"$len case -- object",
			[]byte(`{
				"obj": {
					"$len": 2
				}
			}`),
			[]byte(`{
				"obj": {
					"aaa": "aaa",
					"bbb": "bbb"
				}
			}`),
			true,
		},
		{
			"$len case -- array",
			[]byte(`{
				"array": {
					"$len": 2
				}
			}`),
			[]byte(`{
				"array": [
					"aaa",
					"bbb"
				]
			}`),
			true,
		},
	}

	for _, c := range cases {
		m, err := Parse(c.matcher)
		require.NoError(t, err, c.desc)
		var b interface{}
		require.NoError(t, json.Unmarshal(c.body, &b), c.desc)
		res, err := m.Match(b)
		require.NoError(t, err, c.desc)
		assert.Equal(t, c.res, res, c.desc)
	}
}
