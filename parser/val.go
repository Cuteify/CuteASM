package parser

import (
	"CuteASM/arch/types"
	"CuteASM/lexer"
	"CuteASM/utils"
	"strconv"
	"strings"
)

// 常量定义指令操作数类型
const (
	NUMBER = iota + 1 // 数值字面量类型
	STRING            // 字符串字面量类型
	ADDR              // 内存地址类型
	VAR               // 变量引用类型containsRegister(tokens)
	PSEUDO            // 伪指令类型
	REG               // 寄存器类型
	LABEL             // 标签类型
)

// MemoryAddr 表示汇编指令中的内存地址操作数
type MemoryAddr struct {
	BaseReg      *Reg    // 基址寄存器（如rax）
	IndexReg     *Reg    // 变址寄存器（如rbx）
	Scale        int     // 比例因子（1/2/4/8）
	Displacement float64 // 位移值（如0x100）
	LabelRef     string  // 标签引用（如array_base）
	Length       int     // 数据长度（1/2/4/8）
}

// Reg 表示寄存器操作数
type Reg struct {
	Name string // 寄存器名称（如"ax"）
	Num  int    // 寄存器编号
	Type int
}

// Value 表示汇编指令中的操作数
type Value struct {
	Addr   *MemoryAddr // 内存地址操作数
	Reg    *Reg        // 寄存器操作数
	Var    *VarBlock   // 变量操作数
	String string      // 字符串值
	Pseudo string      // 伪指令名称
	Num    float64     // 数值
	Type   int         // 操作数类型（使用上述常量定义）
}

// Parse 解析token序列为操作数
// 参数：
//
//	p: 解析器实例
//	tokens: 待解析的token序列
func (v *Value) Parse(p *Parser, tokens []lexer.Token) {
	if tokens[0].Type == lexer.SEPARATOR && tokens[0].Value == "$" {
		// 处理变量引用（$开头的标识符）
		v.ParseVar(p, tokens)
		v.Type = VAR
	} else if len(tokens) == 1 && tokens[0].Type == lexer.PSEUDO {
		// 处理伪指令
		v.Pseudo = tokens[0].Value
		v.Type = PSEUDO
	} else if len(tokens) == 1 && tokens[0].Type == lexer.NUMBER {
		// 处理数字字面量
		v.Num, _ = v.parseNumber(tokens[0].Value)
		v.Type = NUMBER
	} else if len(tokens) == 1 && (tokens[0].Type == lexer.STRING || tokens[0].Type == lexer.CHAR) {
		// 处理字符串或字符字面量
		v.String = tokens[0].Value
		v.Type = STRING
	} else if v.isMemoryAddress(tokens) {
		// 处理内存地址表达式（如[BB [rax+0x10]]）
		v.Addr = v.parseMemoryAddress(p, tokens[2:len(tokens)-1]) // 去掉方括号
		v.Addr.Length = utils.GetLength(tokens[0].Value)
		v.Type = ADDR
	} else if containsRegister(tokens) {
		// 处理寄存器操作数
		v.Reg = v.parseRegister(tokens, p)
		v.Type = REG
	} else if isLabel(tokens) {
		// 处理标签引用
		v.String = parseLabel(tokens)
		v.Type = LABEL
	}
}

// isDigit 检查字符串是否全由数字组成
func (v *Value) isDigit(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] < '0' || str[i] > '9' {
			return false
		}
	}
	return true
}

// isMemoryAddress 判断token序列是否表示内存地址
// 内存地址格式: [前缀] [表达式]
func (v *Value) isMemoryAddress(tokens []lexer.Token) bool {
	return len(tokens) >= 3 &&
		tokens[0].Type == lexer.PSEUDO &&
		tokens[1].Value == "[" &&
		tokens[len(tokens)-1].Value == "]"
}

// parseMemoryAddress 解析内存地址表达式
func (v *Value) parseMemoryAddress(p *Parser, tokens []lexer.Token) *MemoryAddr {
	addr := &MemoryAddr{Scale: 1}
	// 按运算符分割表达式
	parts := splitByOperators(tokens, []string{"+", "-", "*"})

	for _, part := range parts {
		switch {
		case isLabelRef(part):
			// 处理标签引用部分
			addr.LabelRef = strings.TrimSuffix(part[0].Value, ":")
		case isRegisterPart(part):
			// 处理寄存器部分
			v.handleRegisterPart(part, p, addr)
		default:
			// 处理位移数值部分
			v.handleDisplacementPart(part, p, addr)
		}
	}
	return addr
}

// isLabelRef 判断token序列是否为标签引用
func isLabelRef(part []lexer.Token) bool {
	return len(part) == 1 &&
		part[0].Type == lexer.NAME &&
		strings.HasSuffix(part[0].Value, ":")
}

// isRegisterPart 判断token序列是否包含寄存器
func isRegisterPart(part []lexer.Token) bool {
	return containsRegister(part)
}

