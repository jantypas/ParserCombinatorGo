# ParseCominator - A simple parser combinator in Golang to handle 

The ParserCombinator is a simple Golang package to handle very basic, CLI style parsing, or 
parsing that might be required in a textual game.  If you are looking to parse larger,
multiple objects, or a programming language, this is
not appropriate... consider using a full parser library or 
parser generator instead.

Let us consider cases where this library might be appropriate.
# Example : Parsing a stock trading application.

In this application, we might have requests such as:
* BUY 100 shares of AAPL
* SELL 50 shares of 
* DISPLAY STOCK AAPL
* DISPLAY PORTFOLIO
* LIQUIDATE

How might we do this?  
In a full parser, we might have a grammar that looks like this:

* Command : BuyAction | SellAction | DipplayAction | SellAllAction ;
* BuyAction : BUY INTEGER SHARES OF STRING ;

Now let's get into how this all works....

# First step -- the Lexer
The lexer is the portion of the code that breaks text into useful objects.
Input streams aren't nice and simple and well-behaved.  Consider our example.  
We could have the user providing Buy 100 Shares of Futzco.
But this isn't guaranteed.  We can parse this as:
* BUY 199 shares of Futzco
* BU Y 100 Sha res of Fut zco
* B U Y 1 0 0 S h a r e s.... 
 
You get the idea.  The lexer's job is to break the input stream into useful objects.  The 
lexer has rules to tell it what "words" are.  In general, we have the following rules:

* Words are split on whitespace
* We understand several "types" of text.
  * STRINGS - A string is a sequence of characters that are not whitespace.
  * INTEGER - A sequence of digits that can be converted to an integer.
  * FLOAT - A sequence of digits that can be converted to a float.
  * QUOTED_STRING - A sequence of characters that are enclosed in quotes.
  * STRING_CHOICE - A set of possible strings.  We check the input against the list and return OK if the string we found is in this list.
  * STRING_LIST - A list of strings we must match against in succession.  If we expect a list of ALPHA, BETA, GAMMA, the next three tokens must be ALPHA, BETA or GAMMA.
  * COMMA - The lexer must recognize commas
  * COLON - The lexer must recognize colons
  * ERROR - We got a token we don't understand.
  * EOF - The end of the input stream.

When we receive a token, we receive two items:
* The token type, which is one of the types above.
* The token value.  The value is the actual contents of the type.  For an integer for example, it's the actual number.

# What the lexer needs:

When a Lexer is created via NewLexer, the Lexer object can take a few fields.
* Input - The input stream to parse.  This is the string the user gave us
* Exclude - A list of strings to exclude from the input stream.  This is useful for ignoring comments, or other text that we don't want give to the parser.

The _Exclude_ string list is a list of words we should remove from the lexing process such that, when the lexer attempts to run, these words are already removed from the input.
This is useful for ignoring comments, or other text that we don't want give to the parser.

So, given a string such as "Please buy 100 shares of Futzco?", with the Lexer above, and an exclude list
of "PLEASE" and "?", we will receive the tokens: "BUY", "100", "SHARES", "OF", "FUTZCO".

# The Parser - Taking the Lexical Tokens and Parsing Them into meaning

The Lexer helps us out by turning strings of characters into "words" but it has no idea what 
words make a valid sentence.  The Lexer could give us "SHARES FUTZCO 100 BUY".  This would be valid at
the lexical level, but it makes no sense.  We need additional rules to tell us what a valid sentence is.
This is the job of the Parser.  It has a set of rules that determine what lexical tokens must follow, in what order, 
to make valid sentences, both at a syntactic and semantic level.

We achieve this by giving the parser a set of parse rules, which it runs one after another, until it 
finds a rule that can be satisfied.  Much like firewall rules, it just walks down the rule list,
checking the lexical tokens against each rule until it either finds a successful match, finds a hard failure where
we know we had a match but it was incorrect at some level, or it runs out of rules to check.  At that point,
it simply informs us that no matching rule could be foumd.  We define these rules in code -- in this case, 
a list of Golang structures.

```aiignore
package Rulebase

import (
	"ParserCombinator/ParserCombinator/ParserCore"
	"errors"
)
```
That import brings in our parser code.  In other languages it might an #include directive.

