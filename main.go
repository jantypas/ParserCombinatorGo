package main

// Parser demo program
// This program demonstrates how to use the ParserCombinator library to parse commands
//
import (
	"ParserCombinator/ParserCombinator/ParserCore"
)

// DataObject is a structure that holds the parsed data we care about
type DataObject struct {
	Command   string
	Distance  int
	Direction string
	XPos      int
	YPos      int
}

// This is a list of rules that we can use to parse commands.
var parserRules = []ParserCore.ParseRule{
	{
		//
		// The MoveRunRule expects input of the form
		// MOVE or RUN followed by an integer distrance and a compass direction such as
		// MOVE 12 SOUTH or RUN 5 EAST
		//
		Name: "MoveOrRunRule",
		Steps: []ParserCore.ParserRuleStep{
			{
				Name:         "MoveOrRun_Step1",
				ParserType:   ParserCore.PARSE_STRING_CHOICE,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"MOVE", "RUN"}, // Accept either one
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					if err != nil {
						// If we don't get a MOVE Or RUN, skip to the next possible rule
						return ParserCore.SKIP_RULE_ON_ERROR, err
					}
					do := (*data).(*DataObject)
					do.Command = token.(string)
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				// Get the integer distance
				Name:       "MoveOrRun_Step2",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					do := (*data).(*DataObject)
					do.Distance = token.(int)
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				// Get the direction, which is a string from a list of choices
				Name:         "MoveOrRun_Step3",
				ParserType:   ParserCore.PARSE_STRING_CHOICE,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"NORTH", "SOUTH", "EAST", "WEST"},
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					do := (*data).(*DataObject)
					do.Direction = token.(string)
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
		},
	},
	{
		//
		// The WhatIsATRule expects a phrase of the form:
		// WHAT IS AT <X>,<Y> where <X> and <Y> are integers ex:
		// WHAT IS AT 92, 12
		//
		Name: "WhatIsAtRule",
		Steps: []ParserCore.ParserRuleStep{
			{
				Name:         "WhatIsAt_Step1",
				ParserType:   ParserCore.PARSE_STRING_LIST,
				Options:      ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
				ParsedValues: []string{"WHAT", "IS", "AT"},
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					(*data).(*DataObject).Command = "WHATIS"
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				// Get the X integer position
				Name:       "WhatIsAt_Step2",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					(*data).(*DataObject).XPos = token.(int)
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				// Get the comma separator
				Name:       "WhatIsAt_Step3",
				ParserType: ParserCore.PARSE_COMMA,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
			{
				// Get the Y integer position
				Name:       "WhatIsAt_Step4",
				ParserType: ParserCore.PARSE_ANY_INTEGER,
				ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
					(*data).(*DataObject).YPos = token.(int)
					return ParserCore.PARSE_SUCCESS, nil
				},
			},
		},
	},
}

func main() {
	do := DataObject{} // Create our special data object
	_, err := ParserCore.Parse("What is at 153,54", parserRules, &do, false)
	if err != nil {
		println("Error parsing input:", err.Error())
		return
	} else {
		println("Completed",
			"[", do.Command, ":", do.Distance, ":", do.Direction,
			" X:", do.XPos, " Y:", do.YPos, "]")
	}
}
