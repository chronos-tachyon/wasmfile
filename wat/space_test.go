package wat

import (
	"testing"
)

func W(t SpaceType, n uint) Space {
	return Space{Type: t, Count: n}
}

func TestSpace_AppendTo(t *testing.T) {
	type testCase struct {
		Name     string
		Input    Space
		ExpectGo string
		Expect   string
	}

	testData := [...]testCase{
		{
			Name:     "SP-0",
			Input:    W(SP, 0),
			ExpectGo: "wat.Space{wat.SP, 0}",
			Expect:   "SP*0",
		},
		{
			Name:     "SP-1",
			Input:    W(SP, 1),
			ExpectGo: "wat.Space{wat.SP, 1}",
			Expect:   "SP*1",
		},
		{
			Name:     "SP-42",
			Input:    W(SP, 42),
			ExpectGo: "wat.Space{wat.SP, 42}",
			Expect:   "SP*42",
		},
		{
			Name:     "CRLF-5",
			Input:    W(CRLF, 5),
			ExpectGo: "wat.Space{wat.CRLF, 5}",
			Expect:   "CRLF*5",
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			str := row.Input.GoString()
			if str != row.ExpectGo {
				t.Errorf("GoString: wrong output\n\texpect: %s\n\tactual: %s", row.ExpectGo, str)
			}

			str = row.Input.String()
			if str != row.Expect {
				t.Errorf("String: wrong output\n\texpect: %s\n\tactual: %s", row.Expect, str)
			}
		})
	}
}
