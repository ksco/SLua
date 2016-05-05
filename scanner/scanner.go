package scanner

import (
	"io"
	"strconv"
	"unicode"

	"github.com/ksco/slua/ascii"
)

type Scanner struct {
	module  string
	reader  io.RuneReader
	current rune
	line    int
	column  int
	buffer  []rune
}

const eof rune = 0

func New(reader io.RuneReader) *Scanner {
	s := new(Scanner)
	s.module = "scanner"
	s.reader = reader
	s.line = 1
	s.current = eof
	return s
}

func (s *Scanner) Scan() *Token {
	if s.current == eof {
		s.current = s.next()
	}

	for s.current != eof {
		switch s.current {
		case ' ', '\t', '\v', '\f': // Skip whitespace
			s.current = s.next()
		case '\r', '\n':
			s.newLine()
		case '-':
			n := s.next()
			if n == '-' {
				s.comment()
			} else {
				s.current = n
				return s.normalToken(TokenSub)
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return s.number(false)

		case '+', '*', '/', '#', '(', ')', ';', ',':
			t := s.current
			s.current = s.next()
			return s.normalToken(string(t))
		case '.':
			n := s.next()
			if n == '.' {
				s.current = s.next()
				return s.normalToken(TokenConcat)
			} else {
				s.buffer = s.buffer[:0]
				s.buffer = append(s.buffer, s.current)
				s.current = n
				return s.number(true)
			}
		case '~':
			n := s.next()
			if n != '=' {
				panic(&Error{
					module: s.module,
					line:   s.line,
					column: s.column,
					str:    "expect '=' after '~'",
				})
			}
			s.current = s.next()
			return s.normalToken(TokenNotEqual)
		case '=':
			return s.xequal(TokenEqual)
		case '>':
			return s.xequal(TokenGreaterEqual)
		case '<':
			return s.xequal(TokenLessEqual)
		case '\'', '"':
			return s.singlelineString()
		default:
			return s.id()
		}
	}
	return NewToken()
}

// Helper functions

func isLetter(ch rune) bool {
	return ascii.IsLetter(byte(ch)) || ch == '_' ||
		ch >= 0x80 && unicode.IsLetter(ch)
}

func (s *Scanner) normalToken(category string) *Token {
	return &Token{
		Line:     s.line,
		Column:   s.column,
		Category: category,
	}
}

func (s *Scanner) stringToken(value, category string) *Token {
	t := s.normalToken(category)
	t.Value = value
	return t
}

func (s *Scanner) numberToken(value float64) *Token {
	t := s.normalToken(TokenNumber)
	t.Value = value
	return t
}

func (s *Scanner) next() rune {
	ch, _, err := s.reader.ReadRune()
	if err == nil {
		s.column++
	}
	return ch
}

func (s *Scanner) newLine() {
	ch := s.next()
	if (ch == '\r' || ch == '\n') && ch != s.current {
		s.current = s.next()
	} else {
		s.current = ch
	}
	s.line++
	s.column = 0
}

func (s *Scanner) comment() {
	s.current = s.next()
	for s.current != '\r' && s.current != '\n' && s.current != eof {
		s.current = s.next()
	}
}

func (s *Scanner) number(point bool) *Token {
	if !point {
		s.buffer = s.buffer[:0]
		for unicode.IsDigit(s.current) {
			s.buffer = append(s.buffer, s.current)
			s.current = s.next()
		}
		if s.current == '.' {
			s.buffer = append(s.buffer, s.current)
			s.current = s.next()
		}
	}
	for unicode.IsDigit(s.current) {
		s.buffer = append(s.buffer, s.current)
		s.current = s.next()
	}
	str := string(s.buffer)
	number, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(&Error{
			module: s.module,
			line:   s.line,
			column: s.column,
			str:    "parse number " + str + " error: invalid syntax",
		})
	}
	return s.numberToken(number)
}

func (s *Scanner) xequal(category string) *Token {
	t := s.current
	ch := s.next()
	if ch == '=' {
		s.current = s.next()
		return s.normalToken(category)
	}
	s.current = ch
	return s.normalToken(string(t))
}

func (s *Scanner) stringChar() {
	if s.current == '\\' {
		s.current = s.next()
		if s.current == 'a' {
			s.buffer = append(s.buffer, '\a')
		} else if s.current == 'b' {
			s.buffer = append(s.buffer, '\b')
		} else if s.current == 'f' {
			s.buffer = append(s.buffer, '\f')
		} else if s.current == 'n' {
			s.buffer = append(s.buffer, '\n')
		} else if s.current == 'r' {
			s.buffer = append(s.buffer, '\r')
		} else if s.current == 't' {
			s.buffer = append(s.buffer, '\t')
		} else if s.current == 'v' {
			s.buffer = append(s.buffer, '\v')
		} else if s.current == '\\' {
			s.buffer = append(s.buffer, '\\')
		} else if s.current == '"' {
			s.buffer = append(s.buffer, '"')
		} else if s.current == '\'' {
			s.buffer = append(s.buffer, '\'')
		} else {
			panic(&Error{
				module: s.module,
				line:   s.line,
				column: s.column,
				str:    "unexpect character after '\\'",
			})
		}
	} else {
		s.buffer = append(s.buffer, s.current)
	}
	s.current = s.next()
}

func (s *Scanner) singlelineString() *Token {
	quote := s.current
	s.current = s.next()
	s.buffer = s.buffer[:0]
	for s.current != quote {
		if s.current == eof {
			panic(&Error{
				module: s.module,
				line:   s.line,
				column: s.column,
				str:    "incomplete string at <eof>",
			})
		}
		if s.current == '\r' || s.current == '\n' {
			panic(&Error{
				module: s.module,
				line:   s.line,
				column: s.column,
				str:    "incomplete string at <eol>",
			})
		}
		s.stringChar()
	}
	s.current = s.next()
	return s.stringToken(string(s.buffer), TokenString)
}

func (s *Scanner) id() *Token {
	if !isLetter(s.current) {
		panic(&Error{
			module: s.module,
			line:   s.line,
			column: s.column,
			str:    "unexpect character",
		})
	}

	s.buffer = s.buffer[:0]
	s.buffer = append(s.buffer, s.current)
	s.current = s.next()
	for isLetter(s.current) || unicode.IsDigit(s.current) {
		s.buffer = append(s.buffer, s.current)
		s.current = s.next()
	}

	str := string(s.buffer)
	if isKeyword(str) {
		return s.stringToken(str, str)
	}
	return s.stringToken(str, TokenID)
}
