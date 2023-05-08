package wat

import (
	"fmt"
)

type hexBytes []byte

func (hex hexBytes) GoString() string {
	const hexDigits = "0123456789abcdef"
	var scratch [256]byte
	out := scratch[:0]
	hexLen := uint(len(hex))
	for i := uint(0); i < hexLen; i++ {
		ch := hex[i]
		if i > 0 {
			out = append(out, ' ')
		}
		out = append(out, hexDigits[ch>>4])
		out = append(out, hexDigits[ch&0xf])
	}
	return string(out)
}

func (hex hexBytes) String() string {
	return hex.GoString()
}

var (
	_ fmt.GoStringer = hexBytes(nil)
	_ fmt.Stringer   = hexBytes(nil)
)
