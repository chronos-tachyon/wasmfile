package wat

import (
	"fmt"
	"strconv"
)

type appenderTo interface {
	AppendTo([]byte, bool) []byte
}

type gutsAppender interface {
	appendGuts([]byte, bool) []byte
}

func appendGuts(out []byte, verbose bool, v any) []byte {
	switch x := v.(type) {
	case nil:
		return out

	case gutsAppender:
		return x.appendGuts(out, verbose)

	case []string:
		for i, str := range x {
			if i > 0 {
				out = append(out, ", "...)
			}
			out = strconv.AppendQuote(out, str)
		}
		return out

	case []*Node:
		for i, node := range x {
			if i > 0 {
				out = append(out, ", "...)
			}
			out = node.AppendTo(out, verbose)
		}
		return out

	default:
		return appendPretty(out, verbose, v)
	}
}

func appendPretty(out []byte, verbose bool, v any) []byte {
	if verbose {
		switch x := v.(type) {
		case nil:
			return append(out, "nil"...)
		case appenderTo:
			return x.AppendTo(out, verbose)
		case fmt.GoStringer:
			str := x.GoString()
			return append(out, str...)
		case fmt.Stringer:
			str := x.String()
			return append(out, str...)
		}
	} else {
		switch x := v.(type) {
		case nil:
			return append(out, "<nil>"...)
		case appenderTo:
			return x.AppendTo(out, verbose)
		case fmt.Stringer:
			str := x.String()
			return append(out, str...)
		case fmt.GoStringer:
			str := x.GoString()
			return append(out, str...)
		}
	}

	switch x := v.(type) {
	case error:
		out = append(out, "err:"...)
		return strconv.AppendQuote(out, x.Error())

	case string:
		return strconv.AppendQuote(out, x)

	case []string:
		out = append(out, "["...)
		for i, str := range x {
			if i > 0 {
				out = append(out, ", "...)
			}
			out = strconv.AppendQuote(out, str)
		}
		out = append(out, "]"...)
		return out

	case []*Node:
		out = append(out, "["...)
		for i, node := range x {
			if i > 0 {
				out = append(out, ", "...)
			}
			out = node.AppendTo(out, verbose)
		}
		out = append(out, "]"...)
		return out

	default:
		var str string
		if verbose {
			str = fmt.Sprintf("%#v", v)
		} else {
			str = fmt.Sprintf("%v", v)
		}
		return append(out, str...)
	}
}
