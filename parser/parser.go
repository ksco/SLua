package parser

import (
	"github.com/ksco/slua/scanner"
	"github.com/ksco/slua/syntax"
)

type Parser struct {
	s              *scanner.Scanner
	module         string
	currentToken   *scanner.Token
	lookAheadToken *scanner.Token
}

func New(s *scanner.Scanner) *Parser {
	p := new(Parser)
	p.s = s
	p.module = "parser"
	p.currentToken = scanner.NewToken()
	p.lookAheadToken = scanner.NewToken()
	return p
}

func (p *Parser) Parse() syntax.SyntaxTree {
	return p.parseChunk()
}

func (p *Parser) nextToken() *scanner.Token {
	if p.lookAheadToken.Category != scanner.TokenEOF {
		p.currentToken = p.lookAheadToken.Clone()
		p.lookAheadToken.Category = scanner.TokenEOF
	} else {
		p.currentToken = p.s.Scan()
	}
	return p.currentToken
}

func (p *Parser) lookAhead() *scanner.Token {
	if p.lookAheadToken.Category == scanner.TokenEOF {
		p.lookAheadToken = p.s.Scan()
	}
	return p.lookAheadToken
}

func (p *Parser) parseChunk() syntax.SyntaxTree {
	block := p.parseBlock()
	if p.nextToken().Category != scanner.TokenEOF {
		panic(&Error{module: p.module, token: p.currentToken,
			str: "expect <eof>"})
	}
	return &syntax.Chunk{Block: block}
}

func (p *Parser) parseBlock() syntax.SyntaxTree {
	block := &syntax.Block{}
	for p.lookAhead().Category != scanner.TokenEOF &&
		p.lookAhead().Category != scanner.TokenEnd &&
		p.lookAhead().Category != scanner.TokenElseif &&
		p.lookAhead().Category != scanner.TokenElse {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
	}
	return block
}

func (p *Parser) parseStatement() syntax.SyntaxTree {
	switch p.lookAhead().Category {
	case ";":
		p.nextToken()
	case scanner.TokenDo:
		p.parseDoStatement()
	case scanner.TokenWhile:
		p.parseWhileStatement()
	case scanner.TokenIf:
		p.parseIfStatement()
	case scanner.TokenLocal:
		p.parseLocalStatement()
	default:
		p.parseOtherStatement()
	}

	return nil
}

func (p *Parser) parseDoStatement() syntax.SyntaxTree {
	p.nextToken()
	assert(p.currentToken.Category == scanner.TokenDo, "not a do statement")
	block := p.parseBlock()
	if p.nextToken().Category != scanner.TokenEnd {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect 'end' for 'do' statement",
		})
	}
	return &syntax.DoStatement{Block: block}
}

func (p *Parser) parseWhileStatement() syntax.SyntaxTree {
	p.nextToken()
	assert(p.currentToken.Category == scanner.TokenWhile,
		"not a while statement")
	exp := p.parseExp()
	if p.nextToken().Category != scanner.TokenDo {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect 'do' for 'while' statement",
		})
	}
	block := p.parseBlock()
	if p.nextToken().Category != scanner.TokenEnd {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect 'end' for 'while' statement",
		})
	}
	return &syntax.WhileStatement{
		Exp:   exp,
		Block: block,
	}
}

func (p *Parser) parseIfStatement() syntax.SyntaxTree {
	p.nextToken()
	assert(p.currentToken.Category == scanner.TokenIf, "not a if statement")
	exp := p.parseExp()
	if p.nextToken().Category != scanner.TokenThen {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect 'then' for 'if' statement",
		})
	}
	trueBranch := p.parseBlock()
	falseBranch := p.parseFalseBranchStatement()
	return &syntax.IfStatement{
		Exp:         exp,
		TrueBranch:  trueBranch,
		FalseBranch: falseBranch,
	}
}

