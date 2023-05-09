package wat

import (
	"testing"
	"unicode/utf8"
)

func TestPosition_AppendTo(t *testing.T) {
	pos := Position{
		ByteOffset: 5555,
		RuneOffset: 4444,
		Line:       1111,
		Column:     42,
		SkipLF:     true,
	}
	expectGo := "wat.Position{B:5555, R:4444, L:1111, C:42, S:true}"
	expect := "L:1112 C:43 @ 5555"

	str := pos.GoString()
	if str != expectGo {
		t.Errorf("GoString: wrong output\n\texpect: %s\n\tactual: %s", expectGo, str)
	}

	str = pos.String()
	if str != expect {
		t.Errorf("String: wrong output\n\texpect: %s\n\tactual: %s", expect, str)
	}
}

func TestPosition_Advance(t *testing.T) {
	type testCase struct {
		Name   string
		Rune   rune
		Size   int
		Input  Position
		Expect Position
	}

	testData := [...]testCase{
		{
			Name:   "RuneError:0",
			Rune:   utf8.RuneError,
			Size:   0,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5555, 4444, 1111, 42, true},
		},
		{
			Name:   "RuneError:1",
			Rune:   utf8.RuneError,
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5555, 4444, 1111, 42, true},
		},
		{
			Name:   "a",
			Rune:   'a',
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5556, 4445, 1111, 43, false},
		},
		{
			Name:   "CR",
			Rune:   '\r',
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5556, 4445, 1112, 0, true},
		},
		{
			Name:   "LF",
			Rune:   '\n',
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5556, 4445, 1111, 42, false},
		},
		{
			Name:   "CR:noSkipLF",
			Rune:   '\r',
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, false},
			Expect: Position{5556, 4445, 1112, 0, true},
		},
		{
			Name:   "LF:noSkipLF",
			Rune:   '\n',
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, false},
			Expect: Position{5556, 4445, 1112, 0, false},
		},
		{
			Name:   "FF",
			Rune:   '\f',
			Size:   1,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5556, 4445, 1111, 42, false},
		},
		{
			Name:   "á",
			Rune:   'á',
			Size:   2,
			Input:  Position{5555, 4444, 1111, 42, true},
			Expect: Position{5557, 4445, 1111, 43, false},
		},
		{
			Name:   "HT:0",
			Rune:   '\t',
			Size:   1,
			Input:  Position{0, 0, 0, 0, false},
			Expect: Position{1, 1, 0, 8, false},
		},
		{
			Name:   "HT:1",
			Rune:   '\t',
			Size:   1,
			Input:  Position{1, 1, 0, 1, false},
			Expect: Position{2, 2, 0, 8, false},
		},
		{
			Name:   "HT:7",
			Rune:   '\t',
			Size:   1,
			Input:  Position{7, 7, 0, 7, false},
			Expect: Position{8, 8, 0, 8, false},
		},
		{
			Name:   "HT:8",
			Rune:   '\t',
			Size:   1,
			Input:  Position{1, 1, 0, 8, false},
			Expect: Position{2, 2, 0, 16, false},
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			pos := row.Input
			pos.Advance(row.Rune, row.Size)
			if pos != row.Expect {
				t.Errorf("Advance: wrong result\n\texpect: %#v\n\tactual: %#v", row.Expect, pos)
			}
		})
	}
}
