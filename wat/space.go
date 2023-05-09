package wat

import (
	"fmt"
	"strconv"
)

type Space struct {
	Type  SpaceType
	Count uint
}

func (sp Space) GoString() string {
	var scratch [16]byte
	return string(sp.AppendTo(scratch[:0], true))
}

func (sp Space) String() string {
	var scratch [16]byte
	return string(sp.AppendTo(scratch[:0], false))
}

func (sp Space) AppendTo(out []byte, verbose bool) []byte {
	if verbose {
		out = append(out, "wat.Space{"...)
		out = sp.Type.AppendTo(out, verbose)
		out = append(out, ", "...)
		out = strconv.AppendUint(out, uint64(sp.Count), 10)
		out = append(out, "}"...)
		return out
	}
	out = sp.Type.AppendTo(out, verbose)
	out = append(out, '*')
	out = strconv.AppendUint(out, uint64(sp.Count), 10)
	return out
}

var (
	_ fmt.GoStringer = Space{}
	_ fmt.Stringer   = Space{}
	_ appenderTo     = Space{}
)
