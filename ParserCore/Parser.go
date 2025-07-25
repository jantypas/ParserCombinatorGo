package ParserCore

import (
	"fmt"
)

var ParserVersion = "1.0.0"

// ParserObject is the main structure for parsing.
// It contains the input string, a debug flag, and a list of tokens to exclude from parsing.
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
	PARSE_RESULT_SKIP_RULE
	PARSE_RESULT_SKIP_STEP // Skip this step, but continue with the next rule
)

var ResultNames = []string{
	"PARSE_RESULT_SUCCESS",
	"PARSE_RESULT_FAILURE",
	"PARSE_RESULT_SKIP_RULE",
	"PARSE_RESULT_SKIP_STEP",
}

// Options allow us to modify text before parsing
const (
	PARSE_OPTION_CONVERT_TO_UPPERCASE = 1 << iota
	PARSE_OPTION_CONVERT_TO_LOWERCASE
	PARSE_OPTION_STRING_IS_OPTIONAL // If this is set, the string is optional
)

// For each of our rules, there are various steps.
// Each step defines a name, the type of object we expect
// any objects we need to use and functions to handle success and failure.
type ParserRuleStep struct {
	Name         string
	ParserType   int
	Options      int
	SkipOnError  int
	ParsedValues []string
	ParseHandler func(err error, token interface{}, tokType int, data *interface{}) (int, error)
}

// ParserRule defines a rule that consists of multiple steps.
type ParseRule struct {
	Name  string
	Steps []ParserRuleStep
}

// parseRule processes a single rule with its steps.
// It uses the Lexer to read tokens and applies the ParseHandler for each step.
// If any step fails, it returns an error including a request to skip to the next rule
func parseRule(l *Lexer, rule ParseRule, data *interface{}, debug bool) (int, error) {
	// For each step in the rule
	for _, step := range rule.Steps {
		IfDebug(debug, fmt.Printf, "%sParse: Trying step: %s for rule: %s%s\n",
			BlueText, step.Name, rule.Name, ResetText)
		IfDebug(debug, fmt.Printf, "       Expecting token type %s with options %d\n", ParserNames[step.ParserType], step.Options)
		// Call the lexer to get the next token -- requesting a specific type to be decoded
		// If the type is wrong return an error
		// If the type is correct, call the ParseHandler to process the token
		var err error
		var value interface{}
		switch step.ParserType {
		case PARSE_ANY_STRING:
			IfDebug(debug, fmt.Printf, "       Parsing ANY STRING\n")
			err, value = parseAnyString(l, step.Options)
		case PARSE_ANY_FLOAT:
			IfDebug(debug, fmt.Printf, "       Parsing ANY FLOAT\n")
			err, value = parseAnyFloat(l, step.Options)
		case PARSE_ANY_INTEGER:
			IfDebug(debug, fmt.Printf, "       Parsing ANY INTEGER\n")
			err, value = parseAnyInteger(l, step.Options)
		case PARSE_ANY_QUOTED_STRING:
			IfDebug(debug, fmt.Printf, "       Parsing ANY QUOTED STRING\n")
			err, value = parseAnyQuotedString(l, step.Options)
		case PARSE_COMMA:
			IfDebug(debug, fmt.Printf, "       Parsing COMMA\n")
			err, value = parseComma(l, step.Options)
		case PARSE_COLON:
			IfDebug(debug, fmt.Printf, "       Parsing COLON\n")
			err, value = parseColon(l, step.Options)
		case PARSE_STRING_CHOICE:
			IfDebug(debug, fmt.Printf, "       Parsing STRING CHOICE\n")
			err, value = parseStringChoice(l, step.ParsedValues, step.Options)
		case PARSE_STRING_LIST:
			IfDebug(debug, fmt.Printf, "       Parsing STRING LIST\n")
			err, value = parseStringList(l, step.ParsedValues, step.Options)
		case PARSE_QUESTION:
			IfDebug(debug, fmt.Printf, "       Parsing QUESTION\n")
			err, value = parseQuestion(l, step.Options)
		case PARSE_LESS_THAN:
			IfDebug(debug, fmt.Printf, "       Parsing LESS THAN\n")
			err, value = parseLessThan(l, step.Options)
		case PARSE_GREATER_THAN:
			IfDebug(debug, fmt.Printf, "       Parsing GREATER THAN\n")
			err, value = parseGreaterThan(l, step.Options)
		case PARSE_EXCLAMATION:
			IfDebug(debug, fmt.Printf, "       Parsing EXCLAMATION\n")
			err, value = parseExclamation(l, step.Options)
		case PARSE_PLUS:
			IfDebug(debug, fmt.Printf, "       Parsing PLUS\n")
			err, value = parsePlus(l, step.Options)
		case PARSE_PERCENT:
			IfDebug(debug, fmt.Printf, "       Parsing PERCENT\n")
			err, value = parsePercent(l, step.Options)
		case PARSE_EQUAL:
			IfDebug(debug, fmt.Printf, "       Parsing EQUAL\n")
			err, value = parseEqual(l, step.Options)
		default:
			IfDebug(debug, fmt.Printf, "       Unknown parser type %d\n", step.ParserType)
			return PARSE_RESULT_FAILURE, fmt.Errorf("unknown parser type %d", step.ParserType)
		}
		if err != nil {
			IfDebug(debug, fmt.Printf, "%s       Error parsing step %s: %v%s\n",
				RedText, step.Name, err, ResetText)
			return step.SkipOnError, err
		} else {
			IfDebug(debug, fmt.Printf, "%s       Successfully parsed step %s: %v%s\n",
				GreenText, step.Name, value, ResetText)
			result, err := step.ParseHandler(err, value, step.ParserType, data)
			IfDebug(debug, fmt.Printf, "%s       ParseHandler returned result %s: Error = %v%s\n",
				GreenText, ResultNames[result], err, ResetText)
		}
	}
	return PARSE_RESULT_SUCCESS, nil
}

// Parse processes the input string using the provided rules.
// It initializes a Lexer with the input string and iterates through the rules.
// For each rule, it attempts to parse the input and calls the ParseHandler for each step.
func (p *ParserObject) Parse(rules []ParseRule, data interface{}) (int, error) {
	IfDebug(p.Debug, fmt.Printf, "%sParser: Parsing input string: %s%s\n",
		BlueText, p.Input, ResetText)
	for _, rule := range rules {
		l := NewLexer(p.Input, p.Exclude)
		result, err := parseRule(l, rule, &data, p.Debug)
		switch result {
		case PARSE_RESULT_SUCCESS:
			return result, nil
		case PARSE_RESULT_FAILURE:
			return result, err
		case PARSE_RESULT_SKIP_RULE:

			continue
		}
	}
	return PARSE_RESULT_FAILURE, fmt.Errorf("no rules matched")
}
