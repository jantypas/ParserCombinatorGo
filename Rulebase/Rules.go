package Rulebase

import (
	"errors"
	"github.com/jantypas/ParserCombinatorGo/ParserCore"
)

// For our stock example, we store the decoded data here
// DataObject is a struct that holds the parsed data from the commands
type DataObject struct {
	Command   string `json:"command"`   // The command to execute, e.g., "MOVE" or "WHAT IS AT"
	NumShares int    `json:"numShares"` // The number of shares to buy or sell
	StockName string `json:"stockName"` // The name of the stock, e.g., "Futzco"
}

// This rule decodes phrases such as
// "BUY 100 SHARES" OF Futzco"
// "SELL 50 SHARES" OF Futzco"
var BuySellStockRule = ParserCore.ParseRule{
	Name: "BuySellStockRule",
	Steps: []ParserCore.ParserRuleStep{
		{
			// Look for BUY or SELL, if we don't find it,skip to the next rule
			Name:         "Command",
			ParserType:   ParserCore.PARSE_STRING_CHOICE,
			ParsedValues: []string{"BUY", "SELL"},
			Options:      ParserCore.PARSE_OPTION_CONVERT_TO_UPPERCASE,
			SkipOnError:  ParserCore.PARSE_RESULT_SKIP_RULE, // If we don't find this, skip to the next rule
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				do := (*data).(*DataObject)
				do.Command = token.(string)
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			// Look for an integer number of shares
			Name:        "NumShares",
			ParserType:  ParserCore.PARSE_ANY_INTEGER,
			SkipOnError: ParserCore.PARSE_RESULT_FAILURE, // If we fail, this rule fails
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
			Name:         "SharesAndOf",
			ParserType:   ParserCore.PARSE_STRING_LIST,
			ParsedValues: []string{"SHARES", "OF"},
			Options:      ParserCore.PARSE_OPTION_CONVERT_TO_UPPERCASE,
			SkipOnError:  ParserCore.PARSE_RESULT_FAILURE,
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
			Name:        "StockName",
			ParserType:  ParserCore.PARSE_ANY_STRING,
			SkipOnError: ParserCore.PARSE_RESULT_FAILURE, // If we fail, this rule fails
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

// Display stock rule
// Takes the form of DISPLAY STOCK name
var DisplayStockRule = ParserCore.ParseRule{
	Name: "DisplayStockRule",
	Steps: []ParserCore.ParserRuleStep{
		{
			// Look for the command "SHOW PORTFOLIO:
			Name:         "Command",
			ParserType:   ParserCore.PARSE_STRING_LIST,
			ParsedValues: []string{"DISPLAY", "STOCK"},
			Options:      ParserCore.PARSE_OPTION_CONVERT_TO_UPPERCASE,
			SkipOnError:  ParserCore.PARSE_RESULT_SKIP_RULE,
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.Command = "DisplayStock"
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
		{
			// Look for the stock name, which is a string
			Name:        "StockName",
			ParserType:  ParserCore.PARSE_ANY_STRING,
			SkipOnError: ParserCore.PARSE_RESULT_FAILURE, // If we fail, this rule fails
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

// Display Portfolio rule
// Takes the form DISPLAY PORTFOLIOO
var DisplayPortfolioRule = ParserCore.ParseRule{
	Name: "Liquidate",
	Steps: []ParserCore.ParserRuleStep{
		{
			// Look for the command "LIQUIDATE"
			Name:         "Command",
			ParserType:   ParserCore.PARSE_STRING_LIST,
			ParsedValues: []string{"DISPLAY", "PORTFOLIO"},
			Options:      ParserCore.PARSE_OPTION_CONVERT_TO_UPPERCASE,
			SkipOnError:  ParserCore.PARSE_RESULT_SKIP_RULE, // If we don't find this, skip to the next rule
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.Command = "DISPLAY-PORTFOLIO"
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var LoquiddateRule = ParserCore.ParseRule{
	Name: "Liquidate",
	Steps: []ParserCore.ParserRuleStep{
		{
			Name:         "Lquidate",
			ParserType:   ParserCore.PARSE_STRING_LIST,
			ParsedValues: []string{"LIQUIDATE"},
			Options:      ParserCore.PARSE_OPTION_CONVERT_TO_UPPERCASE,
			SkipOnError:  ParserCore.PARSE_RESULT_SKIP_RULE, // If we don't find this, skip to the next rule
			ParseHandler: func(err error, token interface{}, tokType int, data *interface{}) (int, error) {
				if err != nil {
					return ParserCore.PARSE_RESULT_FAILURE, err
				}
				do := (*data).(*DataObject)
				do.Command = "LIQUIDATE"
				return ParserCore.PARSE_RESULT_SUCCESS, nil
			},
		},
	},
}

var RuleSet = []ParserCore.ParseRule{
	BuySellStockRule,
	DisplayStockRule,
	DisplayPortfolioRule,
	LoquiddateRule,
}
