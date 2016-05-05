package main

import (
	"fmt"
	"strings"

	"github.com/ksco/slua/scanner"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	r := strings.NewReader("local 你好 = 'Hello, 世界'")
	s := scanner.New(r)

	t := s.Scan()
	for t.Category != scanner.TokenEOF {
		fmt.Printf("Type:%v Value:%v\n", t.Category, t)
		t = s.Scan()
	}
}
