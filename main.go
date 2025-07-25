package main

// Parser demo program
// This program demonstrates how to use the ParserCombinator library to parse commands
//
import (
	"encoding/json"
	"fmt"
	"github.com/jantypas/ParserCombinatorGo/ParserCore"
	"github.com/jantypas/ParserCombinatorGo/Rulebase"
)

func main() {
	txtinput := "Display portfolio"
	// Create our basic opaque data object
	do := Rulebase.DataObject{}
	// Create our parser object with our parameters
	p := ParserCore.ParserObject{
		Debug:   false,                   // Enable debug output
		Input:   txtinput,                // The text we intend to parse
		Exclude: []string{"?", "PLEASE"}, // Ignore these words
	}
	// Do the actual parse step with our rules
	res, err := p.Parse(Rulebase.RuleSet, &do)

	// If we get an error response, we had a problem parsing the input
	if res != ParserCore.PARSE_RESULT_SUCCESS || err != nil {
		println("Error parsing input:", res, err.Error())
		return
	} else {
		// We parsed it, so decode our data object
		txt, err := json.MarshalIndent(do, "", "    ")
		fmt.Println("Given input:", ParserCore.BlueText, txtinput, ParserCore.ResetText, "we get:")
		if err != nil {
			println("Error encoding data object:", err.Error())
		}
		fmt.Println("Completed:\n", ParserCore.BlueText, string(txt), ParserCore.ResetText)
	}
}
