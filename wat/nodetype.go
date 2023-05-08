package wat

import (
	"fmt"
)

type NodeType byte

const (
	RootNode NodeType = iota
	ExprNode
	SpaceNode
	LineCommentNode
	BlockCommentNode
	KeywordNode
	IdentifierNode
	StrNode
	NumNode
)

var nodeTypeNames = [...]string{
	"Root",
	"Expr",
	"Space",
	"LineComment",
	"BlockComment",
	"Keyword",
	"Identifier",
	"Str",
	"Num",
}

func (enum NodeType) GoString() string {
	if enum < NodeType(len(nodeTypeNames)) {
		return nodeTypeNames[enum]
	}
	return fmt.Sprintf("NodeType(%d)", byte(enum))
}

func (enum NodeType) String() string {
	return enum.GoString()
}

var (
	_ fmt.GoStringer = NodeType(0)
	_ fmt.Stringer   = NodeType(0)
)
