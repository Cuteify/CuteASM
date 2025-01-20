package nasm

import (
	"CuteASM/parser"
	"strconv"
	"strings"
)

var Regs = []string{
	"EAX",
	"EBX",
	"ECX",
	"EDX",
	"ESI",
	"EBI",
}

type Nasm struct {
	block *parser.Node
}

var count int = 0

func (nasm *Nasm) Compile(node *parser.Node) (code string) {
	for i := 0; i < len(node.Children); i++ {
		n := node.Children[i]
		switch n.Value.(type) {
		case *parser.SECTION:
			count = 0
			section := n.Value.(*parser.SECTION)
			code += Format("section " + section.Name + "; " + section.Desc)
			code += nasm.Compile(n)
		case *parser.LabelBlock:
			label := n.Value.(*parser.LabelBlock)
			if label.IsFunc {
				code += Format("; ==============================\n; Function:" + label.Name)
			}
			code += Format(label.Name + ":")
			count++
			if label.IsFunc {
				code += Format("PUSH ESP")
				code += Format("MOV EBP, ESP")
				code += Format("SUB ESP, " + strconv.Itoa(label.StackRoom))
			}
			code += nasm.Compile(n)
			count--
			if label.IsFunc {
				code += Format("\n; Function End:" + label.Name + "\n; ==============================\n")
			}
		case *parser.Instruction:
			instruction := n.Value.(*parser.Instruction)
			tmp := instruction.Instruction
			if len(instruction.Args) > 0 {
				tmp += " "
			}
			for i := 0; i < len(instruction.Args); i++ {
				tmp += ParseVal(instruction.Args[i]) + ", "
			}
			if len(instruction.Args) > 0 {
				tmp = tmp[:len(tmp)-2]
			}
			code += Format(tmp)
		}
	}
	return
}

func Format(text string) string {
	return strings.Repeat("    ", count) + text + "\n"
}

func ParseVal(v *parser.Value) (code string) {
	switch v.Type {
	case parser.VAR:
		code += getLengthName(v.Var.Length) + "[ebp" + strconv.Itoa(v.Var.Offset) + "]"
	case parser.NUMBER:
		code += strconv.FormatFloat(v.Num, 'f', -1, 64)
	case parser.REG:
		if v.Reg.Name == "" {
			v.Reg.Name = Regs[v.Reg.Num]
		}
		if v.Reg.Name == "sp" {
			v.Reg.Name = "esp"
		}
		if v.Reg.Name == "bp" {
			v.Reg.Name = "ebp"
		}
		code += v.Reg.Name
	case parser.STRING:
		code += "\"" + v.String + "\""
	}
	return
}

func getLengthName(size int) string {
	switch size {
	case 1:
		return "BYTE"
	case 2:
		return "WORD"
	case 4:
		return "DWORD"
	case 8:
		return "QWORD"
	default:
		return ""
	}
}
