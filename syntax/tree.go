package syntax

import "github.com/ksco/slua/scanner"

type SyntaxTree interface{}

type (
	Chunk struct {
		Block SyntaxTree
	}

	Block struct {
		Stmts []SyntaxTree
	}

	DoStatement struct {
		Block SyntaxTree
	}

	WhileStatement struct {
		Exp   SyntaxTree
		Block SyntaxTree
	}

	IfStatement struct {
		Exp         SyntaxTree
		TrueBranch  SyntaxTree
		FalseBranch SyntaxTree
	}

	ElseifStatement struct {
		Exp         SyntaxTree
		TrueBranch  SyntaxTree
		FalseBranch SyntaxTree
	}

	ElseStatement struct {
		Block SyntaxTree
	}

	LocalNameListStatement struct {
		NameList SyntaxTree
		ExpList  SyntaxTree
	}

	AssignmentStatement struct {
		VarList SyntaxTree
		ExpList SyntaxTree
	}

	VarList struct {
		VarList []SyntaxTree
	}

	Terminator struct {
		Token *scanner.Token
	}

	BinaryExpression struct {
		Left    SyntaxTree
		Right   SyntaxTree
		OpToken *scanner.Token
	}

	UnaryExpression struct {
		Exp     SyntaxTree
		OpToken *scanner.Token
	}

	NameList struct {
		Names []*scanner.Token
	}

	ExpressionList struct {
		ExpList []SyntaxTree
	}
)
