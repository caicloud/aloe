package template

import (
	"testing"

	"github.com/caicloud/aloe/utils/jsonutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cases := []struct {
		desc       string
		raw        string
		snippets   []string
		identitors map[int]identitor
		args       map[int][]identitor
		hasError   bool
	}{
		{
			"empty text",
			"",
			nil,
			map[int]identitor{},
			map[int][]identitor{},
			false,
		},
		{
			"normal text",
			"hello",
			[]string{"hello"},
			map[int]identitor{},
			map[int][]identitor{},
			false,
		},
		{
			"only variable text",
			"%{cluster}",
			[]string{""},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"text before variable",
			"hello%{cluster}",
			[]string{"hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"text after variable",
			"%{cluster}hello",
			[]string{"", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"two variables before text",
			"%{cluster}%{partition}hello",
			[]string{"", "", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: true,
				},
				2: identitor{
					name:  "partition",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"text with %",
			"%%{cluster}%{partition}hello",
			[]string{"%{cluster}", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "partition",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"text with function",
			"%{ cluster() }hello",
			[]string{"", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: false,
				},
			},
			map[int][]identitor{
				1: nil,
			},
			false,
		},
		{
			"text with function with two string args",
			"%{ cluster(`a`, `b`) }hello",
			[]string{"", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: false,
				},
			},
			map[int][]identitor{
				1: []identitor{
					{
						name:  "a",
						isVar: false,
					},
					{
						name:  "b",
						isVar: false,
					},
				},
			},
			false,
		},
		{
			"text with function with two variable args",
			"%{ cluster(a, b) }hello",
			[]string{"", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: false,
				},
			},
			map[int][]identitor{
				1: []identitor{
					{
						name:  "a",
						isVar: true,
					},
					{
						name:  "b",
						isVar: true,
					},
				},
			},
			false,
		},
		{
			"text with %",
			"%%%{cluster}{test}%{partition}hello",
			[]string{"%", "{test}", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster",
					isVar: true,
				},
				2: identitor{
					name:  "partition",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"text with variables with dot",
			"%{cluster.items.0}test%{partition}hello",
			[]string{"", "test", "hello"},
			map[int]identitor{
				1: identitor{
					name:  "cluster.items.0",
					isVar: true,
				},
				2: identitor{
					name:  "partition",
					isVar: true,
				},
			},
			map[int][]identitor{},
			false,
		},
		{
			"single %",
			"%",
			nil,
			nil,
			nil,
			true,
		},
		{
			"unclosed script",
			"%{",
			nil,
			nil,
			nil,
			true,
		},
		{
			"empty variable",
			"%{}",
			nil,
			nil,
			nil,
			true,
		},
	}
	for _, c := range cases {
		in, err := New(c.raw)
		if c.hasError {
			require.NotNil(t, err, "%v: err should not be nil", c.desc)
			continue
		} else {
			require.Nil(t, err, "%v: err should be nil", c.desc)
		}
		temp, ok := in.(*template)
		require.Equal(t, true, ok, "%v: Template is a *template", c.desc)
		assert.Equal(t, c.snippets, temp.snippets, "%v: snippets are not equal", c.desc)
		assert.Equal(t, c.identitors, temp.identitors, "%v: params are not equal", c.desc)
	}
}

func TestRender(t *testing.T) {
	cases := []struct {
		t        *template
		vs       map[string]jsonutil.Variable
		out      string
		hasError bool
	}{
		{
			&template{
				identitors: map[int]identitor{
					1: identitor{
						name:  "cluster",
						isVar: true,
					},
					2: identitor{
						name:  "partition",
						isVar: true,
					},
				},
				snippets: []string{
					`aaa"`,
					`"bbb`,
					`ccc`,
				},
			},
			map[string]jsonutil.Variable{
				"cluster":   jsonutil.NewStringVariable("cluster", "cid"),
				"partition": jsonutil.NewStringVariable("partition", "1.5"),
			},
			`aaa"cid"bbb1.5ccc`,
			false,
		},
		{
			&template{
				identitors: map[int]identitor{
					1: identitor{
						name:  "cluster",
						isVar: true,
					},
					2: identitor{
						name:  "partition",
						isVar: true,
					},
				},
				snippets: []string{
					`{"cluster": "`,
					`", "partition": "`,
					`"}`,
				},
			},
			map[string]jsonutil.Variable{
				"cluster":   jsonutil.NewStringVariable("cluster", "cid"),
				"partition": jsonutil.NewStringVariable("partition", "1.5"),
			},
			`{"cluster": "cid", "partition": "1.5"}`,
			false,
		},
		{
			&template{
				identitors: map[int]identitor{
					1: identitor{
						name:  "cluster.items.[0]",
						isVar: true,
					},
					2: identitor{
						name:  "partition.[0]",
						isVar: true,
					},
				},
				snippets: []string{
					`{"cluster": "`,
					`", "partition": "`,
					`"}`,
				},
			},
			map[string]jsonutil.Variable{
				"cluster": jsonutil.NewVariableMap("cluster",
					map[string]jsonutil.Variable{
						"items": jsonutil.NewVariableArray("items",
							[]jsonutil.Variable{
								jsonutil.NewStringVariable("", "cid"),
							}),
					}),
				"partition": jsonutil.NewVariableArray("partition",
					[]jsonutil.Variable{
						jsonutil.NewStringVariable("", "1.5"),
					}),
			},
			`{"cluster": "cid", "partition": "1.5"}`,
			false,
		},
	}

	for _, c := range cases {
		vm := jsonutil.NewVariableMap("", c.vs)
		out, err := c.t.Render(vm)
		if c.hasError {
			continue
		}
		assert.NoError(t, err, "render should have no error")
		assert.Equal(t, c.out, out, "render result should be same")
	}
}
