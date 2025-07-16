package ParserCore

// Lexer is a simple lexical analyzer that tokenizes input strings.
// It currently understands a limited set of token types including:
// - STRING: Unquoted strings
// - QUOTED_STRING: Strings enclosed in double quotes
// - INTEGER: Whole numbers
// - FLOAT: Decimal numbers
// - COMMA: Comma character
// - COLON: Colon character
// - ERROR: Represents an error in tokenization
// - EOF: End of file marker

import (
	"unicode"
)

type TokenType int

const (
	ERROR TokenType = iota
	EOF
	STRING
	QUOTED_STRING
	INTEGER
	FLOAT
	COMMA
	COLON
)

// The names of the token types for easy reference.
var TokenTypeNames = []string{
	"ERROR_TYPE",
	"EOF",
	"STRING",
	"QUOTED_STRING",
	"INTEGER",
	"FLOAT",
	"COMMA",
	"COLON",
}

// Token represents a single token in the input string.
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// The core lexer object iself
type Lexer struct {
	input  string
	pos    int
	line   int
	column int
}

// NewLexer creates a new Lexer instance with the provided input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

// The workhorse of the system -- NextToken reads the next token from the input string.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: EOF, Line: l.line, Column: l.column}
	}

	switch {
	case l.input[l.pos] == '"':
		return l.readQuotedString()
	case l.input[l.pos] == ',':
		token := Token{Type: COMMA, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == ':':
		token := Token{Type: COLON, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '-':
		return l.readNumber()
	case unicode.IsLetter(rune(l.input[l.pos])):
		return l.readString()
	default:
		return Token{Type: ERROR, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
	}
}

func (l *Lexer) readQuotedString() Token {
	startPos := l.pos
	startColumn := l.column
	l.pos++ // Skip opening quote
	l.column++

	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}

	if l.pos >= len(l.input) {
		return Token{Type: ERROR, Value: "Unterminated string", Line: l.line, Column: startColumn}
	}

	l.pos++ // Skip closing quote
	l.column++

	return Token{
		Type:   QUOTED_STRING,
		Value:  l.input[startPos:l.pos],
		Line:   l.line,
		Column: startColumn,
	}
}

func (l *Lexer) readNumber() Token {
	startPos := l.pos
	startColumn := l.column
	isFloat := false

	// Handle negative sign
	if l.input[l.pos] == '-' {
		l.pos++
		l.column++
		if l.pos >= len(l.input) || !unicode.IsDigit(rune(l.input[l.pos])) {
			return Token{Type: ERROR, Value: "Invalid number", Line: l.line, Column: startColumn}
		}
	}

	for l.pos < len(l.input) {
		if l.input[l.pos] == '.' && !isFloat {
			isFloat = true
			l.pos++
			l.column++
			continue
		}
		if !unicode.IsDigit(rune(l.input[l.pos])) {
			break
		}
		l.pos++
		l.column++
	}

	tokenType := INTEGER
	if isFloat {
		tokenType = FLOAT
	}

	return Token{
		Type:   tokenType,
		Value:  l.input[startPos:l.pos],
		Line:   l.line,
		Column: startColumn,
	}
}

func (l *Lexer) readString() Token {
	startPos := l.pos
	startColumn := l.column

	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos]))) {
		l.pos++
		l.column++
	}

	return Token{
		Type:   STRING,
		Value:  l.input[startPos:l.pos],
		Line:   l.line,
		Column: startColumn,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
			l.pos++
		} else if unicode.IsSpace(rune(l.input[l.pos])) {
			l.column++
			l.pos++
		} else {
			break
		}
	}
}
