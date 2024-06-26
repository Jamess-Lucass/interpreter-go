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

func testStringLiteral(t *testing.T, expression ast.Expression, value string) {
	literal, ok := expression.(*ast.StringLiteral)
	assert.True(t, ok)

	assert.Equal(t, value, literal.Value)
	assert.Equal(t, value, literal.TokenLiteral())
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

func testLetStatement(t *testing.T, statement ast.Statement, name string) {
	letStatement, ok := statement.(*ast.LetStatement)
	assert.True(t, ok)

	assert.Equal(t, "let", statement.TokenLiteral())
	assert.Equal(t, name, letStatement.Name.Value)
	assert.Equal(t, name, letStatement.Name.TokenLiteral())
}

func Test_ParsingLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)
		assert.Len(t, program.Statements, 1)

		statement := program.Statements[0]
		testLetStatement(t, statement, test.expectedIdentifier)

		value := statement.(*ast.LetStatement).Value
		testliteralExpression(t, value, test.expectedValue)
	}
}

func Test_ParsingReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y", "y"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)
		assert.Len(t, program.Statements, 1)

		returnStatement, ok := program.Statements[0].(*ast.ReturnStatement)
		assert.True(t, ok)
		assert.Equal(t, "return", returnStatement.TokenLiteral())

		testliteralExpression(t, returnStatement.Value, test.expectedValue)
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
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
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

func Test_FunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	expression, ok := stmt.Expression.(*ast.FunctionLiteral)
	assert.True(t, ok)
	assert.Len(t, expression.Parameters, 2)

	testliteralExpression(t, expression.Parameters[0], "x")
	testliteralExpression(t, expression.Parameters[1], "y")

	assert.Len(t, expression.Body.Statements, 1)

	bodyStmt, ok := expression.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func Test_FunctionParameter(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)
		assert.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		function, ok := statement.Expression.(*ast.FunctionLiteral)
		assert.True(t, ok)

		assert.Len(t, function.Parameters, len(test.expectedParams))

		for i, identifier := range test.expectedParams {
			testliteralExpression(t, function.Parameters[i], identifier)
		}
	}
}

func Test_CallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	expression, ok := statement.Expression.(*ast.CallExpression)
	assert.True(t, ok)
	assert.Len(t, expression.Arguments, 3)

	testIdentifier(t, expression.Function, "add")
	testliteralExpression(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func Test_CallExpressionParameter(t *testing.T) {
	tests := []struct {
		input             string
		expectedArguments []string
	}{
		{input: "add()", expectedArguments: []string{}},
		// {input: "add(1)", expectedArguments: []string{"1"}},
		{input: "add(1, 2 * 3, 4 + 5)", expectedArguments: []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)

		program := p.Parse()

		assert.Len(t, p.errors, 0)
		assert.NotNil(t, program)
		assert.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		function, ok := statement.Expression.(*ast.CallExpression)
		assert.True(t, ok)

		assert.Len(t, function.Arguments, len(test.expectedArguments))

		for i, argument := range test.expectedArguments {
			assert.Equal(t, argument, function.Arguments[i].String())
		}
	}
}

func Test_StringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	testStringLiteral(t, statement.Expression, "hello world")
}

func Test_ParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	literal, ok := statement.Expression.(*ast.ArrayLiteral)
	assert.True(t, ok)
	assert.Len(t, literal.Elements, 3)

	testIntegerLiteral(t, literal.Elements[0], 1)
	testInfixExpression(t, literal.Elements[1], 2, "*", 2)
	testInfixExpression(t, literal.Elements[2], 3, "+", 3)
}

func Test_ParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	literal, ok := statement.Expression.(*ast.IndexExpressopn)
	assert.True(t, ok)

	testIdentifier(t, literal.Left, "myArray")
	testInfixExpression(t, literal.Index, 1, "+", 1)
}

func Test_ParsingHashLiteralStringKey(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	hashLiteral, ok := statement.Expression.(*ast.HashLiteral)
	assert.True(t, ok)
	assert.Len(t, hashLiteral.Pairs, 3)

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		assert.True(t, ok)

		testIntegerLiteral(t, value, expected[literal.String()])
	}
}

func Test_ParsingHashLiteralIntegerKey(t *testing.T) {
	input := `{1: "one", 2: "two", 3: "three"}`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	hashLiteral, ok := statement.Expression.(*ast.HashLiteral)
	assert.True(t, ok)
	assert.Len(t, hashLiteral.Pairs, 3)

	expected := map[int64]string{
		1: "one",
		2: "two",
		3: "three",
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.IntegerLiteral)
		assert.True(t, ok)

		testStringLiteral(t, value, expected[literal.Value])
	}
}

func Test_ParsingHashLiteralBooleanKey(t *testing.T) {
	input := `{true: "one", false: "two"}`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	hashLiteral, ok := statement.Expression.(*ast.HashLiteral)
	assert.True(t, ok)
	assert.Len(t, hashLiteral.Pairs, 2)

	expected := map[bool]string{
		true:  "one",
		false: "two",
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.Boolean)
		assert.True(t, ok)

		testStringLiteral(t, value, expected[literal.Value])
	}
}

func Test_ParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	hashLiteral, ok := statement.Expression.(*ast.HashLiteral)
	assert.True(t, ok)
	assert.Len(t, hashLiteral.Pairs, 0)
}

func Test_ParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.Parse()

	assert.Len(t, p.errors, 0)
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	hashLiteral, ok := statement.Expression.(*ast.HashLiteral)
	assert.True(t, ok)
	assert.Len(t, hashLiteral.Pairs, 3)

	expected := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		assert.True(t, ok)

		expectedFunc, ok := expected[literal.String()]
		assert.True(t, ok)

		expectedFunc(value)
	}
}
