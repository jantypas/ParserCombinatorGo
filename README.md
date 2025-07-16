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
* SELL 50 shares of TSLA
* DISPLAY PORTFOLIO
* SELL ALL

How might we do this?  
In a full parser, we might have a grammar that looks like this:

* Command : BuyAction | SellAction | DipplayAction | SellAllAction ;
* BuyAction : BUY INTEGER SHARES OF STRING ;

In our package, we define this as a serious of ParserRule and ParserRuleStep objects.
ParserRules define the rules we're going to compare text against, Each rule has a serios of steps.

```
var BuyAction = ParserRule{
   Name: "BuyStock",
    Steps: []ParserRuleStep{ }
}
var SellAction = ParserRule{
   Name: "SellStock",
   Steps: []ParserRuleStep{ }
}
var DisplayAction = ParserRule{
   Name: "DisplayPortfolio",
   Steps: []ParserRuleStep{ }
}
var SellAllAction = ParserRule{
   Name: "SellAll",
   Steps: []ParserRuleStep{ }
}
```

This defines four rules -- rules to buy stocks, sell stocks, display a portfolio and sell all material.
How we do this is defined by adding RuleSteps to reach rule.
A rule step defines how parsing is actually handled.  Using the BuyAction 
rule, it requires the following steps:
* Except a BUY string
* Expect an integer for the number of shares
* Expect the word SHARES
* Expect the word OF
* Expect a string for the stock name

We need to define a rule step for each step.  Something like this:
```aiignore
var BuyAction = ParserRule{
   Name: "BuyStock",
   Steps: []ParserRuleStep{
      {
      },
      {
      },
      {
      },
      {
      },
      {
      },
   }
}
```

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

# The Next Step - The Parser
Once we have a Lexer ready to go, we can give it to a Parser.  The parser object is
just.  When we create the object, we can provide three parameters:

* Debug - Do we enable debugging information
* Input - The input string the user gave us -- we pass this to the Lexer
* Exucde - A set of strings we don't want to see in our parse -- we pass this to the Lexer

Once the parser object is created, we can call Parse on it, passing our parsing rules. and a
pointer to an opaque data object. 
The parser will attempt to parse the input 
using all our rules, and return a final status.  Along the way, each rule
may or may not alter the data object we passed by reference.

We receive one of three status values:

* PARSE_SUCCESS -- we successfully parsed the input with a rule
* PARSE_FAILURE -- we failed to find any rule that could successfully parse this input
* SKIP_ON FAILURE -- In some cases, the rules may fail, but we know that other processing can find a result.  Skip this failure.

# The ParseRule object
A ParseRule object is container for all rule steps needed to complelte that rule.  It 
is nothing more than the name of the rule (for debugging), and an array of ParserRuleSteps which 
will be executed in order much like firewall rules.  The first rule that successfully 
parse the input will return, even if later rules can do a better job of it.

# The ParserRuleStep object
A ParserRuleStep object is a single step in a rule.  It defines for each step of a rule, how to handle
input.  It accepts:

* A name (ex: "RuleStep1") for debugging
* A type, which is one of the types defined above (STRING, INTEGER, etc)
* SkipOnTypeError - A boolean that tells us, if the lexer fails to get the correct type, should we skip to the next rule or not
* Options - A set of bitmapped options to tell us if we need to convert the text to upper or lower case for example
* A user-define function -- this is important.  Assuming we pass the type check, we pass the object and its type, to this function to say "I got this -- is it good?"
The function will parse that data, update the opaque data object if necessary, and 
return one of our status values (SUCCESS, SKIP_ON_FAILURE or FAIL).

