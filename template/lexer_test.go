package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	cases := []struct {
		desc     string
		isNormal bool
		buf      []rune

		expectedToken     []rune
		expectedTokenType Token
		expectedErr       error
		expectedOffset    int
	}{
		{
			"empty text",
			true,
			[]rune(``),
			[]rune(``),
			TextToken,
			nil,
			0,
		},
		{
			"normal text case",
			true,
			[]rune(`this is a normal text`),
			[]rune(`this is a normal text`),
			TextToken,
			nil,
			21,
		},
		{
			"normal variable script case",
			false,
			[]rune(` normalvariable }`),
			[]rune(`normalvariable`),
			VariableNameToken,
			nil,
			17,
		},
		{
			"normal func script case",
			false,
			[]rune(` random() }`),
			[]rune(`random`),
			FuncNameToken,
			nil,
			7,
		},
		{
			"after normal text",
			true,
			[]rune(`  xxx  %{ abc }`),
			[]rune(`  xxx  `),
			TextToken,
			nil,
			9,
		},
		{
			"after normal text with %",
			true,
			[]rune(`  xxx%%xx  %{ abc }`),
			[]rune(`  xxx%xx  `),
			TextToken,
			nil,
			13,
		},
		{
			"after normal text with %",
			true,
			[]rune(`  xxx%%xx  %{ abc }`),
			[]rune(`  xxx%xx  `),
			TextToken,
			nil,
			13,
		},
		{
			"first arg",
			false,
			[]rune("(`a`, `b`) }"),
			[]rune("a"),
			ArgToken,
			nil,
			4,
		},
		{
			"last arg",
			false,
			[]rune(", `b`) }"),
			[]rune("b"),
			ArgToken,
			nil,
			5,
		},
		{
			"end",
			false,
			[]rune(`) } xxxx`),
			[]rune(` xxxx`),
			TextToken,
			nil,
			8,
		},
	}

	for _, c := range cases {
		lexer := Lexer{
			buf:    c.buf,
			normal: c.isNormal,
		}
		token, tokenType, err := lexer.NextToken()
		assert.Equal(t, string(c.expectedToken), string(token), c.desc)
		assert.Equal(t, c.expectedTokenType, tokenType, c.desc)
		assert.Equal(t, c.expectedErr, err, c.desc)
		assert.Equal(t, c.expectedOffset, lexer.offset, c.desc)
	}
}
