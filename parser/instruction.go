package parser

import (
	"CuteASM/arch/types"
	"CuteASM/lexer"
	"strings"
)

type Instruction struct {
	Instruction types.Instruction
	Args        []*Value
	OpSize      int // for backend
}

func (i *Instruction) ParseInstruction(tokens []lexer.Token, p *Parser) {
	// 解析名称和后面的一个空格
	i.Instruction = types.Instruction(strings.ToUpper(tokens[0].Value))
	if len(tokens) <= 1 {
		p.ThisBlock.AddChild(&Node{Value: i})
		return
	}
	tokens = tokens[1:]
	lastCursor := 0
	for e := 0; e < len(tokens); e++ {
		code := tokens[e]
		if code.Type == lexer.SEPARATOR && (code.Value == "," || code.Value == "\n" || code.Value == "\r") {
			val := &Value{}
			val.Parse(p, tokens[lastCursor:e])
			i.Args = append(i.Args, val)
			if code.Value == "\n" || code.Value == "\r" {
				return
			}
			lastCursor = e + 1
		}
	}
	val := &Value{}
	val.Parse(p, tokens[lastCursor:])
	i.Args = append(i.Args, val)
	if len(i.Args) > 4 {
		p.Lexer.Error.MissError("Syntax Error", p.Lexer.Cursor, "Too many arguments")
	}
}

func (i *Instruction) Parse(tokens []lexer.Token, p *Parser) {
	i.ParseInstruction(tokens, p)
	for e := 0; e < len(i.Args); e++ {
		arg := i.Args[e]
		if arg.Type == VAR {
			arg.Var.FindDefind(p)
			valueToMemoryAddr(arg)
		}
	}
	p.ThisBlock.AddChild(&Node{Value: i})
}

func (i *Instruction) IsBuiltin() bool {
	for _, ins := range types.BuiltinI {
		if ins == i.Instruction {
			return true
		}
	}
	return false
}
