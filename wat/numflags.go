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

var numFlagGoNames = [...]string{
	"wat.FlagFloat",
	"wat.FlagInf",
	"wat.FlagNaN",
	"wat.FlagAcanonical",
	"wat.FlagHex",
	"wat.FlagSign",
	"wat.FlagNeg",
	"wat.FlagExpSign",
	"wat.FlagExpNeg",
}

var numFlagNames = [...]string{
	"Float",
	"Inf",
	"NaN",
	"Acanonical",
	"Hex",
	"Sign",
	"Neg",
	"ExpSign",
	"ExpNeg",
}

func (bits NumFlags) GoString() string {
	var scratch [64]byte
	return string(bits.AppendTo(scratch[:0], true))
}

func (bits NumFlags) String() string {
	var scratch [64]byte
	return string(bits.AppendTo(scratch[:0], false))
}

func (bits NumFlags) AppendTo(out []byte, verbose bool) []byte {
	if bits == 0 {
		return append(out, '0')
	}

	names := numFlagNames
	if verbose {
		names = numFlagGoNames
	}

	first := true
	sep := func() {
		if !first {
			out = append(out, '|')
		}
		first = false
	}

	var known NumFlags
	for index := uint(0); index < uint(len(numFlagNames)); index++ {
		name := names[index]
		bit := NumFlags(1) << index
		known |= bit
		if (bits & bit) == 0 {
			continue
		}
		sep()
		out = append(out, name...)
	}

	if unknown := (bits &^ known); unknown != 0 {
		sep()
		out = append(out, '0', 'x')
		out = strconv.AppendUint(out, uint64(unknown), 16)
	}

	return out
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
	_ appenderTo     = NumFlags(0)
)
