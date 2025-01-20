package parser

import (
	"CuteASM/lexer"
	"strconv"
)

const (
	NUMBER = 1
	STRING = 2
	ADDR   = 3
	VAR    = 4
	PSEUDO = 5
	REG    = 6
)

type Instruction struct {
	Instruction string
	Args        []*Value
}

type Addr struct {
	Name string
	Exp  []string
}

type Reg struct {
	Name string
	Num  int
}

type Value struct {
	Addr   *Addr
	Reg    *Reg
	Var    *VarBlock
	String string
	Pseudo string
	Num    float64
	Type   int
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

func (v *Value) Parse(p *Parser, tokens []lexer.Token) {
	if tokens[0].Type == lexer.SEPARATOR && tokens[0].Value == "$" {
		v.Var = &VarBlock{}
		v.Var.Name += "$"
		tokens = tokens[1:]
		for i := 0; i < len(tokens); i++ {
			if tokens[i].Type != lexer.SEPARATOR && tokens[i].Type != lexer.NAME {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
			v.Var.Name += tokens[i].Value
		}
		if len(tokens) >= 2 && tokens[0].Type == lexer.SEPARATOR && tokens[0].Value == "$" {
			name := v.Var.Name[2:]
			if name[0] == 'r' {
				if IsDigit(name[1:]) {
					num, _ := strconv.Atoi(name[1:])
					v.Reg = &Reg{Num: num}
				} else {
					v.Reg = &Reg{Name: name[1:]}
				}
				v.Var = nil
				v.Type = REG
				return
			}
		}
		v.Type = VAR
	} else if len(tokens) == 1 && tokens[0].Type == lexer.PSEUDO {
		v.Pseudo = tokens[0].Value
		v.Type = PSEUDO
	} else if len(tokens) == 1 && tokens[0].Type == lexer.NUMBER {
		v.Num, _ = strconv.ParseFloat(tokens[0].Value, 64)
		v.Type = NUMBER
	} else if len(tokens) == 1 && (tokens[0].Type == lexer.STRING || tokens[0].Type == lexer.CHAR) {
		v.String = tokens[0].Value
		v.Type = STRING
	}
}

func IsDigit(str string) bool {
	strLength := len(str)
	for i := 0; i < strLength; i++ {
		if str[i] < '0' || str[i] > '9' {
			return false
		}
	}
	return true
}