```
// For our stock example, we store the decoded data here
// DataObject is a struct that holds the parsed data from the commands
type DataObject struct {
	Command   string `json:"command"`   // The command to execute, e.g., "MOVE" or "WHAT IS AT"
	NumShares int    `json:"numShares"` // The number of shares to buy or sell
	StockName string `json:"stockName"` // The name of the stock, e.g., "Futzco"
}
```
Here, we define our opaque data strcuture.  The parser cranks through the rules, stuffing values it finds into this structure.
The parser doesn't "know" anything about the data structure, it just knows that it will be passed a pointer to this structure, and,
for values it finds, it will stuff them into the structure. 
```

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
```
This is our first rule.  It parses our BUY and SELL directives.  Notice we have a rule, with a name,
and an array of steps we will walk down, one after another.  Each step parses a piece of input.  When the step
completes, it calls a user-defined function to alter our data object.  For each step, we return a status value and, 
possibly an error.  If the step completes, we return the status of PARSE_RESULT_SUCCESS, which tells the parser that this
step completed successfully.  A value of PARSE_RESULT_FAILURE indicates the step failed and the error
value returned contains an error message that should be returned to the caller.
We can also return PARSE_RESULT_SKIP_STEP, which tells the parser that this step did fail, 
but we don't need to do anything with it.  This is useful for steps that are just there to ensure 
a token exists, but we don't need to do anything with it, or, we can use it to say "OK, this rule failed, 
but try the next rule rather than failing the entire rule chain".

Let's look at this rule in more detail:

Our first step indicates we have a choice of two possible strings, BUY or SELL.  If we don't get either 
a BUY or SELL, this rule fails.  Before we parse the strings, we convert our tokens to upper case.
We also indicate, should the rule fail, don't fail the entire parsing process, skip to the next rule and try again.

This is syntactic phase.  We define what we expect (a choice of two strings).  If we don't get one of them,
we fail, and in this case, we skip to the next rule, but we still haven't handled the semantic case.  This is the
job of the user defined function.  Here, there's nothing much we have to do except, assuming we passed the syntactic case, 
take the string BUY or SELL (we converted it to upper case remember?) and put it in our data structure.
We then return PARSE_RESULT_SUCCECSS - meaning, we can go on to the next step.

In the next step of this rule, we look for an integer value.  A failure here would 
fail the entire rule.  THere's no need to do case conversion etc so it's ignored.  
Assuming we receive an integer, we call our semantic function as before.  In this case, 
we do an additioanl check before storing the value -- we cannot have a negative number of shares nor can 
have zero shares.  If we do, we return a PARSE_RESULT_FAILURE and an error message.

In our next step, we look for the words SHARES and OF in that order.  
If we don't find them, we fail the rule.  This rule is useful for cases where you just need to 
find these words -- but you really don't need to do anything with them.  Thus, in our 
semantic function, we just return PARSE_RESULT_SUCCESS, indicating we found the words.

Finally, we look for the stock name, which is a string.  Again, if we don't find it, we fail the rule.
In this case, we'll accecpt any string, and convert it to upper case.  
We then store the stock name in our data structure via our semantic function.

```
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
```
Here, we define a rule to display a stock.  The rule is called DisplayStockRule.
It has a name, and a series of steps.  The first step looks for the command DISPLAY STOCK.  
If we don't find it, we skip to the next rule.
If we do find it, we set the command in our data structure to "DisplayStock".  
The second step looks for the stock name, which is a string.  
If we don't find it, we fail the rule.  If we do, we store the stock name in our data structure.
```
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
```
This rule is similar to the DisplayStockRule, but it looks for the command DISPLAY PORTFOLIO.
In this case, we're simply looking for the two words DISPLAY and PORTFOLIO in that order.  
If we find them, we set the command in our data structure to "DISPLAY-PORTFOLIO".
```

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
```
This rule is our stock liquidate everything rule -- it simply looks for the word LIQUIDATE.
```
var RuleSet = []ParserCore.ParseRule{
	BuySellStockRule,
	DisplayStockRule,
	DisplayPortfolioRule,
	LoquiddateRule,
}
```
And this is the end our rules, so now we just put them all in a list or _ParseRule_ objects.
Once we have our rules, we can use them with code such as this:

```aiignore
package main

// Parser demo program
// This program demonstrates how to use the ParserCombinator library to parse commands
//
import (
	"ParserCombinator/ParserCombinator/ParserCore"
	"ParserCombinator/ParserCombinator/Rulebase"
	"encoding/json"
	"fmt"
)
```
Here we bring in our parser code and our rules.  The rules are defined in the Rulebase package.

```
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
```
Create our parser object...
```
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
```
Parse out input and produce some nice JSON style output of what we received.

# Benefits of Lexers and Parsers
A qiock aside as tp why lexers and parsers are useful.  We've built a very simple 
lexer and parser for stock trading.  Our application doesn't see the language used.  
English goes in, and data structures come out.  Our application doesn't think in English, it 
only sees the data structures.  We could replace the parse and use one for another language and the 
application would never know the difference.  Want the application to understand Spanish, just change the parser. 
Want Germain, change hte parser --- the application will never know.  Want Klingon, change the parser.

One can even do the reverse with what is called a _renderer_. If a parser takes input and produces data 
structures, a renderer, takes data structures and produces text as output. Much like a parser, changing 
language support simply involves updating the renderer.
