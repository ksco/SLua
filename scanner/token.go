package scanner

import "fmt"

const (
	TokenAnd          string = "and"
	TokenDo                  = "do"
	TokenElse                = "else"
	TokenElseif              = "elseif"
	TokenEnd                 = "end"
	TokenFalse               = "false"
	TokenIf                  = "if"
	TokenLocal               = "local"
	TokenNil                 = "nil"
	TokenNot                 = "not"
	TokenOr                  = "or"
	TokenThen                = "then"
	TokenTrue                = "true"
	TokenWhile               = "while"
	TokenID                  = "<id>"
	TokenString              = "<string>"
	TokenNumber              = "<number>"
	TokenAdd                 = "+"
	TokenSub                 = "-"
	TokenMul                 = "*"
	TokenDiv                 = "/"
	TokenLen                 = "#"
	TokenLeftParen           = "("
	TokenRightParen          = ")"
	TokenAssign              = "="
	TokenSemicolon           = ";"
	TokenComma               = ","
	TokenEqual               = "=="
	TokenNotEqual            = "~="
	TokenLess                = "<"
	TokenLessEqual           = "<="
	TokenGreater             = ">"
	TokenGreaterEqual        = ">="
	TokenConcat              = ".."
	TokenEOF                 = "<eof>"
)

type Token struct {
	Value    interface{}
	Line     int
	Column   int
	Category string
}

func NewToken() *Token {
	token := new(Token)
	token.Category = TokenEOF
	return token
}

func (t *Token) String() string {
	var s string
	if t.Category == TokenNumber || t.Category == TokenID ||
		t.Category == TokenString {
		s = fmt.Sprintf("%v", t.Value)
	} else {
		s = t.Category
	}
	return s
}

func (t *Token) Clone() *Token {
	return &Token{
		Value:    t.Value,
		Line:     t.Line,
		Column:   t.Column,
		Category: t.Category,
	}
}

func isKeyword(id string) bool {
	switch id {
	case TokenAnd, TokenDo, TokenElse, TokenElseif, TokenEnd,
		TokenFalse, TokenIf, TokenLocal, TokenNil, TokenNot,
		TokenOr, TokenThen, TokenTrue, TokenWhile:
		return true
	default:
		return false
	}
}
