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
	if enum < SpaceType(len(spaceTypeNames)) {
		return spaceTypeNames[enum]
	}
	return fmt.Sprintf("SpaceType(%d)", uint(enum))
}

func (enum SpaceType) String() string {
	return enum.GoString()
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
)
