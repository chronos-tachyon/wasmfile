package wat

import (
	"fmt"
	"strconv"
)

type Span struct {
	Begin Position
	End   Position
}

func (span Span) GoString() string {
	var scratch [48]byte
	return string(span.AppendTo(scratch[:0], true))
}

func (span Span) String() string {
	var scratch [32]byte
	return string(span.AppendTo(scratch[:0], false))
}

func (span Span) AppendTo(out []byte, verbose bool) []byte {
	if verbose {
		out = append(out, "wat.Span{"...)
		out = span.Begin.AppendTo(out, verbose)
		out = append(out, ", "...)
		out = span.End.AppendTo(out, verbose)
		out = append(out, "}"...)
		return out
	}
	delta := int64(span.End.RuneOffset - span.Begin.RuneOffset)
	out = span.Begin.AppendTo(out, verbose)
	out = append(out, " ["...)
	out = strconv.AppendInt(out, delta, 10)
	out = append(out, "]"...)
	return out
}

var (
	_ fmt.GoStringer = Span{}
	_ fmt.Stringer   = Span{}
	_ appenderTo     = Span{}
)
