package leb128

import (
	"fmt"
	"math/bits"
)

type Uintish interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint
}

func LenUint[T Uintish](value T) uint {
	u64 := uint64(value)
	// numBits is the minimum number of bits that must be retained,
	// disregarding the 7-bit boundaries imposed by LEB128.
	numBits := 64 - uint(bits.LeadingZeros64(u64))
	return encodedLen(u64, numBits)
}

func PutUintN[T Uintish](out []byte, value T, numBits uint) uint {
	u64 := uint64(value)
	// minBits is the minimum number of bits that must be retained,
	// disregarding the 7-bit boundaries imposed by LEB128.
	minBits := 64 - uint(bits.LeadingZeros64(u64))
	if minBits > numBits {
		panic(fmt.Errorf("uint%d value %d requires %d bits to represent", numBits, u64, minBits))
	}
	return encode(out, u64, minBits)
}

func PutUint64[T Uintish](out []byte, value T) uint {
	return PutUintN(out, value, 64)
}

func PutUint32[T Uintish](out []byte, value T) uint {
	return PutUintN(out, value, 32)
}

func PutUint16[T Uintish](out []byte, value T) uint {
	return PutUintN(out, value, 16)
}

func PutUint8[T Uintish](out []byte, value T) uint {
	return PutUintN(out, value, 8)
}

func AppendUintN[T Uintish](out []byte, value T, numBits uint) []byte {
	const maxSize = 10
	return appendImpl(out, 10, func(p []byte) uint {
		return PutUintN(p, value, numBits)
	})
}

func AppendUint64[T Uintish](out []byte, value T) []byte {
	return AppendUintN(out, value, 64)
}

func AppendUint32[T Uintish](out []byte, value T) []byte {
	return AppendUintN(out, value, 32)
}

func AppendUint16[T Uintish](out []byte, value T) []byte {
	return AppendUintN(out, value, 16)
}

func AppendUint8[T Uintish](out []byte, value T) []byte {
	return AppendUintN(out, value, 8)
}

type Intish interface {
	~int64 | ~int32 | ~int16 | ~int8 | ~int
}

func LenInt[T Intish](value T) uint {
	s64 := int64(value)
	u64 := uint64(s64)

	// numBits is the minimum number of bits that must be retained,
	// disregarding the 7-bit boundaries imposed by LEB128.
	//
	// - For a non-negative value, all leading 0 bits can be discarded
	//   except one.  The retained 0 bit functions as the sign bit.
	//
	// - For a negative value, all leading 1 bits can be discarded except
	//   one.  The retained 1 bit functions as the sign bit.
	//
	// NB: This implies that 0 and -1 can be stored using only a single
	// sign bit with no mantissa bits, which is correct.
	//
	tmp := u64
	if s64 < 0 {
		tmp = ^tmp
	}
	n := uint(bits.LeadingZeros64(tmp))
	numBits := 64 - (n - 1)
	return encodedLen(u64, numBits)
}

func PutIntN[T Intish](out []byte, value T, numBits uint) uint {
	s64 := int64(value)
	u64 := uint64(s64)

	// minBits is the minimum number of bits that must be retained,
	// disregarding the 7-bit boundaries imposed by LEB128.
	//
	// - For a non-negative value, all leading 0 bits can be discarded
	//   except one.  The retained 0 bit functions as the sign bit.
	//
	// - For a negative value, all leading 1 bits can be discarded except
	//   one.  The retained 1 bit functions as the sign bit.
	//
	// NB: This implies that 0 and -1 can be stored using only a single
	// sign bit with no mantissa bits, which is correct.
	//
	tmp := u64
	if s64 < 0 {
		tmp = ^tmp
	}
	n := uint(bits.LeadingZeros64(tmp))
	minBits := 64 - (n - 1)
	if minBits > numBits {
		panic(fmt.Errorf("int%d value %d requires %d bits to represent", numBits, s64, minBits))
	}
	return encode(out, u64, minBits)
}

func PutInt64[T Intish](out []byte, value T) uint {
	return PutIntN(out, value, 64)
}

func PutInt32[T Intish](out []byte, value T) uint {
	return PutIntN(out, value, 32)
}

func PutInt16[T Intish](out []byte, value T) uint {
	return PutIntN(out, value, 16)
}

func PutInt8[T Intish](out []byte, value T) uint {
	return PutIntN(out, value, 8)
}

func AppendIntN[T Intish](out []byte, value T, numBits uint) []byte {
	const maxSize = 10
	return appendImpl(out, 10, func(p []byte) uint {
		return PutIntN(p, value, numBits)
	})
}

func AppendInt64[T Intish](out []byte, value T) []byte {
	return AppendIntN(out, value, 64)
}

func AppendInt32[T Intish](out []byte, value T) []byte {
	return AppendIntN(out, value, 32)
}

func AppendInt16[T Intish](out []byte, value T) []byte {
	return AppendIntN(out, value, 16)
}

func AppendInt8[T Intish](out []byte, value T) []byte {
	return AppendIntN(out, value, 8)
}

func PutInt33(out []byte, u32 uint32, neg bool) uint {
	s64 := int64(u32)
	if neg {
		s64 = -s64
	}
	return PutIntN(out, s64, 33)
}

func AppendInt33(out []byte, u32 uint32, neg bool) []byte {
	s64 := int64(u32)
	if neg {
		s64 = -s64
	}
	return AppendIntN(out, s64, 33)
}

func encodedLen(value uint64, numBits uint) uint {
	var numBytes uint = 1
	for numBits > 7 {
		numBits -= 7
		numBytes++
	}
	return numBytes
}

func encode(out []byte, value uint64, numBits uint) uint {
	var size uint
	for numBits > 7 {
		out[size] = byte(value) | 0x80
		size++
		value >>= 7
		numBits -= 7
	}
	out[size] = byte(value) & 0x7f
	size++
	return size
}

func appendImpl(out []byte, maxSize uint, putfn func([]byte) uint) []byte {
	outLen := uint(len(out))
	outCap := uint(cap(out))
	minCap := (outLen + maxSize)
	if minCap < outCap {
		newCap := outCap << 2
		if newCap < minCap {
			newCap = minCap
		}
		tmp := make([]byte, outLen, newCap)
		copy(tmp, out)
		out = tmp
	}
	outEnd := outLen + maxSize
	size := putfn(out[outLen:outEnd])
	outLen += size
	return out[:outLen]
}
