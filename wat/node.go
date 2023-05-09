package wat

import (
	"fmt"
)

type Node struct {
	Type  NodeType
	Value any
	Span  Span
}

func (node *Node) GoString() string {
	var scratch [1024]byte
	return string(node.AppendTo(scratch[:0], true))
}

func (node *Node) String() string {
	var scratch [1024]byte
	return string(node.AppendTo(scratch[:0], false))
}

func (node *Node) AppendTo(out []byte, verbose bool) []byte {
	if node == nil {
		str := "<nil>"
		if verbose {
			str = "nil"
		}
		return append(out, str...)
	}
	if verbose {
		out = append(out, "&wat.Node{"...)
		out = node.Type.AppendTo(out, verbose)
		out = append(out, ", "...)
		out = appendPretty(out, verbose, node.Value)
		out = append(out, ", "...)
		out = node.Span.AppendTo(out, verbose)
		out = append(out, "}"...)
		return out
	}
	out = node.Type.AppendTo(out, verbose)
	out = append(out, "("...)
	out = appendGuts(out, verbose, node.Value)
	out = append(out, ")"...)
	return out
}

func (node *Node) Equals(other *Node) bool {
	if node == other {
		return true
	}
	if node == nil || other == nil {
		return false
	}
	if node.Type != other.Type {
		return false
	}

	switch node.Type {
	case RootNode:
		fallthrough
	case ExprNode:
		a := node.Value.([]*Node)
		b := other.Value.([]*Node)
		aLen := uint(len(a))
		bLen := uint(len(b))
		if aLen != bLen {
			return false
		}
		for i := uint(0); i < aLen; i++ {
			if !a[i].Equals(b[i]) {
				return false
			}
		}
		return true

	case BlockCommentNode:
		a := node.Value.([]string)
		b := other.Value.([]string)
		aLen := uint(len(a))
		bLen := uint(len(b))
		if aLen != bLen {
			return false
		}
		for i := uint(0); i < aLen; i++ {
			if a[i] != b[i] {
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

func (node *Node) Validate(recursive bool) error {
	if node == nil {
		return nil
	}
	switch node.Type {
	case RootNode:
		fallthrough
	case ExprNode:
		return node.validateChildren(recursive)
	case SpaceNode:
		return node.validateSpace()
	case LineCommentNode:
		fallthrough
	case KeywordNode:
		fallthrough
	case IdentifierNode:
		fallthrough
	case StringNode:
		return node.validateString()
	case BlockCommentNode:
		return node.validateStringList()
	case NumberNode:
		return node.validateNumber()
	default:
		return fmt.Errorf("unknown node type %v", node.Type)
	}
}

func (node *Node) validateChildren(recursive bool) error {
	if node.Value == nil {
		return fmt.Errorf("%v node has nil value, not []*wat.Node", node.Type)
	}
	list, ok := node.Value.([]*Node)
	if !ok {
		return fmt.Errorf("%v node has value of type %T, not []*wat.Node: %#v", node.Type, node.Value, node.Value)
	}
	if recursive {
		for _, child := range list {
			if err := child.Validate(recursive); err != nil {
				return err
			}
		}
	}
	return nil
}

func (node *Node) validateSpace() error {
	if node.Value == nil {
		return fmt.Errorf("%v node has nil value, not wat.Space", node.Type)
	}
	if _, ok := node.Value.(Space); !ok {
		return fmt.Errorf("%v node has value of type %T, not wat.Space: %#v", node.Type, node.Value, node.Value)
	}
	return nil
}

func (node *Node) validateNumber() error {
	if node.Value == nil {
		return fmt.Errorf("%v node has nil value, not wat.Num", node.Type)
	}
	if _, ok := node.Value.(Num); !ok {
		return fmt.Errorf("%v node has value of type %T, not wat.Num: %#v", node.Type, node.Value, node.Value)
	}
	return nil
}

func (node *Node) validateString() error {
	if node.Value == nil {
		return fmt.Errorf("%v node has nil value, not string", node.Type)
	}
	if _, ok := node.Value.(string); !ok {
		return fmt.Errorf("%v node has value of type %T, not string: %#v", node.Type, node.Value, node.Value)
	}
	return nil
}

func (node *Node) validateStringList() error {
	if node.Value == nil {
		return fmt.Errorf("%v node has nil value, not []string", node.Type)
	}
	if _, ok := node.Value.([]string); !ok {
		return fmt.Errorf("%v node has value of type %T, not []string: %#v", node.Type, node.Value, node.Value)
	}
	return nil
}

var (
	_ fmt.GoStringer = (*Node)(nil)
	_ fmt.Stringer   = (*Node)(nil)
	_ appenderTo     = (*Node)(nil)
)
