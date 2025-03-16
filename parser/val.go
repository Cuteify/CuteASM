package parser

import (
	"CuteASM/lexer"
	"strconv"
	"strings"
)

const (
	NUMBER = 1
	STRING = 2
	ADDR   = 3
	VAR    = 4
	PSEUDO = 5
	REG    = 6
)

type MemoryAddr struct {
	BaseReg      *Reg    // 基址寄存器（如rax）
	IndexReg     *Reg    // 变址寄存器（如rbx）
	Scale        int     // 比例因子（1/2/4/8）
	Displacement float64 // 位移值（如0x100）
	LabelRef     string  // 标签引用（如array_base）
	Length       int     // 数据长度（1/2/4/8）
}

type Reg struct {
	Name string
	Num  int
}

type Value struct {
	Addr   *MemoryAddr
	Reg    *Reg
	Var    *VarBlock
	String string
	Pseudo string
	Num    float64
	Type   int
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
				if v.isDigit(name[1:]) {
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
		v.Num, _ = v.parseNumber(tokens[0].Value)
		v.Type = NUMBER
	} else if len(tokens) == 1 && (tokens[0].Type == lexer.STRING || tokens[0].Type == lexer.CHAR) {
		v.String = tokens[0].Value
		v.Type = STRING
	} else if v.isMemoryAddress(tokens) {
		length := 0
		switch tokens[0].Value {
		case "BB":
			length = 1
		case "WW":
			length = 2
		case "DW":
			length = 4
		case "QW":
			length = 8
		}
		v.Addr = v.parseMemoryAddress(p, tokens[2:len(tokens)-1]) // 去掉方括号
		v.Addr.Length = length
		v.Type = ADDR
	}
}

func (v *Value) isDigit(str string) bool {
	strLength := len(str)
	for i := 0; i < strLength; i++ {
		if str[i] < '0' || str[i] > '9' {
			return false
		}
	}
	return true
}

// 判断是否是内存地址表达式
func (v *Value) isMemoryAddress(tokens []lexer.Token) bool {
	return len(tokens) >= 3 &&
		tokens[0].Type == lexer.PSEUDO &&
		tokens[1].Value == "[" &&
		tokens[len(tokens)-1].Value == "]"
}

func (v *Value) parseMemoryAddress(p *Parser, tokens []lexer.Token) *MemoryAddr {
	addr := &MemoryAddr{Scale: 1}

	// 分解表达式各部分
	parts := splitByOperators(tokens, []string{"+", "-", "*"})

	for _, part := range parts {
		switch {
		case part[0].Type == lexer.NAME && len(part) == 0: // 标签引用
			addr.LabelRef = strings.TrimSuffix(part[0].Value, ":")
		case containsRegister(part): // 寄存器处理
			reg := v.parseRegister(part, p)
			if len(part) == 3 {
				if addr.BaseReg == nil {
					addr.BaseReg = reg
				} else {
					addr.IndexReg = reg
				}
			} else if len(part) == 5 && part[1].Value == "*" { // 比例因子
				addr.IndexReg = reg
				scale, _ := strconv.Atoi(part[2].Value)
				addr.Scale = scale
			}
		default: // 数值位移
			num, err := v.parseNumber(part[0].Value)
			if err != nil {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
			addr.Displacement = num
		}
	}

	return addr
}

// 替换原有的strconv.ParseFloat调用
func (v *Value) parseNumber(str string) (num float64, err error) {
	if strings.HasPrefix(str, "0x") {
		tmp, errTmp := strconv.ParseInt(str[2:], 16, 64)
		num = float64(tmp)
		err = errTmp
	} else if strings.HasPrefix(str, "0b") {
		tmp, errTmp := strconv.ParseInt(str[2:], 2, 64)
		num = float64(tmp)
		err = errTmp
	} else if v.isDigit(str) {
		num, err = strconv.ParseFloat(str, 64)
	}
	return
}

func splitByOperators(tokens []lexer.Token, ops []string) [][]lexer.Token {
	var parts [][]lexer.Token
	var current []lexer.Token

	for _, token := range tokens {
		if contains(ops, token.Value) {
			if len(current) > 0 {
				parts = append(parts, current)
				current = nil
			}
		} else {
			current = append(current, token)
		}
	}
	if len(current) > 0 {
		parts = append(parts, current)
	}
	return parts
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func (v *Value) parseRegister(tokens []lexer.Token, p *Parser) (reg *Reg) {
	name := ""
	if tokens[0].Type == lexer.SEPARATOR && tokens[0].Value == "$" {
		tokens = tokens[1:]
		name += "$"
		for i := 0; i < len(tokens); i++ {
			if tokens[i].Type != lexer.SEPARATOR && tokens[i].Type != lexer.NAME {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
			name += tokens[i].Value
		}
		if len(tokens) >= 2 && tokens[0].Type == lexer.SEPARATOR && tokens[0].Value == "$" {
			name = name[2:]
			if name[0] == 'r' {
				if v.isDigit(name[1:]) {
					num, _ := strconv.Atoi(name[1:])
					reg = &Reg{Num: num}
				} else {
					reg = &Reg{Name: name[1:]}
				}
				return reg
			}
		}
	}
	return
}

func containsRegister(tokens []lexer.Token) bool {
	if len(tokens) < 3 {
		return false
	}
	return tokens[0].Value == "$" || tokens[1].Value == "$" || tokens[2].Value[0] == 'r'
}
