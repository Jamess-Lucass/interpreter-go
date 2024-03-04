package parser

import (
	"testing"

	"github.com/Jamess-Lucass/interpreter-go/ast"
	"github.com/Jamess-Lucass/interpreter-go/lexer"
	"github.com/stretchr/testify/assert"
)

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
