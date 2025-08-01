package lexer

// 关键字
var (
	keywords = map[string]int{
		";":        1,
		"(":        1,
		")":        1,
		"[":        1,
		"]":        1,
		",":        1,
		" ":        1,
		"+":        1,
		"-":        1,
		"*":        1,
		"/":        1,
		":":        1,
		"\"":       1,
		"'":        1,
		"$":        1,
		"%":        1,
		"\r":       1,
		"\n":       1,
		"\t":       1,
		"DATA":     8,
		"TEXT":     8,
		"GLOBAL":   8,
		"EXTERN":   8,
		"ENTRY":    8,
		"FOR":      8,
		"IF":       8,
		"ELSE":     8,
		"ENDIF":    8,
		"CONST":    8,
		"CONTINUE": 8,
		"VAR":      8,
		"SECTION":  8,
		"BB":       8,
		"WW":       8,
		"DW":       8,
		"QW":       8,
		"TW":       8,
		"OW":       8,
		"YW":       8,
		"ZW":       8,
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