func (p *Parser) parseFalseBranchStatement() syntax.SyntaxTree {
	if p.lookAhead().Category == scanner.TokenElseif {
		return p.parseElseifStatement()
	} else if p.lookAhead().Category == scanner.TokenElse {
		return p.parseElseStatement()
	} else if p.lookAhead().Category == scanner.TokenEnd {
		p.nextToken()
	} else {
		panic(&Error{
			module: p.module,
			token:  p.lookAheadToken,
			str:    "expect 'end' for 'if' statement",
		})
	}
	return nil
}

func (p *Parser) parseElseifStatement() syntax.SyntaxTree {
	p.nextToken()
	assert(p.currentToken.Category == scanner.TokenElseif,
		"not a 'elseif' statement")
	exp := p.parseExp()
	if p.nextToken().Category != scanner.TokenThen {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect 'then' for 'elseif' statement",
		})
	}
	trueBranch := p.parseBlock()
	falseBranch := p.parseFalseBranchStatement()
	return &syntax.ElseifStatement{
		Exp:         exp,
		TrueBranch:  trueBranch,
		FalseBranch: falseBranch,
	}
}

func (p *Parser) parseElseStatement() syntax.SyntaxTree {
	p.nextToken()
	assert(p.currentToken.Category == scanner.TokenElse,
		"not a 'else' statement")
	block := p.parseBlock()
	if p.nextToken().Category != scanner.TokenEnd {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect 'end' for 'else' statement",
		})
	}
	return &syntax.ElseStatement{Block: block}
}

func (p *Parser) parseLocalStatement() syntax.SyntaxTree {
	p.nextToken()
	assert(p.currentToken.Category == scanner.TokenLocal,
		"not a local statement")
	if p.lookAhead().Category == scanner.TokenID {
		return p.parseLocalNameList()
	} else {
		panic(&Error{
			module: p.module,
			token:  p.lookAheadToken,
			str:    "unexpect token after 'local'",
		})
	}
}

func (p *Parser) parseLocalNameList() syntax.SyntaxTree {
	nameList := p.parseNameList()
	var expList syntax.SyntaxTree
	if p.lookAhead().Category == scanner.TokenAssign {
		p.nextToken()
		expList = p.parseExpList()
	}
	return &syntax.LocalNameListStatement{
		NameList: nameList,
		ExpList:  expList,
	}
}

func (p *Parser) parseNameList() syntax.SyntaxTree {
	if p.nextToken().Category != scanner.TokenID {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "expect <id>",
		})
	}
	nameList := &syntax.NameList{}
	nameList.Names = append(nameList.Names, p.currentToken.Clone())
	for p.lookAhead().Category == scanner.TokenComma {
		p.nextToken()
		if p.nextToken().Category != scanner.TokenID {
			panic(&Error{
				module: p.module,
				token:  p.currentToken,
				str:    "expect <id> after ','",
			})
		}
		nameList.Names = append(nameList.Names, p.currentToken.Clone())
	}
	return nameList
}

func (p *Parser) parseExpList() syntax.SyntaxTree {
	expList := &syntax.ExpressionList{}
	anymore := true
	for anymore {
		expList.ExpList = append(expList.ExpList, p.parseExp())
		if p.lookAhead().Category == scanner.TokenComma {
			p.nextToken()
		} else {
			anymore = false
		}
	}
	return expList
}

func (p *Parser) parseOtherStatement() syntax.SyntaxTree {
	exp := p.parsePrefixExp()
	varList := &syntax.VarList{}
	varList.VarList = append(varList.VarList, exp)
	for p.lookAhead().Category != scanner.TokenAssign {
		if p.lookAhead().Category != scanner.TokenComma {
			panic(&Error{
				module: p.module,
				token:  p.lookAheadToken,
				str:    "expect ',' to split var",
			})
		}
		p.nextToken()
		exp := p.parsePrefixExp()
		varList.VarList = append(varList.VarList, exp)
	}
	p.nextToken()
	expList := p.parseExpList()
	return &syntax.AssignmentStatement{
		VarList: varList,
		ExpList: expList,
	}
}

func (p *Parser) parseExp() syntax.SyntaxTree {
	return p.parseExpImpl(nil, scanner.NewToken(), 0)
}

