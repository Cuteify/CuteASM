package lexer

// 关键字
var (
	keywords = map[string]int{
		";":       1,
		"(":       1,
		")":       1,
		"[":       1,
		"]":       1,
		",":       1,
		" ":       1,
		"+":       1,
		"-":       1,
		"*":       1,
		"/":       1,
		":":       1,
		"\"":      1,
		"'":       1,
		"$":       1,
		"%":       1,
		"\r":      1,
		"\n":      1,
		"\t":      1,
		"MOV":     7,
		"MOVSX":   7,
		"MOVZX":   7,
		"ADD":     7,
		"SUB":     7,
		"INC":     7,
		"DEC":     7,
		"MUL":     7,
		"IMUL":    7,
		"DIV":     7,
		"IDIV":    7,
		"AND":     7,
		"OR":      7,
		"XOR":     7,
		"NOT":     7,
		"SHL":     7,
		"SHR":     7,
		"SAL":     7,
		"SAR":     7,
		"JMP":     7,
		"JE":      7,
		"JNE":     7,
		"JG":      7,
		"JL":      7,
		"JGE":     7,
		"JLE":     7,
		"CALL":    7,
		"RET":     7,
		"CMP":     7,
		"TEST":    7,
		"INT":     7,
		"HLT":     7,
		"NOP":     7,
		"PUSH":    7,
		"POP":     7,
		"SECTION": 8,
		"GLOBAL":  8,
		"BB":      8,
		"WW":      8,
		"DW":      8,
		"QW":      8,
		"VAR":     8,
	}

	// LexToken类型(反查用)
	LexTokenType = map[string]int{
		"SEPARATOR":   0x1,
		"STRING":      0x2,
		"NUMBER":      0x3,
		"NAME":        0x4,
		"CHAR":        0x5,
		"TYPE":        0x6,
		"INSTRUCTION": 0x7,
		"PSEUDO":      0x8,
	}
)

const (
	// LexToken类型
	SEPARATOR   = 0x1
	STRING      = 0x2
	NUMBER      = 0x3
	NAME        = 0x4
	CHAR        = 0x5
	TYPE        = 0x6
	INSTRUCTION = 0x7
	PSEUDO      = 0x8
)
