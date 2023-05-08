package wat

import (
	"fmt"
	"unicode"
)

type Position struct {
	ByteOffset uint64
	RuneOffset uint64
	Line       uint
	Column     uint
	SkipLF     bool
}

func (pos Position) GoString() string {
	return fmt.Sprintf("Position{B:%d, R:%d, L:%d, C:%d, S:%t}", pos.ByteOffset, pos.RuneOffset, pos.Line, pos.Column, pos.SkipLF)
}

func (pos Position) String() string {
	return fmt.Sprintf("L:%d C:%d @ %d", pos.Line+1, pos.Column+1, pos.ByteOffset)
}

func (pos *Position) Advance(ch rune, size int) {
	if size < 0 {
		size = 0
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
)
