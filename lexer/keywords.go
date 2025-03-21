package lexer

// 关键字
var (
	keywords = map[string]int{
		";":             1,
		"(":             1,
		")":             1,
		"[":             1,
		"]":             1,
		",":             1,
		" ":             1,
		"+":             1,
		"-":             1,
		"*":             1,
		"/":             1,
		":":             1,
		"\"":            1,
		"'":             1,
		"$":             1,
		"%":             1,
		"\r":            1,
		"\n":            1,
		"\t":            1,
		"ADD":           7,
		"ADC":           7,
		"AND":           7,
		"CALL":          7,
		"CMP":           7,
		"CMPXCHG":       7,
		"DEC":           7,
		"DIV":           7,
		"IDIV":          7,
		"IMUL":          7,
		"INC":           7,
		"JAE/JNB":       7,
		"JBE/JNA":       7,
		"JC":            7,
		"JCXZ":          7,
		"JE/JZ":         7,
		"JECXZ":         7,
		"JG/JNLE":       7,
		"JGE/JNL":       7,
		"JL/JNGE":       7,
		"JLE/JNG":       7,
		"JMP":           7,
		"JNAE/JB":       7,
		"JNBE/JA":       7,
		"JNE/JNZ":       7,
		"JNO":           7,
		"JNP/JPO":       7,
		"JNS":           7,
		"JO":            7,
		"JP/JPE":        7,
		"JS":            7,
		"LOOP":          7,
		"LOOPE/LOOPZ":   7,
		"LOOPNE/LOOPNZ": 7,
		"MOV":           7,
		"MOVSX":         7,
		"MOVZX":         7,
		"MUL":           7,
		"NEG":           7,
		"NOT":           7,
		"OR":            7,
		"POP":           7,
		"POPA":          7,
		"POPAD":         7,
		"PUSH":          7,
		"PUSHA":         7,
		"PUSHAD":        7,
		"RET/RETF":      7,
		"ROL":           7,
		"RET":           7,
		"ROR":           7,
		"RCL":           7,
		"RCR":           7,
		"SAR":           7,
		"SAL":           7,
		"SHL":           7,
		"SHR":           7,
		"SUB":           7,
		"SBB":           7,
		"TEST":          7,
		"XADD":          7,
		"XOR":           7,
		"XLAT":          7,
		"IN":            7,
		"OUT":           7,
		"LEA":           7,
		"LDS":           7,
		"LES":           7,
		"LFS":           7,
		"LGS":           7,
		"LSS":           7,
		"LAHF":          7,
		"SAHF":          7,
		"PUSHF":         7,
		"POPF":          7,
		"PUSHD":         7,
		"POPD":          7,
		"BSWAP":         7,
		"XCHG":          7,
		"AAD":           7,
		"AAM":           7,
		"AAS":           7,
		"DAA":           7,
		"DAS":           7,
		"AAA":           7,
		"SECTION":       8,
		"GLOBAL":        8,
		"BB":            8,
		"WW":            8,
		"DW":            8,
		"QW":            8,
		"VAR":           8,
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
