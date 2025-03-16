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

func (l *LabelBlock) Parse(p *Parser) {
	beforeCursor := p.Lexer.Cursor
	code := p.Lexer.Next()
	node := &Node{Value: l}
	if p.ThisBlock.Father != nil {
		if _, ok := p.ThisBlock.Value.(*LabelBlock); !ok || !p.ThisBlock.Value.(*LabelBlock).IsFunc {
			p.Back(1)
		}
	}
	if code.Type == lexer.SEPARATOR && code.Value == "(" {
		l.IsFunc = true
		if p.ThisBlock.Father != nil {
			p.Back(1)
		}
		for {
			code := p.Lexer.Next()
			if code.Type == lexer.PSEUDO && utils.GetLength(code.Value) != 0 {
				l.Args = append(l.Args, &ArgBlock{
					Length: utils.GetLength(code.Value),
				})
				code := p.Lexer.Next()
				fmt.Println(code)
				if code.Type == lexer.NAME {
					l.Args[len(l.Args)-1].Name = code.Value
					code := p.Lexer.Next()
					if code.Type != lexer.SEPARATOR || (code.Value != "," && code.Value != ")") {
						p.Error.MissError("", p.Lexer.Cursor, "")
					}
					if code.Type == lexer.SEPARATOR && code.Value == ")" {
						break
					}
				} else {
					p.Error.MissError("", p.Lexer.Cursor, "")
				}
			} else if code.Type == lexer.SEPARATOR && code.Value == ")" {
				break
			} else {
				p.Error.MissError("", p.Lexer.Cursor, code.Value)
			}
		}
	} else {
		p.Lexer.Cursor = beforeCursor
	}
	p.ThisBlock.AddChild(node)
	p.ThisBlock = node
}
