package wat

import (
	"fmt"
	"strconv"
)

type Num struct {
	Flags    NumFlags
	Integer  string
	Fraction string
	Exponent string
}

func (num Num) GoString() string {
	var scratch [64]byte
	out := scratch[:0]
	out = append(out, "Num{"...)
	out = append(out, num.Flags.String()...)
	out = append(out, ", "...)
	out = strconv.AppendQuote(out, num.Integer)
	out = append(out, ", "...)
	out = strconv.AppendQuote(out, num.Fraction)
	out = append(out, ", "...)
	out = strconv.AppendQuote(out, num.Exponent)
	out = append(out, "}"...)
	return string(out)
}

func (num Num) String() string {
	var scratch [64]byte
	out := scratch[:0]

	switch num.Flags & (FlagSign | FlagNeg) {
	case FlagSign | FlagNeg:
		out = append(out, '-')
	case FlagSign:
		out = append(out, '+')
	}

	if num.Flags.HasAny(FlagNaN) {
		out = append(out, 'n', 'a', 'n')
		if num.Flags.HasAny(FlagAcanonical) {
			out = append(out, ':', '0', 'x')
			out = append(out, num.Integer...)
		}
		return string(out)
	}

	if num.Flags.HasAny(FlagInf) {
		out = append(out, 'i', 'n', 'f')
		return string(out)
	}

	expChar := byte('e')
	if num.Flags.HasAny(FlagHex) {
		expChar = 'p'
		out = append(out, '0', 'x')
	}

	out = append(out, num.Integer...)
	if num.Fraction != "" {
		out = append(out, '.')
		out = append(out, num.Fraction...)
	}
	if num.Exponent != "" {
		out = append(out, expChar)
		switch num.Flags & (FlagExpSign | FlagExpNeg) {
		case FlagExpSign | FlagExpNeg:
			out = append(out, '-')
		case FlagExpSign:
			out = append(out, '+')
		}
		out = append(out, num.Exponent...)
	}
	return string(out)
}

var (
	_ fmt.GoStringer = Num{}
	_ fmt.Stringer   = Num{}
)
