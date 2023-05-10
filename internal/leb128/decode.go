package leb128

import (
	"math/bits"
)

func Len(in []byte) (uint, bool) {
	inLen := uint(len(in))
	i := uint(0)
	ok := false
	for i < inLen {
		ch := in[i]
		i++
		if (ch & 0x80) == 0 {
			ok = true
			break
		}
	}
	return i, ok
}

func UintN(in []byte, numBits uint) (rest []byte, out uint64, ok bool) {
	rest, out, ok, _, _ = decode(in, numBits)
	return
}

func Uint64(in []byte) (rest []byte, out uint64, ok bool) {
	return UintN(in, 64)
}

func Uint32(in []byte) (rest []byte, out uint32, ok bool) {
	var u64 uint64
	rest, u64, ok = UintN(in, 32)
	out = uint32(u64)
	return
}

func Uint16(in []byte) (rest []byte, out uint16, ok bool) {
	var u64 uint64
	rest, u64, ok = UintN(in, 16)
	out = uint16(u64)
	return
}

func Uint8(in []byte) (rest []byte, out uint8, ok bool) {
	var u64 uint64
	rest, u64, ok = UintN(in, 8)
	out = uint8(u64)
	return
}

func IntN(in []byte, numBits uint) (rest []byte, out int64, ok bool) {
	var tmp uint64
	var shift uint
	var neg bool
	rest, tmp, ok, shift, neg = decode(in, numBits)
	if ok && neg && shift < 64 {
		// extend the sign bit to fill unused bits of uint64
		tmp |= ^((uint64(1) << shift) - 1)
	}
	out = int64(tmp)
	return
}

func Int64(in []byte) (rest []byte, out int64, ok bool) {
	return IntN(in, 64)
}

func Int33(in []byte) (rest []byte, out uint32, neg bool, ok bool) {
	var s64 int64
	rest, s64, ok = IntN(in, 33)
	out = uint32(s64)
	neg = (s64 < 0)
	return
}

func Int32(in []byte) (rest []byte, out int32, ok bool) {
	var s64 int64
	rest, s64, ok = IntN(in, 32)
	out = int32(s64)
	return
}

func Int16(in []byte) (rest []byte, out int16, ok bool) {
	var s64 int64
	rest, s64, ok = IntN(in, 16)
	out = int16(s64)
	return
}

func Int8(in []byte) (rest []byte, out int8, ok bool) {
	var s64 int64
	rest, s64, ok = IntN(in, 8)
	out = int8(s64)
	return
}

func decode(in []byte, numBits uint) (rest []byte, out uint64, ok bool, shift uint, neg bool) {
	i := uint(0)
	inLen := uint(len(in))
	for i < inLen {
		ch := in[i]
		i++

		isZeroBit7 := (ch & 0x80) == 0
		ch &= 0x7f
		out |= uint64(ch) << shift

		if isZeroBit7 {
			shift += uint(8 - bits.LeadingZeros8(ch))
			ok = (shift <= numBits)
			neg = ((ch & 0x40) != 0)
			break
		}

		shift += 7
	}
	rest = in[i:]
	return
}
