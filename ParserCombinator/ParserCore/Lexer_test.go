package ParserCore

import (
	"fmt"
	"testing"
)

func TestLexer_readString(t *testing.T) {
	l := &Lexer{
		input:          "Please Hello < 5% ! 12 123 > = 123.43 -123 ? -123.45 \"This is a test\":,",
		ignoredStrings: []string{"PLEASE", "?"},
	}
	for {
		tok := l.NextToken()
		if tok.Type == EOF {
			break
		}
		if tok.Type == ERROR {
			t.Errorf("Error reading token: %s at line %d, column %d", tok.Value, tok.Line, tok.Column)
			return
		}
		fmt.Printf("%s %s\n", TokenTypeNames[tok.Type], tok.Value)
	}
}
