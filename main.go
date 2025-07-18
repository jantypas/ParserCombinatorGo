package main

// Parser demo program
// This program demonstrates how to use the ParserCombinator library to parse commands
//
import (
	"ParserCombinator/ParserCombinator/ParserCore"
	"ParserCombinator/ParserCombinator/Rulebase"
)

func main() {
	// Create our basic opaque data object
	do := Rulebase.DataObject{}
	// Create our parser object with our parameters
	p := ParserCore.ParserObject{
		Debug:   true,                             // Enable debug output
		Input:   "Please Buy 12 shares of Futzco", // The text we intend to parse
		Exclude: []string{"?", "PLEASE"},          // Ignore these words
	}
	// Do the actual parse step with our rules
	res, err := p.Parse(Rulebase.RuleSet, &do)

	// If we get an error response, we had a problem parsing the input
	if res != ParserCore.PARSE_RESULT_SUCCESS || err != nil {
		println("Error parsing input:", res, err.Error())
		return
	} else {
		// We parsed it, so decode our data object
		println("Completed",
			"[", do.Command, ":", do.NumShares, ":", do.StockName, "]")
	}
}
