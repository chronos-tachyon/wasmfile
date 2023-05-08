package wat

import (
	"fmt"
)

type nodeSlab struct {
	list [64]Node
	used uint
}

func (slab *nodeSlab) alloc(t NodeType, v any, s Span) *Node {
	i := slab.used
	if i >= uint(len(slab.list)) {
		return nil
	}
	node := &slab.list[i]
	*node = Node{Type: t, Value: v, Span: s}
	i++
	slab.used = i
	return node
}

type Parser struct {
	slabs           []*nodeSlab
	spaceCache      map[Space]*Node
	keywordCache    map[string]*Node
	identifierCache map[string]*Node
	strCache        map[string]*Node
	numCache        map[Num]*Node
	keepSpaces      bool
	keepComments    bool
	disableCaching  bool
}

func (p *Parser) KeepSpace(value bool) {
	p.keepSpaces = value
}

func (p *Parser) KeepComments(value bool) {
	p.keepComments = value
}

func (p *Parser) DisableCaching(value bool) {
	p.disableCaching = value
}

func (p *Parser) Node(t NodeType, v any, s Span) *Node {
	if p == nil {
		return &Node{Type: t, Value: v, Span: s}
	}

	if p.disableCaching {
		return p.createNode(t, v, s)
	}

	switch t {
	case SpaceNode:
		sp := v.(Space)
		if node := p.spaceCache[sp]; node != nil {
			return node
		}
		node := p.createNode(t, v, s)
		if p.spaceCache == nil {
			p.spaceCache = make(map[Space]*Node, 16)
		}
		p.spaceCache[sp] = node
		return node

	case KeywordNode:
		str := v.(string)
		if node := p.keywordCache[str]; node != nil {
			return node
		}
		node := p.createNode(t, v, s)
		if p.keywordCache == nil {
			p.keywordCache = make(map[string]*Node, 16)
		}
		p.keywordCache[str] = node
		return node

	case IdentifierNode:
		str := v.(string)
		if node := p.identifierCache[str]; node != nil {
			return node
		}
		node := p.createNode(t, v, s)
		if p.identifierCache == nil {
			p.identifierCache = make(map[string]*Node, 16)
		}
		p.identifierCache[str] = node
		return node

	case StrNode:
		str := v.(string)
		if node := p.strCache[str]; node != nil {
			return node
		}
		node := p.createNode(t, v, s)
		if p.strCache == nil {
			p.strCache = make(map[string]*Node, 16)
		}
		p.strCache[str] = node
		return node

	case NumNode:
		num := v.(Num)
		if node := p.numCache[num]; node != nil {
			return node
		}
		node := p.createNode(t, v, s)
		if p.numCache == nil {
			p.numCache = make(map[Num]*Node, 16)
		}
		p.numCache[num] = node
		return node

	default:
		return p.createNode(t, v, s)
	}
}

func (p *Parser) createNode(t NodeType, v any, s Span) *Node {
	i := uint(len(p.slabs))
	if i > 0 {
		slab := p.slabs[i-1]
		node := slab.alloc(t, v, s)
		if node != nil {
			return node
		}
	}

	slab := &nodeSlab{}
	p.slabs = append(p.slabs, slab)
	return slab.alloc(t, v, s)
}

func (p *Parser) Parse(lexer *Lexer) (*Node, error) {
	if p == nil {
		p = new(Parser)
	}

	stack := make([]*Node, 0, 16)
	list := make([]*Node, 0, 16)
	root := p.Node(RootNode, list, Span{})
	top := root
	for lexer.HasNext() {
		t := lexer.Next()
		switch t.Type {
		case AcceptToken:
			root.Span.End = t.Span.End
			return root, nil

		case RejectToken:
			return nil, t.Value.(error)

		case OpenParenToken:
			node := p.Node(ExprNode, nil, t.Span)
			top.Value = append(list, node)
			stack = append(stack, node)
			top = node
			list = make([]*Node, 0, 16)

		case CloseParenToken:
			stackLen := len(stack)
			if stackLen < 1 {
				return nil, fmt.Errorf("unmatched '(' at %v", top.Span.Begin)
			}

			stackLen--
			stack[stackLen] = nil
			stack = stack[:stackLen]

			top.Span.End = t.Span.End
			top = root
			if stackLen > 0 {
				top = stack[stackLen-1]
			}
			list = top.Value.([]*Node)

		case SpaceToken:
			if p.keepSpaces {
				list = append(list, p.Node(SpaceNode, t.Value, t.Span))
			}

		case LineCommentToken:
			if p.keepComments {
				list = append(list, p.Node(LineCommentNode, t.Value, t.Span))
			}

		case BlockCommentToken:
			if p.keepComments {
				list = append(list, p.Node(BlockCommentNode, t.Value, t.Span))
			}

		case KeywordToken:
			list = append(list, p.Node(KeywordNode, t.Value, t.Span))

		case IdentifierToken:
			list = append(list, p.Node(IdentifierNode, t.Value, t.Span))

		case StrToken:
			list = append(list, p.Node(StrNode, t.Value, t.Span))

		case NumToken:
			list = append(list, p.Node(NumNode, t.Value, t.Span))

		default:
			return nil, fmt.Errorf("unexpected token %v", t)
		}
		top.Value = list
	}
	panic("unreachable")
}
