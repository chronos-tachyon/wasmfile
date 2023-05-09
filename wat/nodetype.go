package wat

import (
	"fmt"
)

type NodeType byte

const (
	InvalidNode NodeType = iota
	RootNode
	ExprNode
	SpaceNode
	LineCommentNode
	BlockCommentNode
	KeywordNode
	IdentifierNode
	StringNode
	NumberNode
)

var nodeTypeGoNames = [...]string{
	"wat.InvalidNode",
	"wat.RootNode",
	"wat.ExprNode",
	"wat.SpaceNode",
	"wat.LineCommentNode",
	"wat.BlockCommentNode",
	"wat.KeywordNode",
	"wat.IdentifierNode",
	"wat.StringNode",
	"wat.NumberNode",
}

var nodeTypeNames = [...]string{
	"<invalid>",
	"Root",
	"Expr",
	"Space",
	"LineComment",
	"BlockComment",
	"Keyword",
	"Identifier",
	"String",
	"Number",
}

func (enum NodeType) GoString() string {
	var scratch [24]byte
	return string(enum.AppendTo(scratch[:0], true))
}

func (enum NodeType) String() string {
	var scratch [24]byte
	return string(enum.AppendTo(scratch[:0], false))
}

func (enum NodeType) AppendTo(out []byte, verbose bool) []byte {
	names := nodeTypeNames
	if verbose {
		names = nodeTypeGoNames
	}
	var str string
	if enum < NodeType(len(names)) {
		str = names[enum]
	} else {
		str = fmt.Sprintf("wat.NodeType(%d)", byte(enum))
	}
	return append(out, str...)
}

var (
	_ fmt.GoStringer = NodeType(0)
	_ fmt.Stringer   = NodeType(0)
	_ appenderTo     = NodeType(0)
)
