package main

import "ParserCombinator/ParserCombinator/ParserCore"

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
				ParseHandler: func(err error,
					token interface{}, tokenType int,
					data interface{}) (int, error, interface{}) {
					if err != nil {
						return ParserCore.SKIP_RULE_ON_ERROR, err, nil
					} else {
						do := data.(*DataObject)
						do.Command = token.(string)
						return ParserCore.PARSE_SUCCESS, nil, do
					}
				},
			},
			{
				Name:       "MoveOrRun_Step2",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data interface{}) (int, error, interface{}) {
					if err != nil {
						return ParserCore.PARSE_FAILURE, err, "We were expecting an integer distance"
					} else {
						do := data.(*DataObject)
						do.Distance = token.(int)
						return ParserCore.PARSE_SUCCESS, nil, do
					}
				},
			},
			{
				Name:         "MoveOrRun_Step3",
				ParserType:   ParserCore.PARSE_STRING_CHOICE,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"NORTH", "SOUTH", "EAST", "WEST"},
				ParseHandler: func(err error,
					token interface{}, tokenType int,
					data interface{}) (int, error, interface{}) {
					if err != nil {
						return ParserCore.PARSE_FAILURE, err, "We expect NORTH, SOUTH, EAST or WEST"
					} else {
						do := data.(*DataObject)
						do.Direction = token.(string)
						return ParserCore.PARSE_SUCCESS, nil, do
					}
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
				ParseHandler: func(err error, token interface{}, tokType int, data interface{}) (int, error, interface{}) {
					if err != nil {
						return ParserCore.PARSE_FAILURE, err, "We expect WHAT IS AT"
					} else {
						do := data.(*DataObject)
						do.Command = "WHAT IS AT"
						return ParserCore.PARSE_SUCCESS, nil, do
					}
				},
			},
		},
	},
}

func main() {
	lex := ParserCore.NewLexer("MOVE 10 NORTH")
	_, err, i := ParserCore.Parse(*lex, parserRules, &DataObject{})
	if err != nil {
		println("Error parsing input:", err.Error())
		return
	} else {
		do := i.(*DataObject)
		println("Command:", do.Command)
		println("Distance:", do.Distance)
		println("Direction:", do.Direction)
	}
}
