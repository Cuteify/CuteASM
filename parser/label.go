package parser

import (
	"CuteASM/lexer"
	"CuteASM/utils"
	"fmt"
)

type LabelBlock struct {
	IsFunc    bool
	Args      []*ArgBlock
	VarOffset int
	ArgOffset int
	StackRoom int
	//Class      typeSys.Type
	//Return     []typeSys.Type
	Name string
	//BuildFlags []*Build
}

type ArgBlock struct {
	Name   string
	Length int
	Defind *ArgBlock
	Offset int
}

func (l *LabelBlock) Parse(tokens []lexer.Token, p *Parser) {
	l.Name = tokens[0].Value
	node := &Node{Value: l}
	if p.ThisBlock.Father != nil {
		if _, ok := p.ThisBlock.Value.(*LabelBlock); !ok || !p.ThisBlock.Value.(*LabelBlock).IsFunc {
			p.Back(1)
		}
	}
	if len(tokens) >= 4 && tokens[2].Type == lexer.SEPARATOR && tokens[2].Value == "(" {
		l.IsFunc = true
		if p.ThisBlock.Father != nil {
			p.Back(1)
		}
		tokens = tokens[3:]
		for {
			if len(tokens) < 1 {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
			code := tokens[0]
			if len(tokens) > 2 && code.Type == lexer.PSEUDO && utils.GetLength(code.Value) != 0 {
				l.Args = append(l.Args, &ArgBlock{
					Length: utils.GetLength(code.Value),
				})
				code = tokens[1]
				fmt.Println(code)
				if code.Type == lexer.NAME && !p.isInstructions(code) {
					l.Args[len(l.Args)-1].Name = code.Value
					code = tokens[2]
					if code.Type != lexer.SEPARATOR || (code.Value != "," && code.Value != ")") {
						p.Error.MissError("", p.Lexer.Cursor, "")
					}
					if code.Type == lexer.SEPARATOR && code.Value == ")" {
						break
					}
					tokens = tokens[3:]
				} else {
					p.Error.MissError("", p.Lexer.Cursor, "")
				}
			} else if code.Type == lexer.SEPARATOR && code.Value == ")" {
				break
			} else {
				p.Error.MissError("", p.Lexer.Cursor, code.Value)
			}
		}
	}
	p.ThisBlock.AddChild(node)
	p.ThisBlock = node
}
