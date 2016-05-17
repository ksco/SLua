package parser

import (
	"fmt"

	"github.com/ksco/slua/scanner"
)

type Error struct {
	module string
	token  *scanner.Token
	str    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v:%v:%v: '%v' %v", e.module, e.token.Line,
		e.token.Column, e.token.String(), e.str)
}

func assert(cond bool, msg string) {
	if !cond {
		panic("lemon/parser internal error: " + msg)
	}
}
