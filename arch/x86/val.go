// Package x86 implements x86 architecture specific code
package x86

import (
	"CuteASM/arch/types"
	"CuteASM/parser"
)

// 判断操作数是否为寄存器
func IsReg(arg *parser.Value, bits int) bool {
	if arg.Type != parser.REG || arg.Reg == nil {
		return false
	}
	// 根据寄存器名称判断位数
	switch bits {
	case 8:
		return len(arg.Reg.Name) == 1 && (arg.Reg.Name[0] == 'l' || arg.Reg.Name[0] == 'h') // 8位寄存器如al, ah
	case 16:
		return len(arg.Reg.Name) == 2 && arg.Reg.Name[0] == 'a' // 16位寄存器如ax
	case 32:
		return arg.Reg.Name == "eax" // 32位寄存器
	case 64:
		return arg.Reg.Name == "rax" // 64位寄存器
	}
	return false
}

// 判断操作数是否为立即数
func IsImm(arg *parser.Value) bool {
	return arg.Type == parser.NUMBER
}

// 判断操作数是否为内存地址
func IsMem(arg *parser.Value) bool {
	return arg.Type == parser.ADDR
}

// 将Value类型转换为Operand类型
func ValueToOperand(arg *parser.Value) types.Operand {
	switch arg.Type {
	case parser.REG:
		reg := arg.Reg
		if reg == nil {
			return OpNone
		}

		// 根据寄存器类型和名称确定具体操作数类型
		switch reg.Type {
		case types.Reg8:
			return OpReg8
		case types.Reg16:
			return OpReg16
		case types.Reg32:
			return OpReg32
		case types.Reg64:
			return OpReg64
		case types.RegXMM:
			return OpRegXMM
		case types.RegYMM:
			return OpRegYMM
		case types.RegZMM:
			return OpRegZMM
		case types.RegMMX:
			return OpMMX
		case types.RegFPU:
			return OpFPU
		case types.RegTMM:
			return OpTMM
		case types.RegBND:
			return OpBND
		case types.RegSEG:
			// 段寄存器特殊处理
			switch reg.Name {
			case "es":
				return OpSegES
			case "cs":
				return OpSegCS
			case "ss":
				return OpSegSS
			case "ds":
				return OpSegDS
			case "fs":
				return OpSegFS
			case "gs":
				return OpSegGS
			default:
				return OpSysReg
			}
		case types.RegTR:
			// 测试寄存器特殊处理
			if reg.Name == "tr6" {
				return OpTR6
			} else if reg.Name == "tr7" {
				return OpTR7
			}
			return OpSysReg
		case types.RegDR:
			// 调试寄存器特殊处理
			switch reg.Name {
			case "dr0":
				return OpDR0
			case "dr1":
				return OpDR1
			case "dr2":
				return OpDR2
			case "dr3":
				return OpDR3
			case "dr6":
				return OpDR6
			case "dr7":
				return OpDR7
			default:
				return OpSysReg
			}
		case types.RegCR:
			// 控制寄存器特殊处理
			switch reg.Name {
			case "cr0":
				return OpCR0
			case "cr2":
				return OpCR2
			case "cr3":
				return OpCR3
			case "cr4":
				return OpCR4
			case "cr8":
				return OpCR8
			default:
				return OpSysReg
			}
		default:
			return OpSysReg
		}

	case parser.NUMBER:
		// 根据数值大小确定立即数类型
		num := arg.Num
		if num >= -128 && num <= 127 {
			return OpImm8 | OpImm16 | OpImm32 | OpImm64
		} else if num >= -32768 && num <= 32767 {
			return OpImm16 | OpImm32 | OpImm64
		} else if num >= -2147483648 && num <= 2147483647 {
			return OpImm32 | OpImm64
		}
		return OpImm64

	case parser.ADDR:
		// 根据内存地址长度确定内存类型
		if arg.Addr != nil {
			switch arg.Addr.Length {
			case 1:
				return OpMem8
			case 2:
				return OpMem16
			case 4:
				return OpMem32
			case 8:
				return OpMem64
			case 16:
				return OpMem128
			case 32:
				return OpMem256
			case 64:
				return OpMem512
			}
		}
		return OpMem

	case parser.STRING:
		return OpLabel

	case parser.PSEUDO:
		return OpNone

	case parser.LABEL:
		return OpLabel
	}

	return OpNone
}
