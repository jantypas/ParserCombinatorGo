package ParserCore

import "strings"

func (p *ParserObject) StripInputNoise(striptokens []string) {
	// This function strips unwanted tokens from the input string
	for _, token := range striptokens {
		p.Input = strings.ReplaceAll(p.Input, token, " ")
	}
}

func (p *ParserObject) SetDebug(state bool) {
	p.Debug = state
}
