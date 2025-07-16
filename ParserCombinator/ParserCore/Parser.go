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
	ParseHandler func(err error, token interface{}, tokType int, data *interface{}) (int, error)
}

// ParserRule defines a rule that consists of multiple steps.
type ParseRule struct {
	Name  string
	Steps []ParserRuleStep
}

func parseRule(l *Lexer, rule ParseRule, data *interface{}, debug bool) (int, error) {
	var result int
	for _, step := range rule.Steps {
		if debug {
			fmt.Printf("Parse: Trying step: %s for rule: %s\n", step.Name, rule.Name)
		}
		switch step.ParserType {
		case PARSE_ANY_STRING:
			err, value := parseAnyString(l, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed string: %s\n", value)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_ANY_STRING, data)
			if err != nil {
				return result, err
			}
		case PARSE_ANY_FLOAT:
			err, value := parseAnyFloat(l, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed float: %f\n", value)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_ANY_FLOAT, data)
			if err != nil {
				return result, err
			}
		case PARSE_ANY_INTEGER:
			err, value := parseAnyInteger(l, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return result, err
			} else {
				result, err = step.ParseHandler(nil, value, PARSE_ANY_INTEGER, data)
				if err != nil {
					if debug {
						fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
					}
					return result, err
				}
			}
		case PARSE_ANY_QUOTED_STRING:
			err, value := parseAnyQuotedString(l, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed quoted string: %s\n", value)
				}
			}
			result, err := step.ParseHandler(err, value, PARSE_ANY_QUOTED_STRING, data)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return result, err
			}
		case PARSE_COMMA:
			err, value := parseComma(l, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed comma: %s\n", value)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COMMA, data)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return result, err
			}
		case PARSE_COLON:
			err, value := parseColon(l, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed colon: %s\n", value)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_COLON, data)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return result, err
			}
		case PARSE_STRING_CHOICE:
			err, value := parseStringChoice(l, step.ParsedValues, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed string choice: %s\n", value)
				}
			}
			result, err = step.ParseHandler(err, value, PARSE_STRING_CHOICE, data)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return result, err
			}
		case PARSE_STRING_LIST:
			err, _ := parseStringList(l, step.ParsedValues, step.Options)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return PARSE_FAILURE, err
			} else {
				if debug {
					fmt.Printf("Parse: Successfully parsed string list: %v\n", step.ParsedValues)
				}
			}
			result, err = step.ParseHandler(err, nil, PARSE_STRING_LIST, data)
			if err != nil {
				if debug {
					fmt.Printf("Parse: Error in step %s: %v\n", step.Name, err)
				}
				return result, err
			}
		default:
			return PARSE_FAILURE, fmt.Errorf("unknown parser type %d", step.ParserType)
		}
	}
	return PARSE_SUCCESS, nil
}

func Parse(l *Lexer, rules []ParseRule, data interface{}, debug bool) (int, error) {
	if debug {
		fmt.Printf("Parsing input: %s\n", l.input)
	}
	for _, rule := range rules {
		if debug {
			fmt.Printf("Parse: Trying rule: %s\n", rule.Name)
		}
		result, err := parseRule(l, rule, &data, debug)
		if err != nil {
			return result, err
		} else {
			if debug {
				fmt.Printf("Parse: Successfully parsed rule: %s\n", rule.Name)
			}
			return result, nil
		}
	}
	return PARSE_SUCCESS, nil
}
