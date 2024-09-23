package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ekediala/interpreter/lexer"
	"github.com/ekediala/interpreter/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.ReadAndAdvanceToken(); tok.Type != token.EOF; tok = l.ReadAndAdvanceToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
