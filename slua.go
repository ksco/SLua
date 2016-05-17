package main

import (
	"fmt"
	"strings"

	"github.com/ksco/SLua/parser"
	"github.com/ksco/slua/scanner"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	r := strings.NewReader("local さよなら = 'Hello, 世界'")
	s := scanner.New(r)
	p := parser.New(s)
	p.Parse()
}
