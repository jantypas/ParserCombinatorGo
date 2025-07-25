# The Go Mini Parser Combinator Library

Author:  [John Antypas]
Email:   [jantypas@busygeeks.com]
License: [MIT License]
Website: [github.com/jantypas/go-mini-parser-combinator]
Version: [1.0.0]


# What is the mini parser combinator library and what is it for?

Many times during programming, we find we have to parse text from a user. 
Many times the text is simple in nature such as when we're building a CLI, 
and we can simply parse it with standard code.  However, in some cases,
we find ourselves needed to parse complex text streams.  During these
tasks we are often forced to write a parser of some type.  There are many
ways to write a parser, but they all assume a grammar of some type.

Let us assume we are writing a stock trading application that allows the user to:
* Buy a particular number of shares of a particular stock
* Sell a number of shares of a particular stock
* Display the user's holdings of a particular stock
* Display all of the user's holdings
* Liquid an entire stock back to cash
* Liquidate all holdings of all stocks

If we were writing a grammar, we'd typically write it in what is called the ENNF
format like this:

```
stockactions : buystock 
    | sellstock 
    | displaystock
    | displayportfolio
    | liquidatestock
    | liquidateportfolio
    ;

buystock : "buy" STRING INTEGER ; 
sellstock : "sell" STRING INTEGER ;
displaystock : "display" "stock" STRING ;
displayportfolio : "display" "portfolio" ;
liquidatestock : "liquidate" STRING ;
liquidateportfolio : "liquidate" "portfolio" ;

STRING : [a-zA-Z]+ ;
INTEGER : [0-9]+ ;
```

This is an example of an EBNF grammar similar to what the ANTLR tool expects.  ANTLR will take this
grammar, and generate a parser skeleton for us in langauges sucha as Java, C++ or Python.  However,
ANTLR is a complex tool, and it requires us to "staple" in code in the native lanague to actually 
interpret the data we parse.  ANTLR creates what is called the Abstract Syntax Tree or AST.
This tells us what the text should look look like, but, it doesn't tell us if the text is correct
in its meaning -- or rather, ANTLR doesn't help us with the semantic parser.

Put another way, ANTLR can take English text, and using a grammar, tall you if the syntax is 
correct.  Using our stock example, the syntax BUY Futzco -10 is syntactically correct, but
it makes no sense in actual fact as you can't by a negative number of shares of a non-existent stock.
The text is syntactically correct, but semantically incorrect.  ANTLR doesn't help us with this.

Futhermore, ANTLR and other parsers are meant for complex parsing tasks.  They produce relatively large code libraries
in a select set of languages.  There's nothing wrong with this if you need the full power of a parsing tool, but
for lighter asks, it's overkill.

The Go Mini Parser Combinator library is a simple parser combinator library 
that allows you to write lightweight parsers.  This example is written for the Go (Golang)
language, but it can be easily ported to other lanaguages such as Python, Java, C++ or Rust.
It uses no external libraries and is less than 400 lines of Go code, including debugging code.

# What types of objects can it parse?

"Out of the box", we understand a limited set of parsable objects:

* PARSE_ANY_STRING - Look for a single string without white space such as Bob or Apple.
* PARSE_ANY_INTEGER - Look for a single integer such as 123 or 456.
* PARSE_ANY_FLOAT - Look for a single float such as 123.45 or 456.78.
* PARSE_ANY_QUOTED_STRING - Look for a quoted string such as "Bob" or "Apple and Oranges".
* PARSE_STRING_CHOICE - Look for a string that matches one of a set of strings such as "buy", "sell", "display", etc.  If a match is found, this 
parse rule is successful.
* PARSE_STRING_LIST - Look for a list of string tokens that matches a list of provided tokens.  For example, given a token list iof 
{"BUY", "FUTZCO", "STOCK"}, this rule succeeds, if we receive the exact token strings of BUY FUTZCO STOCK.
* PARSE_COMMA, COLON, LESS_THAN, GREATER_THAN, EQUAL, QUESTION, PERCENT, EXCLAMATION - Look for the respective punctuation characters.

These simple type discriminators allow us to build more complex parsers. This is why we call it a parser combinator.
By combining these objects in series, we can build up complex parsing rules.  
In addition string token types can be automatically converted to upper or lower cases, or in case of
the STRING_CHOICE type, the match can be considered optional, meaning, if the match 
would normally fail, succeed anyway and move on.

Parsing rules are a tree objects.  First, we have a set of rules.  Each rule, much like a firewall rule, 
is given the inpuit and the rule is checked.  The rule can return either success (the rule matched), failur (the rule 
received a non-recoverable failure) or SKIP_RULE (the rule failed, but try the next rule).

For each rule, we have a series of steps.  Much like rules, each step checks the type of object, and it can return
success, failure, or SKIP_RULE< but it can also return SKIP_STEP, meaning the rule itself hasn't failed, move to the next 
step.

Once the type check is complete, we pass that step to a user-defined function.  That function has two purposes.  
First, it performs the semantic checks. second, it is passed a user-defined data 
structure.  It can modify the data structure based on what it finds.  This is how the 
output of the parser is captured.
