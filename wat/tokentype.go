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
	StrToken
	NumToken
	OpenParenToken
	CloseParenToken
)

var tokenTypeNames = [...]string{
	"Invalid",
	"Accept",
	"Reject",
	"Space",
	"LineComment",
	"BlockComment",
	"Keyword",
	"Identifier",
	"Str",
	"Num",
	"OpenParen",
	"CloseParen",
}

func (enum TokenType) GoString() string {
	if enum < TokenType(len(tokenTypeNames)) {
		return tokenTypeNames[enum]
	}
	return fmt.Sprintf("TokenType(%d)", uint(enum))
}

func (enum TokenType) String() string {
	return enum.GoString()
}

func (enum TokenType) IsTerminal() bool {
	return enum == AcceptToken || enum == RejectToken
}

var (
	_ fmt.GoStringer = TokenType(0)
	_ fmt.Stringer   = TokenType(0)
)