func (p *Parser) parseExpImpl(left syntax.SyntaxTree, oprand *scanner.Token,
	leftPriority int) syntax.SyntaxTree {
	var exp syntax.SyntaxTree
	p.lookAhead()
	if p.lookAheadToken.Category == scanner.TokenSub ||
		p.lookAheadToken.Category == scanner.TokenLen ||
		p.lookAheadToken.Category == scanner.TokenNot {
		p.nextToken()
		exp = &syntax.UnaryExpression{
			OpToken: p.currentToken,
			Exp:     p.parseExpImpl(nil, scanner.NewToken(), 90),
		}
	} else if isMainExp(p.lookAheadToken) {
		exp = p.parseMainExp()
	} else {
		panic(&Error{
			module: p.module,
			token:  p.lookAheadToken,
			str:    "unexpect token for exp",
		})
	}
	for {
		rightPriority := operatorPriority(p.lookAhead())
		if leftPriority < rightPriority {
			exp = p.parseExpImpl(exp, p.nextToken().Clone(), rightPriority)
		} else if leftPriority == rightPriority {
			if leftPriority == 0 {
				return exp
			}
			assert(left != nil, "left operator is nil")
			exp = &syntax.BinaryExpression{
				Left:    left,
				Right:   exp,
				OpToken: oprand.Clone(),
			}
			oprand = scanner.NewToken()
			leftPriority = 0
			exp = p.parseExpImpl(exp, p.nextToken().Clone(), rightPriority)
		} else {
			if left != nil {
				exp = &syntax.BinaryExpression{
					Left:    left,
					Right:   exp,
					OpToken: oprand.Clone(),
				}
			}
			return exp
		}
	}
}

func (p *Parser) parseMainExp() syntax.SyntaxTree {
	var exp syntax.SyntaxTree
	switch p.lookAhead().Category {
	case scanner.TokenNil, scanner.TokenFalse, scanner.TokenTrue,
		scanner.TokenNumber, scanner.TokenString:
		exp = &syntax.Terminator{Token: p.nextToken().Clone()}
	case scanner.TokenID, scanner.TokenLeftParen:
		exp = p.parsePrefixExp()
	default:
		panic(&Error{
			module: p.module,
			token:  p.lookAheadToken,
			str:    "unexpect token for expression",
		})
	}
	return exp
}

func (p *Parser) parsePrefixExp() syntax.SyntaxTree {
	p.nextToken()
	if p.currentToken.Category != scanner.TokenID &&
		p.currentToken.Category != scanner.TokenLeftParen {
		panic(&Error{
			module: p.module,
			token:  p.currentToken,
			str:    "unexpect token here",
		})
	}
	var exp syntax.SyntaxTree
	if p.currentToken.Category == scanner.TokenLeftParen {
		exp = p.parseExp()
		if p.nextToken().Category != scanner.TokenRightParen {
			panic(&Error{
				module: p.module,
				token:  p.currentToken,
				str:    "expect ')'",
			})
		}
	} else {
		exp = &syntax.Terminator{Token: p.currentToken.Clone()}
	}
	return exp
}

func isMainExp(token *scanner.Token) bool {
	return token.Category == scanner.TokenNil ||
		token.Category == scanner.TokenFalse ||
		token.Category == scanner.TokenTrue ||
		token.Category == scanner.TokenNumber ||
		token.Category == scanner.TokenString ||
		token.Category == scanner.TokenID ||
		token.Category == scanner.TokenLeftParen
}

func operatorPriority(token *scanner.Token) int {
	switch token.Category {
	case scanner.TokenDiv, scanner.TokenMul:
		return 80
	case scanner.TokenAdd, scanner.TokenSub:
		return 70
	case scanner.TokenConcat:
		return 60
	case scanner.TokenGreater, scanner.TokenLess,
		scanner.TokenGreaterEqual, scanner.TokenLessEqual,
		scanner.TokenNotEqual, scanner.TokenEqual:
		return 50
	case scanner.TokenAnd:
		return 40
	case scanner.TokenOr:
		return 30
	default:
		return 0
	}
}
