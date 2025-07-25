package ParserCore

func (p *ParserObject) SetDebug(state bool) {
	p.Debug = state
}

func IfDebug(debug bool, f func(format string, a ...interface{}) (int, error), format string, a ...interface{}) {
	if debug {
		_, err := f(format, a...)
		if err != nil {
			return
		}
	}
}

var GreenText = "\033[32m"
var RedText = "\033[31m"
var BlueText = "\033[34m"
var ResetText = "\033[0m"
