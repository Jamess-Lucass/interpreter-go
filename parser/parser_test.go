package parser

import (
	"fmt"
	"testing"

	"github.com/Jamess-Lucass/interpreter-go/ast"
	"github.com/Jamess-Lucass/interpreter-go/lexer"
	"github.com/stretchr/testify/assert"
)

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) {
	literal, ok := expression.(*ast.IntegerLiteral)
	assert.True(t, ok)

	assert.Equal(t, value, literal.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), literal.TokenLiteral())
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) {
	identifier, ok := expression.(*ast.Identifier)
	assert.True(t, ok)

	assert.Equal(t, value, identifier.Value)
	assert.Equal(t, value, identifier.TokenLiteral())
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) {
	identifier, ok := expression.(*ast.Boolean)
	assert.True(t, ok)

	assert.Equal(t, value, identifier.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), identifier.TokenLiteral())
}

func testliteralExpression(t *testing.T, expression ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, expression, int64(v))
	case int64:
		testIntegerLiteral(t, expression, v)
	case string:
		testIdentifier(t, expression, v)
	case bool:
		testBooleanLiteral(t, expression, v)
	default:
		t.Errorf("type of expression not handled. got=%T", expression)
	}
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) {
	infixExpression, ok := expression.(*ast.InfixExpression)
	assert.True(t, ok)

	testliteralExpression(t, infixExpression.Left, left)

	assert.Equal(t, infixExpression.Operator, operator)

	testliteralExpression(t, infixExpression.Right, right)
}

func Test_ParsingLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		statement := program.Statements[i]

		letStatement, ok := statement.(*ast.LetStatement)
		assert.True(t, ok)

		assert.Equal(t, "let", statement.TokenLiteral())
		assert.Equal(t, test.expectedIdentifier, letStatement.Name.Value)
		assert.Equal(t, test.expectedIdentifier, letStatement.Name.TokenLiteral())
	}
}

func Test_ParsingReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 838383;`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 3)

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		assert.True(t, ok)

		assert.Equal(t, "return", returnStatement.TokenLiteral())
	}
}

func Test_IdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	testIdentifier(t, stmt.Expression, "foobar")
}

func Test_BooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	testBooleanLiteral(t, stmt.Expression, true)
}

func Test_IntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	testIntegerLiteral(t, stmt.Expression, int64(5))
}

func Test_PrefixExpression(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)
		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		expression, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok)

		assert.Equal(t, test.operator, expression.Operator)

		testliteralExpression(t, expression.Right, test.integerValue)
	}
}

func Test_InfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)
		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		testInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue)
	}
}

func Test_OperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)

		assert.Equal(t, test.expected, program.String())
	}
}

func Test_IfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	expression, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	testInfixExpression(t, expression.Condition, "x", "<", "y")

	assert.Len(t, expression.Consequence.Statements, 1)

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIdentifier(t, consequence.Expression, "x")

	assert.Nil(t, expression.Alternative)
}

func Test_IfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	expression, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	testInfixExpression(t, expression.Condition, "x", "<", "y")

	assert.Len(t, expression.Consequence.Statements, 1)

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIdentifier(t, consequence.Expression, "x")

	assert.Len(t, expression.Alternative.Statements, 1)

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIdentifier(t, alternative.Expression, "y")
}
