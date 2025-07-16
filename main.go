package main

import (
	"ParserCombinator/ParserCombinator/ParserCore"
	"fmt"
)

type DataObject struct {
	Command   string
	Distance  int
	Direction string
}

var parserRules = []ParserCore.ParseRule{
	{
		Name: "MoveOrRunRule",
		Steps: []ParserCore.ParserRuleStep{
			{
				Name:         "MoveOrRun_Step1",
				ParserType:   ParserCore.PARSE_STRING_CHOICE,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"MOVE", "RUN"},
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					if err != nil {
						return ParserCore.SKIP_RULE_ON_ERROR, err
					}
					if data == nil {
						return ParserCore.PARSE_FAILURE, fmt.Errorf("data is nil")
					}
					do, ok := (*data).(*DataObject)
					if !ok {
						return ParserCore.PARSE_FAILURE, fmt.Errorf("invalid data type: expected *DataObject")
					}
					do.Command = token.(string)
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				Name:       "MoveOrRun_Step2",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					if err != nil {
						return ParserCore.SKIP_RULE_ON_ERROR, err
					}
					do, err := handleData(data)
					if err != nil {
						return ParserCore.PARSE_FAILURE, err
					}
					if value, ok := token.(int); ok {
						do.Distance = value
						return ParserCore.PARSE_SUCCESS, nil
					}
					return ParserCore.PARSE_FAILURE, fmt.Errorf("expected integer for distance")
				},
			},
			{
				Name:         "MoveOrRun_Step3",
				ParserType:   ParserCore.PARSE_STRING_CHOICE,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"NORTH", "SOUTH", "EAST", "WEST"},
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					if err != nil {
						return ParserCore.SKIP_RULE_ON_ERROR, err
					}
					do, err := handleData(data)
					if err != nil {
						return ParserCore.PARSE_FAILURE, err
					}
					if value, ok := token.(string); ok {
						do.Direction = value
						return ParserCore.PARSE_SUCCESS, nil
					}
					return ParserCore.PARSE_FAILURE, fmt.Errorf("expected string for direction")
				},
			},
		},
	},
	{
		Name: "WhatIsAtRule",
		Steps: []ParserCore.ParserRuleStep{
			{
				Name:         "WhatIsAt_Step1",
				ParserType:   ParserCore.PARSE_STRING_LIST,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"WHAT", "IS", "AT"},
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				Name:       "WhatIsAt_Step2",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					if err != nil {
						return ParserCore.SKIP_RULE_ON_ERROR, err
					}
					do, err := handleData(data)
					if err != nil {
						return ParserCore.PARSE_FAILURE, err
					}
					if value, ok := token.(int); ok {
						do.Distance = value
						return ParserCore.PARSE_SUCCESS, nil
					}
					return ParserCore.PARSE_FAILURE, fmt.Errorf("expected integer for distance")
				},
			},
			{
				Name:       "WhatIsAt_Step3",
				ParserType: ParserCore.PARSE_COMMA,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				Name:       "WhatIsAt_Step4",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					if err != nil {
						return ParserCore.SKIP_RULE_ON_ERROR, err
					}
					do, err := handleData(data)
					if err != nil {
						return ParserCore.PARSE_FAILURE, err
					}
					if value, ok := token.(int); ok {
						do.Distance = value
						return ParserCore.PARSE_SUCCESS, nil
					}
					return ParserCore.PARSE_FAILURE, fmt.Errorf("expected integer for distance")
				},
			},
		},
	},
}

func main() {
	lex := ParserCore.NewLexer("MOVE 10 NORTH")
	_, err := ParserCore.Parse(lex, parserRules, &DataObject{}, true)
	if err != nil {
		println("Error parsing input:", err.Error())
		return
	} else {
		println("Completed")
	}
}
func handleData(data *interface{}) (*DataObject, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	do, ok := (*data).(*DataObject)
	if !ok {
		return nil, fmt.Errorf("invalid data type: expected *DataObject")
	}
	return do, nil
}
