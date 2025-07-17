package Rulebase

import "ParserCombinator/ParserCombinator/ParserCore"

type DataObject struct {
	Command   string // The command to execute, e.g., "MOVE" or "RUN"
	Distance  int    // The distance to move or run
	Direction string // The direction to move or run, e.g., "NORTH", "SOUTH", "EAST", "WEST"
	X         int    // X coordinate for "WHAT IS AT"
	Y         int    // Y coordinate for "WHAT IS AT"
}

var moveRunRule = ParserCore.ParseRule{
	Name: "MoveOrRunRule",
	Steps: []ParserCore.ParserRuleStep{
		{
			Name:            "MoveOrRun_Step1",
			ParserType:      ParserCore.PARSE_STRING_CHOICE,
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			ParsedValues:    []string{"MOVE", "RUN"}, // Accept either one
			SkipOnTypeError: true,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					// If we don't get a MOVE Or RUN, skip to the next possible rule
					return ParserCore.PARSE_RESULT_SKIP, err
				}
				do := (*data).(*DataObject)
				do.Command = token.(string)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			Name:            "MoveOrRun_Step2",
			ParserType:      ParserCore.PARSE_ANY_INTEGER,
			SkipOnTypeError: false,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				do := (*data).(*DataObject)
				do.Distance = token.(int)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			Name:            "MoveOrRun_Step3",
			ParserType:      ParserCore.PARSE_STRING_CHOICE,
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			SkipOnTypeError: false,
			ParsedValues:    []string{"NORTH", "SOUTH", "EAST", "WEST"},
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				do := (*data).(*DataObject)
				do.Direction = token.(string)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var WhatIsRule = ParserCore.ParseRule{
	Name: "WhatIsRule",
	Steps: []ParserCore.ParserRuleStep{
		{
			Name:            "WhatIs_Step1",
			ParserType:      ParserCore.PARSE_STRING_LIST,
			ParsedValues:    []string{"WHAT", "IS", "AT"},
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			SkipOnTypeError: true, // If we don't get the expected string, skip to the next rule
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					// If we don't get the expected string, skip to the next rule
					return ParserCore.PARSE_RESULT_SKIP, err
				}
				do := (*data).(*DataObject)
				do.Command = "WHAT IS AT"
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			Name:            "WhatIs_Step2",
			ParserType:      ParserCore.PARSE_ANY_INTEGER,
			SkipOnTypeError: false,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				do := (*data).(*DataObject)
				do.X = token.(int)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			Name:            "WhatIs_Step3",
			ParserType:      ParserCore.PARSE_COMMA,
			SkipOnTypeError: false,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					// If we don't get a comma, we can't continue
					return ParserCore.PARSE_RESULT_SKIP, err
				}
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			Name:            "WhatIs_Step4",
			ParserType:      ParserCore.PARSE_ANY_INTEGER,
			SkipOnTypeError: false,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					// If we don't get an integer, we can't continue
					return ParserCore.PARSE_RESULT_SKIP, err
				}
				do := (*data).(*DataObject)
				do.Y = token.(int)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var Rules = []ParserCore.ParseRule{
	moveRunRule,
	WhatIsRule,
}
