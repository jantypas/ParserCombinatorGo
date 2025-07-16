package ParserCore

import (
	"fmt"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	lex := NewLexer(`"Hello, World!" -123 -45.67 MOVE, RUN: "Test String"`)
	// while lex.NextToken().Type != EOF {
	for {
		tok := lex.NextToken()
		if tok.Type == EOF {
			break
		}
		fmt.Printf("Token Type %-20s Token Value %s\n", TokenTypeNames[tok.Type], tok.Value)
	}
}
