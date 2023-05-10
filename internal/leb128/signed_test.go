package leb128

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSigned(t *testing.T) {
	type testCase struct {
		Input   int64
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
			Input:   0x0002,
			Encoded: []byte{0x02},
		},
		{
			Input:   0x003f,
			Encoded: []byte{0x3f},
		},
		{
			Input:   0x0040,
			Encoded: []byte{0xc0, 0x00},
		},
		{
			Input:   0x007f,
			Encoded: []byte{0xff, 0x00},
		},
		{
			Input:   0x0080,
			Encoded: []byte{0x80, 0x01},
		},
		{
			Input:   0x3fff,
			Encoded: []byte{0xff, 0xff, 0x00},
		},
		{
			Input:   0x4000,
			Encoded: []byte{0x80, 0x80, 0x01},
		},
		{
			Input:   123456,
			Encoded: []byte{0xc0, 0xc4, 0x07},
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
			Encoded: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00},
		},
		{
			Input:   -0x0001,
			Encoded: []byte{0x7f},
		},
		{
			Input:   -0x0002,
			Encoded: []byte{0x7e},
		},
		{
			Input:   -0x0040,
			Encoded: []byte{0x40},
		},
		{
			Input:   -0x0041,
			Encoded: []byte{0xbf, 0x7f},
		},
		{
			Input:   -0x0080,
			Encoded: []byte{0x80, 0x7f},
		},
		{
			Input:   -0x0081,
			Encoded: []byte{0xff, 0x7e},
		},
		{
			Input:   -123456,
			Encoded: []byte{0xc0, 0xbb, 0x78},
		},
		{
			Input:   -0x80000000,
			Encoded: []byte{0x80, 0x80, 0x80, 0x80, 0x78},
		},
		{
			Input:   -0x80000001,
			Encoded: []byte{0xff, 0xff, 0xff, 0xff, 0x77},
		},
		{
			Input:   -0x8000000000000000,
			Encoded: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01},
		},
	}

	for _, row := range testData {
		name := fmt.Sprintf("%+d", row.Input)
		t.Run(name, func(t *testing.T) {
			var scratch [10]byte
			encoded := AppendInt64(scratch[:0], row.Input)
			if !bytes.Equal(encoded, row.Encoded) {
				t.Errorf("AppendInt64 gives wrong result:\n\texpect: %v\n\tactual: %v", PrettyBytes(row.Encoded), PrettyBytes(encoded))
				return
			}
			rest, decoded, ok := Int64(encoded)
			if !ok {
				t.Error("Int64 fails unexpectedly")
				return
			}
			if decoded != row.Input {
				t.Errorf("Int64 gives wrong result: expect %+#05x, got %+#05x", row.Input, decoded)
			}
			if len(rest) > 0 {
				t.Errorf("Int64 gives leftover bytes: %v", PrettyBytes(rest))
			}
		})
	}
}
