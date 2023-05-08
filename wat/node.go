package wat

import (
	"fmt"
	"strconv"
)

type Node struct {
	Type  NodeType
	Value any
	Span  Span
}

func (node *Node) Equals(other *Node) bool {
	if node == nil {
		return (other == nil)
	}
	if other == nil {
		return false
	}
	if node.Type != other.Type {
		return false
	}
	switch node.Type {
	case RootNode:
		fallthrough
	case ExprNode:
		av := node.Value.([]*Node)
		bv := other.Value.([]*Node)
		avLen := uint(len(av))
		bvLen := uint(len(bv))
		if avLen != bvLen {
			return false
		}
		for i := uint(0); i < avLen; i++ {
			if !av[i].Equals(bv[i]) {
				return false
			}
		}
		return true

	case BlockCommentNode:
		av := node.Value.([]string)
		bv := other.Value.([]string)
		avLen := uint(len(av))
		bvLen := uint(len(bv))
		if avLen != bvLen {
			return false
		}
		for i := uint(0); i < avLen; i++ {
			if av[i] != bv[i] {
				return false
			}
		}
		return true

	case SpaceNode:
		a := node.Value.(Space)
		b := other.Value.(Space)
		if a == b {
			return true
		}
		aIsNL := (a.Type == LF || a.Type == CR || a.Type == CRLF)
		bIsNL := (b.Type == LF || b.Type == CR || b.Type == CRLF)
		return (a.Count == b.Count) && (aIsNL == bIsNL)

	default:
		return (node.Value == other.Value)
	}
}

func (node *Node) AppendTo(out []byte) []byte {
	if node == nil {
		return append(out, "nil"...)
	}
	out = append(out, node.Type.String()...)
	out = append(out, '(')
	switch node.Type {
	case RootNode:
		fallthrough
	case ExprNode:
		for i, child := range node.Value.([]*Node) {
			if i > 0 {
				out = append(out, ',', ' ')
			}
			out = child.AppendTo(out)
		}

	case SpaceNode:
		out = append(out, node.Value.(Space).String()...)

	case LineCommentNode:
		fallthrough
	case KeywordNode:
		fallthrough
	case IdentifierNode:
		fallthrough
	case StrNode:
		out = strconv.AppendQuote(out, node.Value.(string))

	case BlockCommentNode:
		for i, str := range node.Value.([]string) {
			if i > 0 {
				out = append(out, ',', ' ')
			}
			out = strconv.AppendQuote(out, str)
		}

	case NumNode:
		out = append(out, node.Value.(Num).String()...)

	default:
		out = append(out, fmt.Sprintf("%#v", node.Value)...)
	}
	out = append(out, ')')
	return out
}

func (node *Node) GoString() string {
	var scratch [1024]byte
	return string(node.AppendTo(scratch[:0]))
}

func (node *Node) String() string {
	return node.GoString()
}

var (
	_ fmt.GoStringer = (*Node)(nil)
	_ fmt.Stringer   = (*Node)(nil)
)
