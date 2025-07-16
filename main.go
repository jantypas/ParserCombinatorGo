package main

// Parser demo program
// This program demonstrates how to use the ParserCombinator library to parse commands
//
import (
	"ParserCombinator/ParserCombinator/ParserCore"
	"ParserCombinator/ParserCombinator/Rulebase"
)

func main() {
	do := Rulebase.DataObject{} // Create our special data object
	p := ParserCore.ParserObject{
		Debug:   true, // Enable debug output
		Input:   "What is at 153,54",
		Exclude: []string{"?", "PLEASE"},
	}
	_, err := p.Parse(Rulebase.Rules, &do)

	if err != nil {
		println("Error parsing input:", err.Error())
		return
	} else {
		println("Completed",
			"[", do.Command, ":", do.Distance, ":", do.Direction, " X:", do.X, " Y:", do.Y, "]")
	}
}
