package wat

import (
	"fmt"
)

type Space struct {
	Type  SpaceType
	Count uint
}

func (sp Space) GoString() string {
	return fmt.Sprintf("%v*%d", sp.Type, sp.Count)
}

func (sp Space) String() string {
	return sp.GoString()
}

var (
	_ fmt.GoStringer = Space{}
	_ fmt.Stringer   = Space{}
)
