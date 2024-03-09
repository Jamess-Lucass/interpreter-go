package ast

import (
	"bytes"
	"strings"

	"github.com/Jamess-Lucass/interpreter-go/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // token.IDINT
	Value string
}

var _ Expression = (*Identifier)(nil)

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

var _ Statement = (*LetStatement)(nil)

func (s *LetStatement) statementNode() {}

func (s *LetStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token token.Token // token.RETURN
	Value Expression
}

var _ Statement = (*ReturnStatement)(nil)

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

var _ Statement = (*ExpressionStatement)(nil)

func (s *ExpressionStatement) statementNode() {}

func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

var _ Expression = (*IntegerLiteral)(nil)

func (s *IntegerLiteral) expressionNode() {}

func (s *IntegerLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *IntegerLiteral) String() string {
	return s.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

var _ Expression = (*PrefixExpression)(nil)

func (s *PrefixExpression) expressionNode() {}

func (s *PrefixExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(s.Operator)
	out.WriteString(s.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

var _ Expression = (*InfixExpression)(nil)

func (s *InfixExpression) expressionNode() {}

func (s *InfixExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(s.Left.String())
	out.WriteString(" " + s.Operator + " ")
	out.WriteString(s.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

var _ Expression = (*Boolean)(nil)

func (s *Boolean) expressionNode() {}

func (s *Boolean) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Boolean) String() string {
	return s.Token.Literal
}

type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

var _ Statement = (*BlockStatement)(nil)

func (s *BlockStatement) statementNode() {}

func (s *BlockStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range s.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

var _ Expression = (*IfExpression)(nil)

func (s *IfExpression) expressionNode() {}

func (s *IfExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(s.Condition.String())
	out.WriteString(" ")

	if s.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(s.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

var _ Expression = (*FunctionLiteral)(nil)

func (s *FunctionLiteral) expressionNode() {}

func (s *FunctionLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range s.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(s.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(s.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

var _ Expression = (*CallExpression)(nil)

func (s *CallExpression) expressionNode() {}

func (s *CallExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *CallExpression) String() string {
	var out bytes.Buffer

	arguments := []string{}
	for _, a := range s.Arguments {
		arguments = append(arguments, a.String())
	}

	out.WriteString(s.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(arguments, ", "))
	out.WriteString(")")

	return out.String()
}
