package wat

import (
	"fmt"
)

type SpaceType byte

const (
	SP SpaceType = iota
	HT
	LF
	CR
	CRLF
)

var spaceTypeGoNames = [...]string{
	"wat.SP",
	"wat.HT",
	"wat.LF",
	"wat.CR",
	"wat.CRLF",
}

var spaceTypeNames = [...]string{
	"SP",
	"HT",
	"LF",
	"CR",
	"CRLF",
}

var spaceTypeRunes = [...]rune{
	' ',
	'\t',
	'\n',
	'\r',
	'\n',
}

var spaceTypeTexts = [...]string{
	" ",
	"\t",
	"\n",
	"\r",
	"\r\n",
}

func (enum SpaceType) GoString() string {
	var scratch [16]byte
	return string(enum.AppendTo(scratch[:0], true))
}

func (enum SpaceType) String() string {
	var scratch [16]byte
	return string(enum.AppendTo(scratch[:0], false))
}

func (enum SpaceType) AppendTo(out []byte, verbose bool) []byte {
	names := spaceTypeNames
	if verbose {
		names = spaceTypeGoNames
	}
	var str string
	if enum < SpaceType(len(names)) {
		str = names[enum]
	} else {
		str = fmt.Sprintf("wat.SpaceType(%d)", byte(enum))
	}
	return append(out, str...)
}

func (enum SpaceType) Rune() rune {
	if enum < SpaceType(len(spaceTypeRunes)) {
		return spaceTypeRunes[enum]
	}
	return ' '
}

func (enum SpaceType) Text() string {
	if enum < SpaceType(len(spaceTypeTexts)) {
		return spaceTypeTexts[enum]
	}
	return " "
}

var (
	_ fmt.GoStringer = SpaceType(0)
	_ fmt.Stringer   = SpaceType(0)
	_ appenderTo     = SpaceType(0)
)
