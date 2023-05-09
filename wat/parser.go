package wat

import (
	"fmt"
)

type Parser struct {
	slabs           []*nodeSlab
	spaceCache      map[Space]*Node
	numCache        map[Num]*Node
	keywordCache    map[string]*Node
	identifierCache map[string]*Node
	strCache        map[string]*Node
	keepSpaces      bool
	keepComments    bool
	disableCaching  bool
}

func (parser *Parser) KeepSpaces(value bool) *Parser {
	parser.keepSpaces = value
	return parser
}

func (parser *Parser) KeepComments(value bool) *Parser {
	parser.keepComments = value
	return parser
}

func (parser *Parser) DisableCaching(value bool) *Parser {
	parser.disableCaching = value
	return parser
}

func (parser *Parser) Parse(lexer TokenStream) (*Node, error) {
	if parser == nil {
		parser = new(Parser)
	}

	stack := make([]*Node, 0, 16)
	list := make([]*Node, 0, 16)
	root := parser.Node(RootNode, list, Span{})
	top := root

	add := func(child *Node) *Node {
		if err := child.Validate(false); err != nil {
			panic(err)
		}
		list = append(list, child)
		top.Value = list
		return child
	}

	for lexer.HasNext() {
		token := lexer.Next()
		if err := token.Validate(); err != nil {
			return nil, err
		}

		switch token.Type {
		case AcceptToken:
			if len(stack) > 0 {
				return nil, fmt.Errorf("unmatched '(' at %v", token.Span.Begin)
			}
			root.Span.End = token.Span.End
			return root, nil
		case RejectToken:
			return nil, token.Value.(error)
		case OpenParenToken:
			exprList := make([]*Node, 0, 16)
			exprNode := add(parser.Node(ExprNode, exprList, token.Span))
			stack = append(stack, exprNode)
			list = exprList
			top = exprNode
		case CloseParenToken:
			stackLen := len(stack)
			if stackLen < 1 {
				return nil, fmt.Errorf("unmatched ')' at %v", token.Span.Begin)
			}
			top.Span.End = token.Span.End
			stackLen--
			stack[stackLen] = nil
			stack = stack[:stackLen]
			top = root
			if stackLen > 0 {
				top = stack[stackLen-1]
			}
			list = top.Value.([]*Node)
		case SpaceToken:
			if parser.keepSpaces {
				add(parser.Node(SpaceNode, token.Value, token.Span))
			}
		case LineCommentToken:
			if parser.keepComments {
				add(parser.Node(LineCommentNode, token.Value, token.Span))
			}
		case BlockCommentToken:
			if parser.keepComments {
				add(parser.Node(BlockCommentNode, token.Value, token.Span))
			}
		case KeywordToken:
			add(parser.Node(KeywordNode, token.Value, token.Span))
		case IdentifierToken:
			add(parser.Node(IdentifierNode, token.Value, token.Span))
		case StringToken:
			add(parser.Node(StringNode, token.Value, token.Span))
		case NumberToken:
			add(parser.Node(NumberNode, token.Value, token.Span))
		default:
			return nil, fmt.Errorf("unexpected token %v", token)
		}
	}
	panic("unreachable")
}

func (parser *Parser) Node(tt NodeType, tv any, ts Span) *Node {
	if parser == nil {
		return &Node{Type: tt, Value: tv, Span: ts}
	}

	if !parser.disableCaching {
		switch tt {
		case SpaceNode:
			sp := tv.(Space)
			if node := parser.spaceCache[sp]; node != nil {
				return node
			}
			node := parser.createNode(tt, tv, ts)
			if parser.spaceCache == nil {
				parser.spaceCache = make(map[Space]*Node, 16)
			}
			parser.spaceCache[sp] = node
			return node

		case NumberNode:
			num := tv.(Num)
			if node := parser.numCache[num]; node != nil {
				return node
			}
			node := parser.createNode(tt, tv, ts)
			if parser.numCache == nil {
				parser.numCache = make(map[Num]*Node, 16)
			}
			parser.numCache[num] = node
			return node

		case KeywordNode:
			str := tv.(string)
			if node := parser.keywordCache[str]; node != nil {
				return node
			}
			node := parser.createNode(tt, tv, ts)
			if parser.keywordCache == nil {
				parser.keywordCache = make(map[string]*Node, 16)
			}
			parser.keywordCache[str] = node
			return node

		case IdentifierNode:
			str := tv.(string)
			if node := parser.identifierCache[str]; node != nil {
				return node
			}
			node := parser.createNode(tt, tv, ts)
			if parser.identifierCache == nil {
				parser.identifierCache = make(map[string]*Node, 16)
			}
			parser.identifierCache[str] = node
			return node

		case StringNode:
			str := tv.(string)
			if node := parser.strCache[str]; node != nil {
				return node
			}
			node := parser.createNode(tt, tv, ts)
			if parser.strCache == nil {
				parser.strCache = make(map[string]*Node, 16)
			}
			parser.strCache[str] = node
			return node
		}
	}

	return parser.createNode(tt, tv, ts)
}

func (parser *Parser) createNode(tt NodeType, tv any, ts Span) *Node {
	i := uint(len(parser.slabs))
	if i > 0 {
		slab := parser.slabs[i-1]
		node := slab.alloc(tt, tv, ts)
		if node != nil {
			return node
		}
	}

	slab := &nodeSlab{}
	parser.slabs = append(parser.slabs, slab)
	return slab.alloc(tt, tv, ts)
}

type nodeSlab struct {
	list [64]Node
	used uint
}

func (slab *nodeSlab) alloc(tt NodeType, tv any, ts Span) *Node {
	i := slab.used
	if i >= uint(len(slab.list)) {
		return nil
	}
	node := &slab.list[i]
	*node = Node{Type: tt, Value: tv, Span: ts}
	i++
	slab.used = i
	return node
}
