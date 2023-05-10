package wat

import (
	"errors"
	"testing"
)

func TestToken_AppendTo(t *testing.T) {
	type testCase struct {
		Name     string
		Input    Token
		ExpectGo string
		Expect   string
	}

	testData := [...]testCase{
		{
			Name:     "Invalid",
			Input:    Token{Type: InvalidToken},
			ExpectGo: `wat.Token{wat.InvalidToken, nil, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `<invalid>()`,
		},
		{
			Name:     "Accept",
			Input:    Token{Type: AcceptToken},
			ExpectGo: `wat.Token{wat.AcceptToken, nil, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Accept()`,
		},
		{
			Name:     "Reject",
			Input:    Token{Type: RejectToken, Value: errors.New("blah")},
			ExpectGo: `wat.Token{wat.RejectToken, err:"blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Reject(err:"blah")`,
		},
		{
			Name:     "Space",
			Input:    Token{Type: SpaceToken, Value: Space{LF, 1}},
			ExpectGo: `wat.Token{wat.SpaceToken, wat.Space{wat.LF, 1}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Space(LF*1)`,
		},
		{
			Name:     "LineComment",
			Input:    Token{Type: LineCommentToken, Value: "blah"},
			ExpectGo: `wat.Token{wat.LineCommentToken, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `LineComment("blah")`,
		},
		{
			Name:     "BlockComment",
			Input:    Token{Type: BlockCommentToken, Value: L("a", "b", "c")},
			ExpectGo: `wat.Token{wat.BlockCommentToken, ["a", "b", "c"], wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `BlockComment("a", "b", "c")`,
		},
		{
			Name:     "Keyword",
			Input:    Token{Type: KeywordToken, Value: "blah"},
			ExpectGo: `wat.Token{wat.KeywordToken, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Keyword("blah")`,
		},
		{
			Name:     "Identifier",
			Input:    Token{Type: IdentifierToken, Value: "blah"},
			ExpectGo: `wat.Token{wat.IdentifierToken, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Identifier("blah")`,
		},
		{
			Name:     "String",
			Input:    Token{Type: StringToken, Value: "blah"},
			ExpectGo: `wat.Token{wat.StringToken, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `String("blah")`,
		},
		{
			Name:     "Number",
			Input:    Token{Type: NumberToken, Value: N(0, "0")},
			ExpectGo: `wat.Token{wat.NumberToken, wat.Num{0, "0"}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Number(0)`,
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
