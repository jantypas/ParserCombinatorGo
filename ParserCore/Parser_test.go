package ParserCore

import (
	"fmt"
	"testing"
)

type DataObject struct {
	TestString string
	TestInt    int
	TestFloat  float64
	TestQuote  string
}

func TestParserObject_Parse(t *testing.T) {
	DO := DataObject{}
	Rules := []ParseRule{
		{
			Name: "TestRule",
			Steps: []ParserRuleStep{
				{
					Name:        "ParseString",
					ParserType:  PARSE_ANY_STRING,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						do := (*data).(*DataObject)
						do.TestString = token.(string)
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseInteger",
					ParserType:  PARSE_ANY_INTEGER,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_FAILURE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						fmt.Printf("Parsing INT %v", token)
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						do := (*data).(*DataObject)
						do.TestInt = token.(int)
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseFloat",
					ParserType:  PARSE_ANY_FLOAT,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						do := (*data).(*DataObject)
						do.TestFloat = token.(float64)
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseQuotedString",
					ParserType:  PARSE_ANY_QUOTED_STRING,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						do := (*data).(*DataObject)
						do.TestQuote = token.(string)
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseComma",
					ParserType:  PARSE_COMMA,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Comma is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParsQuestion",
					ParserType:  PARSE_QUESTION,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Colon is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseLessThan",
					ParserType:  PARSE_LESS_THAN,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// LessThan is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseGreaterThan",
					ParserType:  PARSE_GREATER_THAN,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// GreaterThan is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseExclamation",
					ParserType:  PARSE_EXCLAMATION,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Exclamation is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParsePlus",
					ParserType:  PARSE_PLUS,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Plus is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParseEqual",
					ParserType:  PARSE_EQUAL,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Equal is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParsePercent",
					ParserType:  PARSE_PERCENT,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Percent is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:        "ParsColon",
					ParserType:  PARSE_COLON,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						// Colon is just a separator, no action needed
						return PARSE_RESULT_SUCCESS, nil
					},
				},
			},
		},
	}

	p := ParserObject{
		Debug:   true,
		Input:   "String 123 123.12 \"Test 123\",?<>!+=%:",
		Exclude: []string{"TEST"},
	}
	parse, err := p.Parse(Rules, &DO)
	if err != nil {
		t.Errorf("Parse failed with error: %v", err)
		return
	} else {
		if parse != PARSE_RESULT_SUCCESS {
			t.Errorf("Parse did not succeed, got: %d", parse)
			return
		}
		fmt.Printf("====== Parse Result = %v\n", DO)
	}
}

func TestParserObject_OptionalString(t *testing.T) {
	DO := DataObject{}
	Rules := []ParseRule{
		{
			Name: "TestOptionalString",
			Steps: []ParserRuleStep{
				{
					Name:        "ParseString",
					ParserType:  PARSE_ANY_STRING,
					Options:     PARSE_OPTION_CONVERT_TO_UPPERCASE,
					SkipOnError: PARSE_RESULT_SKIP_RULE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						do := (*data).(*DataObject)
						do.TestString = token.(string)
						return PARSE_RESULT_SUCCESS, nil
					},
				},
				{
					Name:         "ParseOptionalString",
					ParserType:   PARSE_STRING_CHOICE,
					ParsedValues: []string{"OPTIONAL", "STRING"},
					Options:      PARSE_OPTION_CONVERT_TO_UPPERCASE | PARSE_OPTION_STRING_IS_OPTIONAL,
					SkipOnError:  PARSE_RESULT_FAILURE,
					ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
						if err != nil {
							return PARSE_RESULT_FAILURE, err
						}
						do := (*data).(*DataObject)
						if token != nil {
							do.TestString += " " + token.(string)
						}
						return PARSE_RESULT_SUCCESS, nil
					},
				},
			},
		},
	}

	p := ParserObject{
		Debug:   true,
		Input:   "Bob Binklestein",
		Exclude: []string{"TEST"},
	}
	parse, err := p.Parse(Rules, &DO)
	if err != nil {
		t.Errorf("Parse failed with error: %v", err)
		return
	} else {
		if parse != PARSE_RESULT_SUCCESS {
			t.Errorf("Parse did not succeed, got: %d", parse)
			return
		}
		fmt.Printf("====== Parse Result = %v\n", DO)
	}
}