// handleRegisterPart 处理寄存器表达式部分
func (v *Value) handleRegisterPart(part []lexer.Token, p *Parser, addr *MemoryAddr) {
	reg := v.parseRegister(part, p)
	if reg == nil {
		return
	}

	// 根据表达式结构设置基址/变址寄存器
	switch {
	case len(part) == 3:
		if addr.BaseReg == nil {
			addr.BaseReg = reg
		} else {
			addr.IndexReg = reg
		}
	case len(part) == 5 && part[1].Value == "*":
		addr.IndexReg = reg
		scale, _ := strconv.Atoi(part[2].Value)
		addr.Scale = scale
	}
}

// handleDisplacementPart 处理数值位移部分
func (v *Value) handleDisplacementPart(part []lexer.Token, p *Parser, addr *MemoryAddr) {
	if len(part) == 0 {
		return
	}

	num, err := v.parseNumber(part[0].Value)
	if err != nil {
		p.Error.MissError("", p.Lexer.Cursor, "")
		return
	}
	addr.Displacement = num
}

// parseNumber 解析不同进制的数值字符串
func (v *Value) parseNumber(str string) (num float64, err error) {
	// 处理十六进制
	if strings.HasPrefix(str, "0x") {
		tmp, errTmp := strconv.ParseInt(str[2:], 16, 64)
		num = float64(tmp)
		err = errTmp
		// 处理二进制
	} else if strings.HasPrefix(str, "0b") {
		tmp, errTmp := strconv.ParseInt(str[2:], 2, 64)
		num = float64(tmp)
		err = errTmp
		// 处理十进制
	} else if v.isDigit(str) {
		num, err = strconv.ParseFloat(str, 64)
	}
	return
}

// splitByOperators 按指定运算符分割token序列
func splitByOperators(tokens []lexer.Token, ops []string) [][]lexer.Token {
	var parts [][]lexer.Token
	var current []lexer.Token

	// 遍历token进行分割
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

// contains 检查字符串是否在列表中
func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

// parseRegister 解析寄存器token序列
func (v *Value) parseRegister(tokens []lexer.Token, p *Parser) (reg *Reg) {
	if containsRegister(tokens) {
		lengthP := strings.ToUpper(tokens[1].Value[:1])
		regType := 0
		switch lengthP {
		case "R":
			regType = types.Reg64
		case "E":
			regType = types.Reg32
		case "N":
			regType = types.Reg16
		case "L":
			regType = types.Reg8
		case "D":
			regType = types.RegDR
		case "C":
			regType = types.RegCR
		case "T":
			regType = types.RegTR
		case "F":
			regType = types.RegFPU
		case "S":
			regType = types.RegSEG
		case "B":
			regType = types.RegBND
		case "X":
			regType = types.RegXMM
		case "Y":
			regType = types.RegYMM
		case "Z":
			regType = types.RegZMM
		case "M":
			regType = types.RegMMX
		case "A":
			regType = types.RegTMM
		}
		name := tokens[1].Value
		if v.isDigit(name[1:]) {
			num, _ := strconv.Atoi(name[1:])
			reg = &Reg{Num: num}
		} else {
			reg = &Reg{Name: name[1:]}
		}
		reg.Type = regType
		return reg
	}
	return
}

// containsRegister 判断token序列是否包含寄存器
func containsRegister(tokens []lexer.Token) bool {
	if len(tokens) < 2 {
		return false
	}
	return (tokens[0].Value == "%" && tokens[0].Type == lexer.SEPARATOR) && (tokens[1].Type == lexer.NAME)
}

// isLabel 判断token序列是否构成标签
func isLabel(tokens []lexer.Token) bool {
	for _, token := range tokens {
		if token.Type != lexer.NAME && token.Type != lexer.SEPARATOR {
			return false
		}
	}
	return true
}

// parseLabel 从token序列构建标签字符串
func parseLabel(tokens []lexer.Token) (label string) {
	for _, token := range tokens {
		label += token.Value
	}
	return
}

// ParseVar 解析变量引用操作数
func (v *Value) ParseVar(p *Parser, tokens []lexer.Token) {
	v.Var = &VarBlock{}
	v.Var.Name += "$"
	tokens = tokens[1:]
	// 拼接变量名各部分
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type != lexer.SEPARATOR && tokens[i].Type != lexer.NAME {
			p.Error.MissError("", p.Lexer.Cursor, "")
		}
		v.Var.Name += tokens[i].Value
	}
	v.Type = ADDR
}

func valueToMemoryAddr(arg *Value) {
	arg.Addr = &MemoryAddr{}
	arg.Addr.Length = arg.Var.Length
	// 获取变量的偏移
	arg.Addr.Displacement = float64(arg.Var.Offset)
	arg.Var = nil
	arg.Type = ADDR
}
