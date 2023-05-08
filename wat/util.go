package wat

import (
	"fmt"
)

type HexBytes []byte

func (hex HexBytes) GoString() string {
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

func (hex HexBytes) String() string {
	return hex.GoString()
}

var (
	_ fmt.GoStringer = HexBytes(nil)
	_ fmt.Stringer   = HexBytes(nil)
)
