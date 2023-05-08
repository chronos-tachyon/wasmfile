package wat

import (
	"fmt"
	"strconv"
)

type NumFlags uint32

const (
	FlagFloat NumFlags = (1 << iota)
	FlagInf
	FlagNaN
	FlagAcanonical
	FlagHex
	FlagSign
	FlagNeg
	FlagExpSign
	FlagExpNeg
)

var numFlagNames = [...]string{
	"FlagFloat",
	"FlagInf",
	"FlagNaN",
	"FlagAcanonical",
	"FlagHex",
	"FlagSign",
	"FlagNeg",
	"FlagExpSign",
	"FlagExpNeg",
}

func (bits NumFlags) GoString() string {
	if bits == 0 {
		return "0"
	}

	var scratch [64]byte
	var known NumFlags
	out := scratch[:0]
	first := true
	for index := uint(0); index < uint(len(numFlagNames)); index++ {
		name := numFlagNames[index]
		bit := NumFlags(1) << index
		known |= bit

		if (bits & bit) == 0 {
			continue
		}

		if !first {
			out = append(out, '|')
		}
		out = append(out, name...)
		first = false
	}
	if unknown := (bits &^ known); unknown != 0 {
		if !first {
			out = append(out, '|')
		}
		out = append(out, "0x"...)
		out = strconv.AppendUint(out, uint64(unknown), 16)
	}
	return string(out)
}

func (bits NumFlags) String() string {
	return bits.GoString()
}

func (bits NumFlags) HasAll(mask NumFlags) bool {
	return (bits & mask) == mask
}

func (bits NumFlags) HasAny(mask NumFlags) bool {
	return (bits & mask) != 0
}

func (bits NumFlags) HasNone(mask NumFlags) bool {
	return (bits & mask) == 0
}

var (
	_ fmt.GoStringer = NumFlags(0)
	_ fmt.Stringer   = NumFlags(0)
)
