package wat

import (
	"fmt"
)

type Span struct {
	Begin Position
	End   Position
}

func (span Span) GoString() string {
	return fmt.Sprintf("Span{%#v, %#v}", span.Begin, span.End)
}

func (span Span) String() string {
	p := span.Begin.ByteOffset
	q := span.End.ByteOffset
	return fmt.Sprintf("%v [%d]", span.Begin, int64(q-p))
}

var (
	_ fmt.GoStringer = Span{}
	_ fmt.Stringer   = Span{}
)
