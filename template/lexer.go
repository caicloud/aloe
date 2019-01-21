package template

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	// ErrUnclosedScript defines missing '}' error
	ErrUnclosedScript = errors.New("unclosed script, missing '}'")
	// ErrUnrecognizedToken defines unreognized token error
	ErrUnrecognizedToken = errors.New("only %% and %{ is allowed")
	// ErrUnclosedParenthesis defines missing ')' error
	ErrUnclosedParenthesis = errors.New("unclosed parenthesis, missing ')'")
)

// Lexer is lexer for template
type Lexer struct {
	buf []rune

	offset int

	normal bool
}

// Token defines template token
type Token int

const (
	// VariableNameToken defines variable name token
	VariableNameToken Token = iota
	// FuncNameToken defines func name token
	FuncNameToken
	// ArgToken defines function arg token
	ArgToken
	// ArgVariableToken defines function arg variable name token
	ArgVariableToken
	// TextToken defines normal text token
	TextToken

	// UnknownToken defines token unknown
	UnknownToken
)

// NewLexer new a lexer for raw
func NewLexer(raw []rune) *Lexer {
	return &Lexer{
		buf:    raw,
		normal: true,
	}
}

// IsEnd judge whether lexer is end
func (lr *Lexer) IsEnd() bool {
	if lr.offset == len(lr.buf) && lr.normal {
		return true
	}
	return false
}

func (lr *Lexer) top() (rune, bool) {
	for i, c := range lr.buf[lr.offset:] {
		switch c {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			lr.offset += i
			return c, true
		}
	}
	return 0, false
}

func (lr *Lexer) readUntil(delim rune) ([]rune, bool) {
	bs := []rune{}
	for _, c := range lr.buf[lr.offset:] {
		if c == delim {
			return bs, true
		}
		lr.offset++
		bs = append(bs, c)
	}
	return bs, false
}

// NextToken returns next token of the template
func (lr *Lexer) NextToken() ([]rune, Token, error) {
	if lr.normal {
		text, err := lr.nextTextToken()
		if err != nil {
			return nil, UnknownToken, err
		}
		return text, TextToken, nil
	}
	b, ok := lr.top()
	if !ok {
		return nil, UnknownToken, ErrUnclosedScript
	}
	switch b {
	case '(', ',':
		lr.offset++
		b, ok := lr.top()
		if !ok {
			return nil, UnknownToken, ErrUnclosedParenthesis
		}
		if b == ')' {
			return lr.NextToken()
		}
		return lr.nextArg()
	case ')':
		lr.offset++
		next, ok := lr.top()
		if !ok || next != '}' {
			return nil, UnknownToken, ErrUnclosedScript
		}
		lr.offset++
		lr.normal = true
		return lr.NextToken()
	default:
		token, err := lr.nextName()
		if err != nil {
			return nil, UnknownToken, err
		}
		b, ok := lr.top()
		if !ok {
			return nil, UnknownToken, ErrUnclosedScript
		}
		if b == '}' {
			lr.offset++
			lr.normal = true
			return token, VariableNameToken, nil
		}
		return token, FuncNameToken, nil
	}
}

// tryReadBeginToken try to read BeginToken
// %%: it's a normal text token, bool will be false and text will be returned
// %{: it's a begin token, bool will be true
// others: error will be returned
func (lr *Lexer) tryReadBeginToken() ([]rune, bool, error) {
	if len(lr.buf) == lr.offset+1 {
		return nil, false, ErrUnrecognizedToken
	}
	switch lr.buf[lr.offset+1] {
	case '%':
		lr.offset = lr.offset + 2
		bs, _ := lr.readUntil('%')
		return append([]rune(`%`), bs...), false, nil
	case '{':
		lr.offset = lr.offset + 2
		return nil, true, nil
	default:
		return nil, false, ErrUnrecognizedToken
	}

}

func (lr *Lexer) nextTextToken() ([]rune, error) {
	bs, ok := lr.readUntil('%')
	if !ok {
		return bs, nil
	}
	token := bs
	for lr.offset < len(lr.buf) {
		bs, isBeginToken, err := lr.tryReadBeginToken()
		if err != nil {
			return nil, err
		}
		token = append(token, bs...)
		if isBeginToken {
			lr.normal = false
			return token, nil
		}
	}
	return token, nil
}

func (lr *Lexer) nextArg() ([]rune, Token, error) {
	b, ok := lr.top()
	if !ok {
		return nil, UnknownToken, ErrUnclosedParenthesis
	}
	switch b {
	case '`':
		lr.offset++
		bs, ok := lr.readUntil('`')
		if !ok {
			return nil, UnknownToken, fmt.Errorf("unclosed quote `")
		}
		lr.offset++
		return bs, ArgToken, nil
	default:
		name, err := lr.nextName()
		if err != nil {
			return nil, UnknownToken, err
		}
		return name, ArgVariableToken, nil
	}
}

func (lr *Lexer) nextName() ([]rune, error) {
	first := lr.buf[lr.offset]
	if !unicode.IsLetter(first) && first != '_' {
		return nil, fmt.Errorf("variable name or function name should begin with unicode letter or '_'")
	}
	lr.offset++
	token := []rune{first}
	for i, c := range lr.buf[lr.offset:] {
		if unicode.IsLetter(c) || c == '_' || c == '.' || unicode.IsDigit(c) {
			token = append(token, c)
		} else {
			lr.offset += i
			break
		}
	}
	return token, nil

}
