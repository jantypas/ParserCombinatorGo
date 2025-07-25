/*
 * Copyright 2025 John Antypas
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ParserCore

import (
	"fmt"
	"strings"
)

func convertString(s string, opt int) string {
	if opt&PARSE_OPTION_CONVERT_TO_UPPERCASE != 0 {
		return strings.ToUpper(s)
	}
	if opt&PARSE_OPTION_CONVERT_TO_LOWERCASE != 0 {
		return strings.ToLower(s)
	}
	return s
}

func parseAnyString(l *Lexer, opt int) (error, string) {
	tok := l.NextToken()
	if tok.Type == STRING {
		value := convertString(tok.Value, opt)
		return nil, value
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
		value := convertString(tok.Value, opt)
		return nil, value
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
		value = convertString(value, opt)
		for _, choice := range choices {
			if value == choice {
				return nil, value
			}
		}
		// If we reach here, the value is not in the choices
		if opt&PARSE_OPTION_STRING_IS_OPTIONAL != 0 {
			l.PushBack(tok)
			return nil, "" // Return empty string if the option is optional
		} else {
			return fmt.Errorf("expected one of %v, got %s at line %d, column %d", choices, tok.Value, tok.Line, tok.Column), ""
		}
	} else {
		return fmt.Errorf("expected STRING, got %s at line %d, column %d", tok.Value, tok.Line, tok.Column), ""
	}
}

func parseStringList(l *Lexer, sl []string, opt int) (error, []string) {
	for sitem := range sl {
		tok := l.NextToken()
		if tok.Type == STRING {
			value := convertString(tok.Value, opt)
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
