package ParserCore

import (
	"fmt"
)

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
)

// When we parse something, here are possible error codes.
const (
	PARSE_SUCCESS = iota
	PARSE_FAILURE
	SKIP_RULE_ON_ERROR
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
	Name         string
	ParserType   int
	Options      int
	ParsedValues []string
	ParseHandler func(err error, token interface{}, tokType int, data interface{}) (int, error, interface{})
}

// ParserRule defines a rule that consists of multiple steps.
type ParseRule struct {
	Name  string
	Steps []ParserRuleStep
}

func parseRule(l Lexer, rule ParseRule, data interface{}) (int, error, interface{}) {
	for _, step := range rule.Steps {
		switch step.ParserType {
		case PARSE_ANY_STRING:
			err, value := parseAnyString(l, step.Options)
			return step.ParseHandler(err, value, PARSE_ANY_STRING, data)
		case PARSE_ANY_FLOAT:
			err, value := parseAnyFloat(l, step.Options)
			return step.ParseHandler(err, value, PARSE_ANY_FLOAT, data)
		case PARSE_ANY_INTEGER:
			err, value := parseAnyInteger(l, step.Options)
			return step.ParseHandler(err, value, PARSE_ANY_INTEGER, data)
		case PARSE_ANY_QUOTED_STRING:
			err, value := parseAnyQuotedString(l, step.Options)
			return step.ParseHandler(err, value, PARSE_ANY_QUOTED_STRING, data)
		case PARSE_COMMA:
			err, value := parseComma(l, step.Options)
			return step.ParseHandler(err, value, PARSE_COMMA, data)
		case PARSE_COLON:
			err, value := parseColon(l, step.Options)
			return step.ParseHandler(err, value, PARSE_COLON, data)
		case PARSE_STRING_CHOICE:
			err, value := parseStringChoice(l, step.ParsedValues, step.Options)
			return step.ParseHandler(err, value, PARSE_STRING_CHOICE, data)
		case PARSE_STRING_LIST:
			err, _ := parseStringList(l, step.ParsedValues, step.Options)
			return step.ParseHandler(err, nil, PARSE_STRING_LIST, data)
		default:
			return PARSE_FAILURE, fmt.Errorf("unknown parser type %d", step.ParserType), nil
		}
	}
	return PARSE_FAILURE, fmt.Errorf("no steps defined for rule %s", rule.Name), nil
}

func Parse(l Lexer, rules []ParseRule, data interface{}) (int, error, interface{}) {
	for _, rule := range rules {
		result, err, parsedData := parseRule(l, rule, data)
		if err != nil {
			return result, err, parsedData
		}
		if result == PARSE_SUCCESS {
			return result, nil, parsedData
		}
	}
	return PARSE_FAILURE, fmt.Errorf("no matching rule found"), nil
}
