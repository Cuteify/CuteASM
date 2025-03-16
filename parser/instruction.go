package parser

import (
	"CuteASM/lexer"
)

type Instruction struct {
	Instruction string
	Args        []*Value
}

func (i *Instruction) ParseInstruction(p *Parser) {
	tmp := []lexer.Token{}
	lastIsSep := true
	for {
		code := p.Lexer.Next()
		if code.Type == lexer.SEPARATOR && (code.Value == "," || code.Value == "\n" || code.Value == "\r") {
			lastIsSep = false
			if len(tmp) != 0 {
				val := &Value{}
				val.Parse(p, tmp)
				i.Args = append(i.Args, val)
				tmp = tmp[:0]
			}
			if code.Value == "\n" || code.Value == "\r" {
				break
			}
			if lastIsSep {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
		} else {
			lastIsSep = false
			tmp = append(tmp, code)
		}
	}
}

func (i *Instruction) Parse(p *Parser) {
	i.ParseInstruction(p)
	for e := 0; e < len(i.Args); e++ {
		arg := i.Args[e]
		if arg.Type == VAR {
			arg.Var.FindDefind(p)
		}
	}
	p.ThisBlock.AddChild(&Node{Value: i})
}
