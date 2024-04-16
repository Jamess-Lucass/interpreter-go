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

	return Eval(program)
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
