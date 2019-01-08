package template

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	UnclosedScriptError      = errors.New("unclosed script, missing '}'")
	UnrecognizedToken        = errors.New("only %% and %{ is allowed")
	UnclosedParenthesisError = errors.New("unclosed parenthesis, missing ')'")
)

type Lexer struct {
	buf []rune

	offset int

	normal bool
}

// Token defines template token
type Token int

const (
	VariableNameToken Token = iota
	FuncNameToken
	ArgToken
	ArgVariableToken
	TextToken

	UnknownToken
)

func NewLexer(raw []rune) *Lexer {
	return &Lexer{
		buf:    raw,
		normal: true,
	}
}

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

func (lr *Lexer) peek() (rune, bool) {
	b, ok := lr.top()
	if !ok {
		return 0, false
	}
	lr.offset++
	return b, true
}

// peekNext try to peek bs and return whether bs are successfully peeked
// if skipSpace, spaces can be skipped if they are before bs
func (lr *Lexer) peekNext(bs []rune, skipSpace bool) bool {
	if len(lr.buf[lr.offset:]) < len(bs) {
		return false
	}
	hasSkipped, index, skipOff := !skipSpace, 0, 0
	for _, c := range lr.buf[lr.offset:] {
		skipOff++
		switch c {
		case ' ', '\n', '\r', '\t':
			if skipSpace && !hasSkipped {
				continue
			}
		default:
			// set non-space flag
			hasSkipped = true
		}
		if bs[index] != c {
			return false
		}
		if hasSkipped {
			index++
		}
		if index >= len(bs) {
			break
		}
	}
	lr.offset += skipOff
	return true
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
		return nil, UnknownToken, UnclosedScriptError
	}
	switch b {
	case '(', ',':
		lr.offset++
		b, ok := lr.top()
		if !ok {
			return nil, UnknownToken, UnclosedParenthesisError
		}
		if b == ')' {
			return lr.NextToken()
		}
		return lr.nextArg()
	case ')':
		lr.offset++
		next, ok := lr.top()
		if !ok || next != '}' {
			return nil, UnknownToken, UnclosedScriptError
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
			return nil, UnknownToken, UnclosedScriptError
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
		return nil, false, UnrecognizedToken
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
		return nil, false, UnrecognizedToken
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
		return nil, UnknownToken, UnclosedParenthesisError
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
		if unicode.IsLetter(c) || c == '_' || unicode.IsDigit(c) {
			token = append(token, c)
		} else {
			lr.offset += i
			break
		}
	}
	return token, nil

}
