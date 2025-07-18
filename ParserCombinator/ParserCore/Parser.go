package ParserCore

import (
	"fmt"
)

type ParserObject struct {
	Debug   bool // Debug flag to control debug output
	Input   string
	Exclude []string // List of tokens to exclude from parsing
}

// Constants
// These are the types of objects we can parse
const (
	PARSE_ANY_STRING = iota
	PARSE_ANY_INTEGER
	PARSE_ANY_FLOAT
	PARSE_ANY_QUOTED_STRING
	PARSE_COMMA
	PARSE_COLON
	PARSE_STRING_CHOICE
	PARSE_STRING_LIST
	PARSE_QUESTION
	PARSE_LESS_THAN
	PARSE_GREATER_THAN
	PARSE_EXCLAMATION
	PARSE_PLUS
	PARSE_PERCENT
	PARSE_EQUAL
)

// ParserNames is a list of names for the parser types.
var ParserNames = []string{
	"PARSE_ANY_STRING",
	"PARSE_ANY_INTEGER",
	"PARSE_ANY_FLOAT",
	"PARSE_ANY_QUOTED_STRING",
	"PARSE_COMMA",
	"PARSE_COLON",
	"PARSE_STRING_CHOICE",
	"PARSE_STRING_LIST",
	"PARSE_QUESTION",
	"PARSE_LESS_THAN",
	"PARSE_GREATER_THAN",
	"PARSE_EXCLAMATION",
	"PARSE_PLUS",
	"PARSE_PERCENT",
	"PARSE_EQUAL",
}

// When we parse something, here are possible error codes.
const (
	PARSE_RESULT_SUCCESS = iota
	PARSE_RESULT_FAILURE
	PARSE_RESULT_SKIP
)

// Options allow us to modify text before parsing
const (
	PARSE_OPTION_CONVER_TO_UPPERCASE = 1 << iota
	PARSE_OPTION_CONVERT_TO_LOWERCASE
)

// For each of our rules, there are various steps.
// Each step defines a name, the type of object we expect
// any objects we need to use and functions to handle success and failure.
type ParserRuleStep struct {
	Name            string
	ParserType      int
	Options         int
	SkipOnTypeError bool // If true, skip this step if the token type is not as expected
	ParsedValues    []string
	ParseHandler    func(err error, token interface{}, tokType int, data *interface{}) (int, error)
}

// ParserRule defines a rule that consists of multiple steps.
type ParseRule struct {
	Name  string
	Steps []ParserRuleStep
}

func ifDebug(debug bool, f func(format string, a ...interface{}) (int, error), format string, a ...interface{}) {
	if debug {
		f(format, a...)
	}
}

// parseRule processes a single rule with its steps.
// It uses the Lexer to read tokens and applies the ParseHandler for each step.
// If any step fails, it returns an error including a request to skip to the next rule
func parseRule(l *Lexer, rule ParseRule, data *interface{}, debug bool) (int, error) {
	var result int
	// For each step in the rule
	for _, step := range rule.Steps {
		ifDebug(debug, fmt.Printf, "Parse: Trying step: %s for rule: %s\n", step.Name, rule.Name)
		ifDebug(debug, fmt.Printf, "       Expecting token type %s with options %d\n", ParserNames[step.ParserType], step.Options)
		// Call the lexer to get the next token -- requesting a specific type to be decoded
		// If the type is wrong return an error
		// If the type is correct, call the ParseHandler to process the token
		switch step.ParserType {
		case PARSE_ANY_STRING:
			err, value := parseAnyString(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_ANY_STRING, data)
			if err != nil {
				return result, err
			}
		case PARSE_ANY_FLOAT:
			err, value := parseAnyFloat(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}

			result, err = step.ParseHandler(err, value, PARSE_ANY_FLOAT, data)
			if err != nil {
				return result, err
			}
		case PARSE_ANY_INTEGER:
			err, value := parseAnyInteger(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}

			result, err = step.ParseHandler(nil, value, PARSE_ANY_INTEGER, data)
			if err != nil {
				return result, err
			}
		case PARSE_ANY_QUOTED_STRING:
			err, value := parseAnyQuotedString(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}

			result, err := step.ParseHandler(err, value, PARSE_ANY_QUOTED_STRING, data)
			if err != nil {
				return result, err
			}
		case PARSE_COMMA:
			err, value := parseComma(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}

			result, err = step.ParseHandler(err, value, PARSE_COMMA, data)
			if err != nil {
				return result, err
			}
		case PARSE_COLON:
			err, value := parseColon(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		case PARSE_STRING_CHOICE:
			err, value := parseStringChoice(l, step.ParsedValues, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}

			result, err = step.ParseHandler(err, value, PARSE_STRING_CHOICE, data)
			if err != nil {
				return result, err
			}
		case PARSE_STRING_LIST:
			err, _ := parseStringList(l, step.ParsedValues, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}

			result, err = step.ParseHandler(err, nil, PARSE_STRING_LIST, data)
			if err != nil {
				return result, err
			}
		case PARSE_QUESTION:
			err, value := parseQuestion(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		case PARSE_LESS_THAN:
			err, value := parseLessThan(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		case PARSE_GREATER_THAN:
			err, value := parseGreaterThan(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}

		case PARSE_EXCLAMATION:
			err, value := parseExclamation(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		case PARSE_PLUS:
			err, value := parsePlus(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		case PARSE_PERCENT:
			err, value := parsePercent(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		case PARSE_EQUAL:
			err, value := parseEqual(l, step.Options)
			if err != nil {
				ifDebug(debug, fmt.Printf, "Parse: Error parsing step %s in rule %s: %s\n", step.Name, rule.Name, err.Error())
				if step.SkipOnTypeError {
					return PARSE_RESULT_SKIP, err
				} else {
					return PARSE_RESULT_FAILURE, fmt.Errorf("error parsing step %s in rule %s: %w", step.Name, rule.Name, err)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				return result, err
			}
		default:
			return PARSE_RESULT_FAILURE, fmt.Errorf("unknown parser type %d", step.ParserType)
		}
	}
	return PARSE_RESULT_SUCCESS, nil
}

// Parse processes the input string using the provided rules.
// It initializes a Lexer with the input string and iterates through the rules.
// For each rule, it attempts to parse the input and calls the ParseHandler for each step.
func (p *ParserObject) Parse(rules []ParseRule, data interface{}) (int, error) {
	if p.Debug {
		fmt.Printf("Parsing input: %s\n", p.Input)
	}
	for _, rule := range rules {
		l := NewLexer(p.Input, p.Exclude)
		if p.Debug {
			fmt.Printf("Parse: Trying rule: %s\n", rule.Name)
		}
		result, err := parseRule(l, rule, &data, p.Debug)
		switch result {
		case PARSE_RESULT_SUCCESS:
			if p.Debug {
				fmt.Printf("Parse: Rule %s succeeded\n", rule.Name)
			}
			return result, err
		case PARSE_RESULT_FAILURE:
			if p.Debug {
				fmt.Printf("Parse: Rule %s failed with error: %s\n", rule.Name, err.Error())
			}
			return result, err
		case PARSE_RESULT_SKIP:
			if p.Debug {
				fmt.Printf("Parse: Rule %s skipped due to error: %s\n", rule.Name, err.Error())
			}
			continue
		}
	}
	return PARSE_RESULT_FAILURE, fmt.Errorf("no rules matched")
}
