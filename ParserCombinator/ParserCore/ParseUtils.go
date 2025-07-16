package ParserCore

import (
	"fmt"
	"strings"
)

func (p *ParserObject) StripInputNoise(striptokens []string) {
	// This function strips unwanted tokens from the input string
	for _, token := range striptokens {
		p.Input = strings.ReplaceAll(p.Input, token, "")
		if p.Debug {
			fmt.Println("Stripped token:", token, "leaving", p.Input)
		}
	}
}

func (p *ParserObject) SetDebug(state bool) {
	p.Debug = state
}
