package scanner

import "fmt"

type Error struct {
	module string
	line   int
	column int
	str    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v:%v:%v %v", e.module, e.line, e.column, e.str)
}
