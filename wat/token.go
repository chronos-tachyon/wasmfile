package wat

import (
	"fmt"
)

type Token struct {
	Type  TokenType
	Value any
	Span  Span
}

func (token Token) GoString() string {
	return fmt.Sprintf("Token{%v, %#v, %#v}", token.Type, token.Value, token.Span)
}

func (token Token) String() string {
	switch x := token.Value.(type) {
	case nil:
		return token.Type.String() + "()"
	case error:
		return fmt.Sprintf("%v(%q)", token.Type, x.Error())
	default:
		return fmt.Sprintf("%v(%#v)", token.Type, token.Value)
	}
}

func (token Token) IsTerminal() bool {
	return token.Type.IsTerminal()
}

var (
	_ fmt.GoStringer = Token{}
	_ fmt.Stringer   = Token{}
)
