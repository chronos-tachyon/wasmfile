package leb128

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUnsigned(t *testing.T) {
	type testCase struct {
		Input   uint64
		Encoded []byte
	}

	testData := [...]testCase{
		{
			Input:   0x0000,
			Encoded: []byte{0x00},
		},
		{
			Input:   0x0001,
			Encoded: []byte{0x01},
		},
		{
			Input:   0x003f,
			Encoded: []byte{0x3f},
		},
		{
			Input:   0x0040,
			Encoded: []byte{0x40},
		},
		{
			Input:   0x007f,
			Encoded: []byte{0x7f},
		},
		{
			Input:   0x0080,
			Encoded: []byte{0x80, 0x01},
		},
		{
			Input:   0x3fff,
			Encoded: []byte{0xff, 0x7f},
		},
		{
			Input:   0x4000,
			Encoded: []byte{0x80, 0x80, 0x01},
		},
		{
			Input:   624485,
			Encoded: []byte{0xe5, 0x8e, 0x26},
		},
		{
			Input:   0x7fffffff,
			Encoded: []byte{0xff, 0xff, 0xff, 0xff, 0x07},
		},
		{
			Input:   0x80000000,
			Encoded: []byte{0x80, 0x80, 0x80, 0x80, 0x08},
		},
		{
			Input:   0x7fffffffffffffff,
			Encoded: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f},
		},
		{
			Input:   0x8000000000000000,
			Encoded: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01},
		},
	}

	for _, row := range testData {
		name := fmt.Sprintf("%d", row.Input)
		t.Run(name, func(t *testing.T) {
			var scratch [10]byte
			encoded := AppendUint64(scratch[:0], row.Input)
			if !bytes.Equal(encoded, row.Encoded) {
				t.Errorf("AppendUint64 gives wrong result:\n\texpect: %v\n\tactual: %v", PrettyBytes(row.Encoded), PrettyBytes(encoded))
			}
			rest, decoded, ok := Uint64(encoded)
			if !ok {
				t.Error("Uint64 fails unexpectedly")
				return
			}
			if decoded != row.Input {
				t.Errorf("Uint64 gives wrong result: expect %#04x, got %#04x", row.Input, decoded)
			}
			if len(rest) > 0 {
				t.Errorf("Uint64 gives leftover bytes: %v", PrettyBytes(rest))
			}
		})
	}
}

type PrettyBytes []byte

func (pb PrettyBytes) String() string {
	const hex = "0123456789abcdef"
	if len(pb) <= 0 {
		return "[]"
	}
	buffer := make([]byte, 0, len(pb)*6)
	buffer = append(buffer, '[')
	for i, ch := range pb {
		if i > 0 {
			buffer = append(buffer, ',', ' ')
		}
		buffer = append(buffer, '0', 'x', hex[ch>>4], hex[ch&0xf])
	}
	buffer = append(buffer, ']')
	return string(buffer)
}
