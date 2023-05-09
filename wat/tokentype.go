package wat

import (
	"fmt"
)

type TokenType byte

const (
	InvalidToken TokenType = iota
	AcceptToken
	RejectToken
	SpaceToken
	LineCommentToken
	BlockCommentToken
	KeywordToken
	IdentifierToken
	StringToken
	NumberToken
	OpenParenToken
	CloseParenToken
)

var tokenTypeGoNames = [...]string{
	"wat.InvalidToken",
	"wat.AcceptToken",
	"wat.RejectToken",
	"wat.SpaceToken",
	"wat.LineCommentToken",
	"wat.BlockCommentToken",
	"wat.KeywordToken",
	"wat.IdentifierToken",
	"wat.StringToken",
	"wat.NumberToken",
	"wat.OpenParenToken",
	"wat.CloseParenToken",
}

var tokenTypeNames = [...]string{
	"<invalid>",
	"Accept",
	"Reject",
	"Space",
	"LineComment",
	"BlockComment",
	"Keyword",
	"Identifier",
	"String",
	"Number",
	"OpenParen",
	"CloseParen",
}

var tokenTypeNodeTypes = [...]NodeType{
	InvalidNode,
	InvalidNode,
	InvalidNode,
	SpaceNode,
	LineCommentNode,
	BlockCommentNode,
	KeywordNode,
	IdentifierNode,
	StringNode,
	NumberNode,
	InvalidNode,
	InvalidNode,
}

func (enum TokenType) GoString() string {
	var scratch [24]byte
	return string(enum.AppendTo(scratch[:0], true))
}

func (enum TokenType) String() string {
	var scratch [24]byte
	return string(enum.AppendTo(scratch[:0], false))
}

func (enum TokenType) AppendTo(out []byte, verbose bool) []byte {
	names := tokenTypeNames
	if verbose {
		names = tokenTypeGoNames
	}
	var str string
	if enum < TokenType(len(names)) {
		str = names[enum]
	} else {
		str = fmt.Sprintf("wat.TokenType(%d)", byte(enum))
	}
	return append(out, str...)
}

func (enum TokenType) NodeType() NodeType {
	if enum < TokenType(len(tokenTypeNodeTypes)) {
		return tokenTypeNodeTypes[enum]
	}
	return InvalidNode
}

func (enum TokenType) IsTerminal() bool {
	return enum == AcceptToken || enum == RejectToken
}

var (
	_ fmt.GoStringer = TokenType(0)
	_ fmt.Stringer   = TokenType(0)
	_ appenderTo     = TokenType(0)
)
