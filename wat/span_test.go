package wat

import (
	"testing"
)

func TestSpan_AppendTo(t *testing.T) {
	span := Span{
		Begin: Position{
			ByteOffset: 100,
			RuneOffset: 100,
			Line:       5,
			Column:     0,
			SkipLF:     false,
		},
		End: Position{
			ByteOffset: 200,
			RuneOffset: 190,
			Line:       9,
			Column:     0,
			SkipLF:     false,
		},
	}
	expectGo := "wat.Span{wat.Position{B:100, R:100, L:5, C:0, S:false}, wat.Position{B:200, R:190, L:9, C:0, S:false}}"
	expect := "L:6 C:1 @ 100 [90]"

	str := span.GoString()
	if str != expectGo {
		t.Errorf("GoString: wrong output\n\texpect: %s\n\tactual: %s", expectGo, str)
	}

	str = span.String()
	if str != expect {
		t.Errorf("String: wrong output\n\texpect: %s\n\tactual: %s", expect, str)
	}
}
