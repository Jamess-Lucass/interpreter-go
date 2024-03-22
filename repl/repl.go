package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Jamess-Lucass/interpreter-go/evaluator"
	"github.com/Jamess-Lucass/interpreter-go/lexer"
	"github.com/Jamess-Lucass/interpreter-go/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.Parse()
		if len(p.Errors()) > 0 {
			io.WriteString(out, "Whoops!, we ran into an issue!\n")
			io.WriteString(out, "Parser errors:\n")

			for _, msg := range p.Errors() {
				io.WriteString(out, fmt.Sprintf("\t%s\n", msg))
			}

			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
