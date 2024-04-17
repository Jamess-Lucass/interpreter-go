package evaluator

import (
	"fmt"
	"testing"

	"github.com/Jamess-Lucass/interpreter-go/lexer"
	"github.com/Jamess-Lucass/interpreter-go/object"
	"github.com/Jamess-Lucass/interpreter-go/parser"
	"github.com/stretchr/testify/assert"
)

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Parse()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func Test_EvalIntgerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Integer)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value)
	}
}

func Test_EvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Boolean)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value, fmt.Sprintf("for input: %v", test.input))
	}
}

func Test_EvalIntegerOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Integer)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value)
	}
}

func Test_BangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Boolean)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value)
	}
}

func Test_IfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)

		integer, ok := test.expected.(int)
		if ok {
			result, ok := evaluated.(*object.Integer)
			assert.True(t, ok)

			assert.Equal(t, int64(integer), result.Value)
		} else {
			assert.Equal(t, NULL, evaluated)
		}
	}
}
func Test_ReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
		`,
			10,
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Integer)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value)
	}
}

func Test_ErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Error)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Message)
	}
}

func Test_LetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Integer)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value)
	}
}

func Test_FunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Function)
	assert.True(t, ok)

	assert.Len(t, result.Parameters, 1)
	assert.Equal(t, result.Parameters[0].String(), "x")
	assert.Equal(t, "(x + 2)", result.Body.String())
}

func Test_FunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		result, ok := evaluated.(*object.Integer)
		assert.True(t, ok)

		assert.Equal(t, test.expected, result.Value)
	}
}

func Test_Closures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Integer)
	assert.True(t, ok)

	assert.Equal(t, int64(4), result.Value)
}
