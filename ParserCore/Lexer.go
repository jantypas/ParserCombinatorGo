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
	"strings"
	"unicode"
)

var LexerVersion = "1.0.0"

type TokenType int

// TokenType represents the different types of tokens that can be recognized by the lexer.
const (
	ERROR TokenType = iota
	EOF
	STRING
	QUOTED_STRING
	INTEGER
	FLOAT
	COMMA
	COLON
	QUESTION
	LESS_THAN
	GREATER_THAN
	EXCLAMATION
	PLUS
	PERCENT
	EQUAL
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
	"QUESTION",
	"LESS_THAN",
	"GREATER_THAN",
	"EXCLAMATION",
	"PLUS",
	"PERCENT",
	"EQUAL",
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
	input          string
	pos            int
	line           int
	column         int
	ignoredStrings []string
	lastToken      *Token // Added field to store last token
}

// NewLexer creates a new Lexer instance with the provided input string.
func NewLexer(input string, exclude []string) *Lexer {
	return &Lexer{
		input:          input,
		pos:            0,
		line:           1,
		column:         1,
		ignoredStrings: exclude,
		lastToken:      nil,
	}
}

// SetIgnoredStrings sets the list of strings to be ignored during tokenization
func (l *Lexer) SetIgnoredStrings(ignored []string) {
	l.ignoredStrings = ignored
}

// PushBack pushes back the last read token so it can be read again
func (l *Lexer) PushBack(token Token) {
	l.lastToken = &token
}

// The workhorse of the system -- NextToken reads the next token from the input string.
func (l *Lexer) NextToken() Token {
	if l.lastToken != nil {
		token := *l.lastToken
		l.lastToken = nil
		return token
	}

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
	case l.input[l.pos] == '?':
		token := Token{Type: QUESTION, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == '<':
		token := Token{Type: LESS_THAN, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == '>':
		token := Token{Type: GREATER_THAN, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == '!':
		token := Token{Type: EXCLAMATION, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == '+':
		token := Token{Type: PLUS, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == '%':
		token := Token{Type: PERCENT, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case l.input[l.pos] == '=':
		token := Token{Type: EQUAL, Value: string(l.input[l.pos]), Line: l.line, Column: l.column}
		l.pos++
		l.column++
		return token
	case unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '-':
		return l.readNumber()
	case unicode.IsLetter(rune(l.input[l.pos])):
		if token := l.readString(); !l.shouldIgnore(token.Value) {
			return token
		}
		return l.NextToken()
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

func (l *Lexer) shouldIgnore(value string) bool {
	for _, ignored := range l.ignoredStrings {
		if strings.EqualFold(value, ignored) {
			return true
		}
	}
	return false
}
