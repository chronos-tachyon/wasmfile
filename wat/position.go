package wat

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Position struct {
	ByteOffset uint64
	RuneOffset uint64
	Line       uint
	Column     uint
	SkipLF     bool
}

func (pos Position) GoString() string {
	var scratch [48]byte
	return string(pos.AppendTo(scratch[:0], true))
}

func (pos Position) String() string {
	var scratch [32]byte
	return string(pos.AppendTo(scratch[:0], false))
}

func (pos Position) AppendTo(out []byte, verbose bool) []byte {
	if verbose {
		out = append(out, "wat.Position{B:"...)
		out = strconv.AppendUint(out, pos.ByteOffset, 10)
		out = append(out, ',', ' ', 'R', ':')
		out = strconv.AppendUint(out, pos.RuneOffset, 10)
		out = append(out, ',', ' ', 'L', ':')
		out = strconv.AppendUint(out, uint64(pos.Line), 10)
		out = append(out, ',', ' ', 'C', ':')
		out = strconv.AppendUint(out, uint64(pos.Column), 10)
		out = append(out, ',', ' ', 'S', ':')
		out = strconv.AppendBool(out, pos.SkipLF)
		out = append(out, '}')
		return out
	}
	out = append(out, 'L', ':')
	out = strconv.AppendUint(out, uint64(pos.Line)+1, 10)
	out = append(out, ' ', 'C', ':')
	out = strconv.AppendUint(out, uint64(pos.Column)+1, 10)
	out = append(out, ' ', '@', ' ')
	out = strconv.AppendUint(out, pos.ByteOffset, 10)
	return out
}

func (pos *Position) Advance(ch rune, size int) {
	if size < 1 || (size == 1 && ch == utf8.RuneError) {
		return
	}

	pos.ByteOffset += uint64(size)
	pos.RuneOffset++

	skipLF := pos.SkipLF
	pos.SkipLF = false

	switch {
	case ch == '\r':
		pos.Line++
		pos.Column = 0
		pos.SkipLF = true

	case ch == '\n':
		if skipLF {
			break
		}
		pos.Line++
		pos.Column = 0

	case ch == '\t':
		width := 8 - (pos.Column & 7)
		pos.Column += width

	case ch < 0x20:
		// pass

	case unicode.IsControl(ch):
		// pass

	default:
		pos.Column++
	}
}

var (
	_ fmt.GoStringer = Position{}
	_ fmt.Stringer   = Position{}
	_ appenderTo     = Position{}
)
