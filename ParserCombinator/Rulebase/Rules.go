package Rulebase

import (
	"ParserCombinator/ParserCombinator/ParserCore"
	"errors"
)

// For our stock example, we store the decoded data here
// DataObject is a struct that holds the parsed data from the commands
type DataObject struct {
	Command   string // The command to execute, e.g., "MOVE" or "WHAT IS AT"
	NumShares int    // The number of shares to buy or sell
	StockName string // The name of the stock, e.g., "Futzco"
}

// This rule decodes phrases such as
// "BUY 100 SHARES" OF Futzco"
// "SELL 50 SHARES" OF Futzco"
var BuySellStockRule = ParserCore.ParseRule{
	Name: "BuySellStockRule",
	Steps: []ParserCore.ParserRuleStep{
		{
			// Look for BUY or SELL, if we don't find it,skip to the next rule
			Name:            "Command",
			ParserType:      ParserCore.PARSE_STRING_CHOICE,
			ParsedValues:    []string{"BUY", "SELL"},
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			SkipOnTypeError: true,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				do := (*data).(*DataObject)
				do.Command = token.(string)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			// Look for an integer number of shares
			Name:            "NumShares",
			ParserType:      ParserCore.PARSE_ANY_INTEGER,
			SkipOnTypeError: false, // If we fail, this rule fails
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				dp := (*data).(*DataObject)
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				value := token.(int)
				if value <= 0 {
					return ParserCore.PARSE_RESULT_FAILURE,
						errors.New("Share number must be 1 or greater")
				}
				dp.NumShares = value
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			// Look for the words SHARES and OF in order
			Name:            "SharesAndOf",
			ParserType:      ParserCore.PARSE_STRING_LIST,
			ParsedValues:    []string{"SHARES", "OF"},
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			SkipOnTypeError: false,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				// We don't need to do anything with this token, just ensure it exists
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			// Look for the stock name, which is a string
			Name:            "StockName",
			ParserType:      ParserCore.PARSE_ANY_STRING,
			SkipOnTypeError: false, // If we fail, this rule fails
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.StockName = token.(string)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var WhatIsAtRule = ParserCore.ParseRule{
	Name: "ShowPortfolio",
	Steps: []ParserCore.ParserRuleStep{
		{
			// Look for the command "SHOW PORTFOLIO:
			Name:            "Command",
			ParserType:      ParserCore.PARSE_STRING_LIST,
			ParsedValues:    []string{"SHOW", "PORTFOLIO"},
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			SkipOnTypeError: true, // If we don't find this, skip to the next rule
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.Command = "ShowPortfolio"
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var LiquidateRule = ParserCore.ParseRule{
	Name: "Liquidate",
	Steps: []ParserCore.ParserRuleStep{
		{
			// Look for the command "LIQUIDATE"
			Name:            "Command",
			ParserType:      ParserCore.PARSE_STRING_LIST,
			ParsedValues:    []string{"LIQUIDATE"},
			Options:         ParserCore.PARSE_OPTION_CONVER_TO_UPPERCASE,
			SkipOnTypeError: true, // If we don't find this, skip to the next rule
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.Command = "Liquidate"
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			// Look for the symbol name
			Name:            "Symbol",
			ParserType:      ParserCore.PARSE_ANY_STRING,
			SkipOnTypeError: false, // If we fail, this rule fails
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.StockName = token.(string)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var RuleSet = []ParserCore.ParseRule{
	BuySellStockRule,
	WhatIsAtRule,
	LiquidateRule,
}
