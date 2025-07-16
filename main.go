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
	_, err := ParserCore.Parse("What is at 153,54", Rulebase.Rules, &do, true)
	if err != nil {
		println("Error parsing input:", err.Error())
		return
	} else {
		println("Completed",
			"[", do.Command, ":", do.Distance, ":", do.Direction, " X:", do.X, " Y:", do.Y, "]")
	}
}
