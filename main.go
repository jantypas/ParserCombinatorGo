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
		Debug:   true,                       // Enable debug output
		Input:   "Please What is at 153,54", // The text we intend to parse
		Exclude: []string{"?", "PLEASE"},    // Ignore these words
	}
	// Do the actual parse step with our rules
	_, err := p.Parse(Rulebase.Rules, &do)

	// If we get an error response, we had a problem parsing the input
	if err != nil {
		println("Error parsing input:", err.Error())
		return
	} else {
		// We parsed it, so decode our data object
		println("Completed",
			"[", do.Command, ":", do.Distance, ":", do.Direction, " X:", do.X, " Y:", do.Y, "]")
	}
}
