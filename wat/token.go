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
	var scratch [48]byte
	return string(token.AppendTo(scratch[:0], true))
}

func (token Token) String() string {
	var scratch [32]byte
	return string(token.AppendTo(scratch[:0], false))
}

func (token Token) AppendTo(out []byte, verbose bool) []byte {
	if verbose {
		out = append(out, "wat.Token{"...)
		out = token.Type.AppendTo(out, verbose)
		out = append(out, ", "...)
		out = appendPretty(out, verbose, token.Value)
		out = append(out, ", "...)
		out = token.Span.AppendTo(out, verbose)
		out = append(out, "}"...)
		return out
	}
	out = token.Type.AppendTo(out, verbose)
	out = append(out, "("...)
	out = appendGuts(out, verbose, token.Value)
	out = append(out, ")"...)
	return out
}

func (token Token) IsTerminal() bool {
	return token.Type.IsTerminal()
}

func (token Token) Validate() error {
	switch token.Type {
	case AcceptToken:
		fallthrough
	case OpenParenToken:
		fallthrough
	case CloseParenToken:
		return token.validateNil()
	case RejectToken:
		return token.validateError()
	case SpaceToken:
		return token.validateSpace()
	case LineCommentToken:
		return token.validateString()
	case BlockCommentToken:
		return token.validateStringList()
	case KeywordToken:
		return token.validateString()
	case IdentifierToken:
		return token.validateString()
	case StringToken:
		return token.validateString()
	case NumberToken:
		return token.validateNumber()
	default:
		return fmt.Errorf("unknown token type %v", token.Type)
	}
}

func (token Token) validateNil() error {
	if token.Value != nil {
		return fmt.Errorf("%v token has non-nil value: %#v", token.Type, token.Value)
	}
	return nil
}

func (token Token) validateError() error {
	if token.Value == nil {
		return fmt.Errorf("%v token has nil value, not error", token.Type)
	}
	if _, ok := token.Value.(error); !ok {
		return fmt.Errorf("%v token has value of type %T, not error: %#v", token.Type, token.Value, token.Value)
	}
	return nil
}

func (token Token) validateSpace() error {
	if token.Value == nil {
		return fmt.Errorf("%v token has nil value, not wat.Space", token.Type)
	}
	if _, ok := token.Value.(Space); !ok {
		return fmt.Errorf("%v token has value of type %T, not wat.Space: %#v", token.Type, token.Value, token.Value)
	}
	return nil
}

func (token Token) validateNumber() error {
	if token.Value == nil {
		return fmt.Errorf("%v token has nil value, not wat.Num", token.Type)
	}
	if _, ok := token.Value.(Num); !ok {
		return fmt.Errorf("%v token has value of type %T, not wat.Num: %#v", token.Type, token.Value, token.Value)
	}
	return nil
}

func (token Token) validateString() error {
	if token.Value == nil {
		return fmt.Errorf("%v token has nil value, not string", token.Type)
	}
	if _, ok := token.Value.(string); !ok {
		return fmt.Errorf("%v token has value of type %T, not string: %#v", token.Type, token.Value, token.Value)
	}
	return nil
}

func (token Token) validateStringList() error {
	if token.Value == nil {
		return fmt.Errorf("%v token has nil value, not []string", token.Type)
	}
	if _, ok := token.Value.([]string); !ok {
		return fmt.Errorf("%v token has value of type %T, not []string: %#v", token.Type, token.Value, token.Value)
	}
	return nil
}

var (
	_ fmt.GoStringer = Token{}
	_ fmt.Stringer   = Token{}
	_ appenderTo     = Token{}
)
