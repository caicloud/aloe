package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cases := []struct {
		raw      string
		snippts  []string
		varNames []string
		hasError bool
	}{
		{
			"",
			[]string{""},
			nil,
			false,
		},
		{
			"hello",
			[]string{"hello"},
			nil,
			false,
		},
		{
			"%{cluster}",
			[]string{"", ""},
			[]string{"cluster"},
			false,
		},
		{
			"hello%{cluster}",
			[]string{"hello", ""},
			[]string{"cluster"},
			false,
		},
		{
			"%{cluster}hello",
			[]string{"", "hello"},
			[]string{"cluster"},
			false,
		},
		{
			"%{cluster}%{partition}hello",
			[]string{"", "", "hello"},
			[]string{"cluster", "partition"},
			false,
		},
		{
			"%%{cluster}%{partition}hello",
			[]string{"%{cluster}", "hello"},
			[]string{"partition"},
			false,
		},
		{
			"%%%{cluster}{test}%{partition}hello",
			[]string{"%", "{test}", "hello"},
			[]string{"cluster", "partition"},
			false,
		},
		{
			"%",
			nil,
			nil,
			true,
		},
		{
			"%{",
			nil,
			nil,
			true,
		},
		{
			"%{}",
			nil,
			nil,
			true,
		},
	}
	for _, c := range cases {
		in, err := New(c.raw)
		if c.hasError {
			assert.NotNil(t, err, "err should not be nil")
			continue
		}
		temp, ok := in.(*template)
		assert.Equal(t, true, ok, "Template is a *template")
		assert.Equal(t, c.snippts, temp.snippts, "snippts are not equal")
		assert.Equal(t, c.varNames, temp.varNames, "params are not equal")
	}
}

func TestRender(t *testing.T) {
	cases := []struct {
		t        *template
		vs       map[string]Variable
		out      string
		hasError bool
	}{
		{
			&template{
				[]string{"cluster", "partition"},
				[]string{
					`aaa"`,
					`"bbb`,
					`ccc`,
				},
			},
			map[string]Variable{
				"cluster": {
					Raw:  []byte("cid"),
					Name: "cluster",
					Type: StringType,
				},
				"partition": {
					Raw:  []byte("1.5"),
					Name: "partition",
					Type: NumberType,
				},
			},
			`aaa"cid"bbb1.5ccc`,
			false,
		},
		{
			&template{
				[]string{"cluster", "partition"},
				[]string{
					`{"cluster": "`,
					`", "partition": "`,
					`"}`,
				},
			},
			map[string]Variable{
				"cluster": {
					Raw:  []byte("cid"),
					Name: "cluster",
					Type: StringType,
				},
				"partition": {
					Raw:  []byte("1.5"),
					Name: "partition",
					Type: NumberType,
				},
			},
			`{"cluster": "cid", "partition": "1.5"}`,
			false,
		},
	}

	for _, c := range cases {
		out, err := c.t.Render(c.vs)
		if c.hasError {
			continue
		}
		assert.NoError(t, err, "render should have no error")
		assert.Equal(t, c.out, out, "render result should be same")
	}
}
