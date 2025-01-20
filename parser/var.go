package parser

import (
	"CuteASM/utils"
	"fmt"
)

type VarBlock struct {
	Name   string
	Offset int
	Length int
}

func (v *VarBlock) Parse(p *Parser) {
	instruction := &Instruction{}
	instruction.ParseInstruction(p)
	if len(instruction.Args) == 2 {
		if len(instruction.Args[0].Var.Name) > 2 && instruction.Args[0].Var.Name[1] == '$' {
			switch instruction.Args[0].Var.Name[2:] {
			case "stackRoom":
				oldThisBlock := p.ThisBlock
				for {
					if p.ThisBlock.Father == nil {
						p.Error.MissError("", p.Lexer.Cursor, "")
					}
					switch p.ThisBlock.Value.(type) {
					case *LabelBlock:
						lb := p.ThisBlock.Value.(*LabelBlock)
						if lb.IsFunc {
							lb.StackRoom += int(instruction.Args[1].Num)
							fmt.Println(instruction.Args[1].Num)
							goto end2
						}
						fmt.Println(lb.StackRoom)
					}
					p.Back(1)
				}
			end2:
				p.ThisBlock = oldThisBlock
			}
			return
		} else {
			v.Name = instruction.Args[0].Var.Name
			v.Length = utils.GetLength(instruction.Args[1].Pseudo)
			if v.Length == 0 {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
		}
	}
	oldThisBlock := p.ThisBlock
	for {
		if p.ThisBlock.Father == nil {
			p.Error.MissError("", p.Lexer.Cursor, "")
		}
		switch p.ThisBlock.Value.(type) {
		case *LabelBlock:
			lb := p.ThisBlock.Value.(*LabelBlock)
			if lb.IsFunc {
				lb.VarOffset += v.Length
				v.Offset = -lb.VarOffset
				lb.StackRoom += v.Length
				goto end
			}
		}
		p.Back(1)
	}
end:
	p.ThisBlock = oldThisBlock
	node := &Node{
		Value: v,
	}
	p.ThisBlock.AddChild(node)
}

func (v *VarBlock) FindDefind(p *Parser) {
	oldThisBlock := p.ThisBlock
	for {
		if p.ThisBlock.Father == nil {
			p.Error.MissError("", p.Lexer.Cursor, "")
		}
		for i := len(p.ThisBlock.Children) - 1; i > 0; i-- {
			switch p.ThisBlock.Children[i].Value.(type) {
			case *VarBlock:
				varBlock := p.ThisBlock.Children[i].Value.(*VarBlock)
				v.Offset = varBlock.Offset
				v.Length = varBlock.Length
				goto end
			}
		}
		p.Back(1)
	}
end:
	p.ThisBlock = oldThisBlock
}
