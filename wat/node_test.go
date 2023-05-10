package wat

import (
	"testing"
)

func TestNode_AppendTo(t *testing.T) {
	type testCase struct {
		Name     string
		Input    *Node
		ExpectGo string
		Expect   string
	}

	testData := [...]testCase{
		{
			Name:     "Nil",
			Input:    nil,
			ExpectGo: `nil`,
			Expect:   `<nil>`,
		},
		{
			Name:     "Invalid",
			Input:    &Node{Type: InvalidNode},
			ExpectGo: `&wat.Node{wat.InvalidNode, nil, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `<invalid>()`,
		},
		{
			Name:     "Space",
			Input:    &Node{Type: SpaceNode, Value: Space{LF, 1}},
			ExpectGo: `&wat.Node{wat.SpaceNode, wat.Space{wat.LF, 1}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Space(LF*1)`,
		},
		{
			Name:     "LineComment",
			Input:    &Node{Type: LineCommentNode, Value: "blah"},
			ExpectGo: `&wat.Node{wat.LineCommentNode, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `LineComment("blah")`,
		},
		{
			Name:     "BlockComment",
			Input:    &Node{Type: BlockCommentNode, Value: L("a", "b", "c")},
			ExpectGo: `&wat.Node{wat.BlockCommentNode, ["a", "b", "c"], wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `BlockComment("a", "b", "c")`,
		},
		{
			Name:     "Keyword",
			Input:    &Node{Type: KeywordNode, Value: "blah"},
			ExpectGo: `&wat.Node{wat.KeywordNode, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Keyword("blah")`,
		},
		{
			Name:     "Identifier",
			Input:    &Node{Type: IdentifierNode, Value: "blah"},
			ExpectGo: `&wat.Node{wat.IdentifierNode, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Identifier("blah")`,
		},
		{
			Name:     "String",
			Input:    &Node{Type: StringNode, Value: "blah"},
			ExpectGo: `&wat.Node{wat.StringNode, "blah", wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `String("blah")`,
		},
		{
			Name:     "Number",
			Input:    &Node{Type: NumberNode, Value: N(0, "0")},
			ExpectGo: `&wat.Node{wat.NumberNode, wat.Num{0, "0"}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect:   `Number(0)`,
		},
		{
			Name: "Expr",
			Input: &Node{
				Type: ExprNode,
				Value: []*Node{
					&Node{Type: NumberNode, Value: N(0, "1")},
					&Node{Type: NumberNode, Value: N(0, "2")},
					&Node{Type: NumberNode, Value: N(0, "3")},
				},
			},
			ExpectGo: `&wat.Node{wat.ExprNode, [` +
				`&wat.Node{wat.NumberNode, wat.Num{0, "1"}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}` +
				`, ` +
				`&wat.Node{wat.NumberNode, wat.Num{0, "2"}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}` +
				`, ` +
				`&wat.Node{wat.NumberNode, wat.Num{0, "3"}, wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}` +
				`], wat.Span{wat.Position{B:0, R:0, L:0, C:0, S:false}, wat.Position{B:0, R:0, L:0, C:0, S:false}}}`,
			Expect: `Expr(Number(1), Number(2), Number(3))`,
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
