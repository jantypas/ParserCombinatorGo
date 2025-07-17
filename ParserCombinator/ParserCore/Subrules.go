package ParserCore

import (
	"fmt"
	"strings"
)

func parseAnyString(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == STRING {
		if opt&PARSE_OPTION_CONVER_TO_UPPERCASE != 0 {
			// Convert to uppercase if the option is set
			tok.Value = strings.ToUpper(tok.Value)
		}
		if opt&PARSE_OPTION_CONVERT_TO_LOWERCASE != 0 {
			// Convert to lowercase if the option is set
			tok.Value = strings.ToLower(tok.Value)
		}
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected STRING, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseAnyInteger(l *Lexer, opt int) (error, int) {
	tok := l.NextToken()
	if tok.Type == INTEGER {
		var value int
		_, err := fmt.Sscanf(tok.Value, "%d", &value)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s at line %d, column %d", tok.Value, tok.Line, tok.Column), 0
		}
		return nil, value
	} else {
		return fmt.Errorf("expected INTEGER, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), 0
	}
}

func parseAnyFloat(l *Lexer, opt int) (error, float64) {
	tok := l.NextToken()
	if tok.Type == FLOAT {
		var value float64
		_, err := fmt.Sscanf(tok.Value, "%f", &value)
		if err != nil {
			return fmt.Errorf("invalid float value: %s at line %d, column %d", tok.Value, tok.Line, tok.Column), 0.0
		}
		return nil, value
	} else {
		return fmt.Errorf("expected FLOAT, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), 0.0
	}
}

func parseAnyQuotedString(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == QUOTED_STRING {
		if opt&PARSE_OPTION_CONVER_TO_UPPERCASE != 0 {
			// Convert to uppercase if the option is set
			tok.Value = strings.ToUpper(tok.Value)
		}
		if opt&PARSE_OPTION_CONVERT_TO_LOWERCASE != 0 {
			// Convert to lowercase if the option is set
			tok.Value = strings.ToLower(tok.Value)
		}
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected QUOTED_STRING, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseComma(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == COMMA {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected COMMA, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseColon(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == COLON {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected COLON, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseStringChoice(l *Lexer, choices []string, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == STRING {
		value := tok.Value
		if opt&PARSE_OPTION_CONVER_TO_UPPERCASE != 0 {
			value = strings.ToUpper(value)
		}
		if opt&PARSE_OPTION_CONVERT_TO_LOWERCASE != 0 {
			value = strings.ToLower(value)
		}
		for _, choice := range choices {
			if value == choice {
				return nil, value
			}
		}
		return fmt.Errorf("expected one of %v, got %s at line %d, column %d", choices, tok.Value, tok.Line, tok.Column), ""
	} else {
		return fmt.Errorf("expected STRING, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseStringList(l *Lexer, sl []string, opt int) (error, []string) {
	var value string
	for sitem := range sl {
		tok := l.NextToken()
		if tok.Type == STRING {
			if opt&PARSE_OPTION_CONVER_TO_UPPERCASE != 0 {
				value = strings.ToUpper(tok.Value)
			}
			if opt&PARSE_OPTION_CONVERT_TO_LOWERCASE != 0 {
				value = strings.ToLower(tok.Value)
			}
			if value != sl[sitem] {
				return fmt.Errorf("Expected one of %v, got %s at line %d, column %d", sl, tok.Value, tok.Line, tok.Column), nil
			}
		}
	}
	return nil, sl
}

func parseQuestion(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == QUESTION {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected QUESTION, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseLessThan(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == LESS_THAN {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected LESS_THAN, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseGreaterThan(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == GREATER_THAN {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected GREATER_THAN, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseExclamation(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == EXCLAMATION {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected EXCLAMATION, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseEqual(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == EQUAL {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected EQUALS, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parsePlus(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == PLUS {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected PLUS, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parsePercent(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == PERCENT {
		return nil, tok.Value
	} else {
		return fmt.Errorf("expected PERCENT, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}
