package ParserCore

import (
	"testing"
)

func Test_convertString(t *testing.T) {
	result := convertString("Hello World", PARSE_OPTION_CONVERT_TO_UPPERCASE)
	if result != "HELLO WORLD" {
		t.Errorf("convertString() failed, expected 'HELLO WORLD', got '%s'", result)
	}
	result = convertString("Hello World", PARSE_OPTION_CONVERT_TO_LOWERCASE)
	if result != "hello world" {
		t.Errorf("convertString() failed, expected 'hello world', got '%s'", result)
	}
	result = convertString("Hello World", 0)
	if result != "Hello World" {
		t.Errorf("convertString() failed, expected 'Hello World', got '%s'", result)
	}
}

func Test_parseAnyString(t *testing.T) {
	l := NewLexer("Hello World", nil)
	err, value := parseAnyString(l, PARSE_OPTION_CONVERT_TO_UPPERCASE)
	if err != nil || value != "HELLO" {
		t.Errorf("parseAnyString() failed, expected 'HELLO WORLD', got '%s' with error '%v'", value, err)
	}

	l = NewLexer("123", nil)
	err, value = parseAnyString(l, 0)
	if err == nil || value != "" {
		t.Errorf("parseAnyString() failed, expected error for non-string input, got '%s' with error '%v'", value, err)
	}
}

func Test_parseAnyInteger(t *testing.T) {
	l := NewLexer("123", nil)
	err, value := parseAnyInteger(l, PARSE_OPTION_CONVERT_TO_UPPERCASE)
	if err != nil || value != 123 {
		t.Errorf("parseAnyInteger() failed, expected 123, got '%d' with error '%v'", value, err)
	}

	l = NewLexer("BOB", nil)
	err, value = parseAnyInteger(l, 0)
	if err == nil || value != 0 {
		t.Errorf("parseAnyInteger() failed, expected error for non-input input, got '%d' with error '%v'", value, err)
	}
}

func Test_parseAnyFloat(t *testing.T) {
	l := NewLexer("123.123", nil)
	err, value := parseAnyFloat(l, PARSE_OPTION_CONVERT_TO_UPPERCASE)
	if err != nil || value != 123.123 {
		t.Errorf("parseAnyFloat() failed, expected 123.123, got '%f' with error '%v'", value, err)
	}

	l = NewLexer("BOB", nil)
	err, value = parseAnyFloat(l, 0)
	if err == nil || value != 0 {
		t.Errorf("parseAnyFloat() failed, expected error for non-input input, got '%f' with error '%v'", value, err)
	}
}

func Test_parseAnyQuotedString(t *testing.T) {
	l := NewLexer("\"Hello World\"", nil)
	err, value := parseAnyQuotedString(l, PARSE_OPTION_CONVERT_TO_UPPERCASE)
	if err != nil || value != "\"HELLO WORLD\"" {
		t.Errorf("parseAnyQuotedString() failed, expected 'HELLO WORLD', got '%s' with error '%v'", value, err)
	}

	l = NewLexer("123", nil)
	err, value = parseAnyQuotedString(l, 0)
	if err == nil || value != "" {
		t.Errorf("parseAnyQuotedString() failed, expected error for non-quoted input, got '%s' with error '%v'", value, err)
	}
}

func Test_parseComma(t *testing.T) {
	l := NewLexer(",", nil)
	err, value := parseComma(l, 0)
	if err != nil || value != "," {
		t.Errorf("parseComma() failed, expected ',', got '%s' with error '%v'", value, err)
	}
}

func Test_parseColon(t *testing.T) {
	l := NewLexer(":", nil)
	err, value := parseColon(l, 0)
	if err != nil || value != ":" {
		t.Errorf("parseColon() failed, expected ':', got '%s' with error '%v'", value, err)
	}
}

func Test_parseStringChoice(t *testing.T) {
	l := NewLexer("choice1", nil)
	choices := []string{"choice1", "choice2", "choice3"}
	err, value := parseStringChoice(l, choices, 0)
	if err != nil || value != "choice1" {
		t.Errorf("parseStringChoice() failed, expected 'choice1', got '%s' with error '%v'", value, err)
	}
	l = NewLexer("invalid_choice", nil)
	choices = []string{"alpha", "beta", "gamma"}
	err, value = parseStringChoice(l, choices, PARSE_OPTION_STRING_IS_OPTIONAL)
	if err != nil || value != "" {
		t.Errorf("parseStringChoice() failed, expected error for invalid choice, got '%s' with error '%v'", value, err)
	}
}

func Test_parseStringList(t *testing.T) {
	l := NewLexer("item1 item2 item3", nil)
	err, value := parseStringList(l, []string{"item1", "item2", "item3"}, 0)
	if err != nil {
		t.Errorf("parseStringList() failed, expected ['item1', 'item2', 'item3'], got '%v' with error '%v'", value, err)
	}
}

func Test_parseQuestion(t *testing.T) {
	l := NewLexer("?", nil)
	err, value := parseQuestion(l, 0)
	if err != nil || value != "?" {
		t.Errorf("parseQuestion() failed, expected '?', got '%s' with error '%v'", value, err)
	}
}

func Test_parseLessThan(t *testing.T) {
	l := NewLexer("<", nil)
	err, value := parseLessThan(l, 0)
	if err != nil || value != "<" {
		t.Errorf("parseLessThan() failed, expected '<', got '%s' with error '%v'", value, err)
	}
}

func Test_parseGreaterThan(t *testing.T) {
	l := NewLexer(">", nil)
	err, value := parseGreaterThan(l, 0)
	if err != nil || value != ">" {
		t.Errorf("parseGreaterThan() failed, expected '>', got '%s' with error '%v'", value, err)
	}
}

func Test_parseExclamation(t *testing.T) {
	l := NewLexer("!", nil)
	err, value := parseExclamation(l, 0)
	if err != nil || value != "!" {
		t.Errorf("parseExclamation() failed, expected '!', got '%s' with error '%v'", value, err)
	}
}

func Test_parsePlus(t *testing.T) {
	l := NewLexer("+", nil)
	err, value := parsePlus(l, 0)
	if err != nil || value != "+" {
		t.Errorf("parsePlus() failed, expected '+', got '%s' with error '%v'", value, err)
	}
}

func Test_parsePercent(t *testing.T) {
	l := NewLexer("%", nil)
	err, value := parsePercent(l, 0)
	if err != nil || value != "%" {
		t.Errorf("parsePercent() failed, expected percent, got '%s' with error '%v'", value, err)
	}
}

func Test_parseEqual(t *testing.T) {
	l := NewLexer("=", nil)
	err, value := parseEqual(l, 0)
	if err != nil || value != "=" {
		t.Errorf("parseEqual() failed, expected '=', got '%s' with error '%v'", value, err)
	}
}
