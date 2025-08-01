package types

type Register int
type Instruction string

// 修改后的指令映射结构
type InstructionMap map[Instruction]Opcode

type Architecture struct {
	RegisterList map[string]Register
	WordSize     int
	Instructions InstructionMap
}

type Opcode struct {
	Opcode   OpBytes
	Operands []Operands
}

type OpBytes []byte

var BuiltinI = []Instruction{
	"ADD",
	"AND",
	"CALL",
	"CMP",
	"DIV",
	"HALT",
	"JMP",
	"JMPN",
	"JMPZ",
	"LOAD",
	"MOV",
	"MUL",
	"NEG",
	"NOT",
	"OR",
	"POP",
	"PUSH",
	"RET",
	"SHIFTL",
	"SHIFTR",
	"STORE",
	"SUB",
	"XOR",
	"XCHG",
}

var Pseudo = []Instruction{
	"DATA",
	"TEXT",
	"GLOBAL",
	"EXTERN",
	"ENTRY",
	"FOR",
	"IF",
	"ELSE",
	"ENDIF",
	"CONST",
	"CONTINUE",
	"VAR",
	"SECTION",
	"BB",
	"WW",
	"DW",
	"QW",
	"TW",
	"OW",
	"YW",
	"ZW",
}

type Operands [4]Operand

type Operand int64

func (r Operand) Set(opt Operand) Operand {
	return r | opt
}

func (r Operand) Clear(opt Operand) Operand {
	return r &^ opt
}

func (r Operand) Has(opt Operand) bool {
	return (r&opt) != 0 || opt == r
}

func (r Operand) Is(opt Operand) bool {
	return r == opt
}

func (os Operands) Is(opts Operands) bool {
	return os[0].Is(opts[0]) && os[1].Is(opts[1]) && os[2].Is(opts[2]) && os[3].Is(opts[3])
}

func (os Operands) Has(opts Operands) bool {
	return os[0].Has(opts[0]) && os[1].Has(opts[1]) && os[2].Has(opts[2]) && os[3].Has(opts[3])
}

const (
	Reg8 = 1 << iota
	Reg16
	Reg32
	Reg64

	RegXMM
	RegYMM
	RegZMM

	RegMMX
	RegFPU
	RegTMM
	RegBND

	RegSEG
	RegTR
	RegDR
	RegCR

	RegFlag
)

// 定义操作码映射类型
type OpcodeMap map[Operands]OpBytes
