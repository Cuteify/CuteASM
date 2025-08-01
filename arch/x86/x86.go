package x86

import (
	"CuteASM/arch/types"
	"CuteASM/parser"
)

// X86Builtin x86架构内置指令实现
type X86Builtin struct {
	arch *types.Architecture
}

// 添加空行解决EOF问题

// NewX86Builtin 创建x86内置指令实现实例
func NewX86Builtin(arch *types.Architecture) *X86Builtin {
	return &X86Builtin{arch: arch}
}

// Add 实现ADD指令
func (b *X86Builtin) Add(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, addOpMap)
}

// Mov 实现MOV指令
func (b *X86Builtin) Mov(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, movOpMap)
}

// ... 其他指令实现 ...

// RegisterList 获取寄存器列表
func RegisterList() map[string]types.Register {
	return RegLookup
}

// New 创建x86架构实例
func New() *types.Architecture {
	arch := &types.Architecture{
		RegisterList: RegLookup,
		WordSize:     32,
	}
	// 移除未使用的builtin变量初始化
	// arch.Builtin = NewX86Builtin(arch)
	return arch
}

// And 实现AND指令
func (b *X86Builtin) And(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, andOpMap)
}

// Call 实现CALL指令
func (b *X86Builtin) Call(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, callOpMap)
}

// Cmp 实现CMP指令
func (b *X86Builtin) Cmp(i *parser.Instruction) types.OpBytes {
	tmp := opMapHandler(i, cmpOpMap, true)
	if tmp == nil {
		if len(i.Args) != 2 {
			return nil // 参数数量错误
		}
		// 颠倒操作数顺序
		i.Args[0], i.Args[1] = i.Args[1], i.Args[0]
		return opMapHandler(i, cmpOpMap)
	}
	return tmp
}

// Div 实现DIV指令
func (b *X86Builtin) Div(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, divOpMap)
}

// Halt 实现HALT指令
func (b *X86Builtin) Halt(i *parser.Instruction) types.OpBytes {
	return []byte{0xF4} // HLT
}

// Jmp 实现JMP指令
func (b *X86Builtin) Jmp(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 1 {
		return nil // 参数数量错误
	}
	// 实现JMP指令编码
	if i.Args[0].Type == parser.NUMBER {
		return []byte{0xE9} // JMP imm32
	}
	return nil
}

// JmpNeg 实现JMPN指令
func (b *X86Builtin) JmpNeg(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 1 {
		return nil // 参数数量错误
	}
	// 实现JMPN指令编码
	if i.Args[0].Type == parser.NUMBER {
		return []byte{0x0F, 0x81} // JNO imm32
	}
	return nil
}

// JmpZero 实现JMPZ指令
func (b *X86Builtin) JmpZero(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 1 {
		return nil // 参数数量错误
	}
	// 实现JMPZ指令编码
	if i.Args[0].Type == parser.NUMBER {
		return []byte{0x0F, 0x84} // JE imm32
	}
	return nil
}

// Load 实现LOAD指令
func (b *X86Builtin) Load(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 2 {
		return nil // 参数数量错误
	}
	// 实现LOAD指令编码
	if IsReg(i.Args[0], 32) && IsMem(i.Args[1]) {
		return []byte{0x8B} // MOV EAX, [mem]
	}
	return nil
}

// Mul 实现MUL指令
func (b *X86Builtin) Mul(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, mulOpMap)
}

// Neg 实现NEG指令
func (b *X86Builtin) Neg(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 1 {
		return nil // 参数数量错误
	}
	// 实现NEG指令编码
	if IsReg(i.Args[0], 32) {
		return []byte{0xF7, 0xD0} // NEG EAX
	}
	return nil
}

// Not 实现NOT指令
func (b *X86Builtin) Not(i *parser.Instruction) types.OpBytes {
	//return opMapHandler(i, notOpMap)
	return nil
}

// Or 实现OR指令
func (b *X86Builtin) Or(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, orOpMap)
}

// Pop 实现POP指令
func (b *X86Builtin) Pop(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, popOpMap)
}

// Push 实现PUSH指令
func (b *X86Builtin) Push(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, pushOpMap)
}

// Ret 实现RET指令
func (b *X86Builtin) Ret(i *parser.Instruction) types.OpBytes {
	return []byte{0xC3} // RET
}

// ShiftL 实现SHIFTL指令
func (b *X86Builtin) ShiftL(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 1 {
		return nil // 参数数量错误
	}
	// 实现SHIFTL指令编码
	if IsReg(i.Args[0], 32) {
		return []byte{0xD1, 0xE0} // SHL EAX, 1
	}
	return nil
}

// ShiftR 实现SHIFTR指令
func (b *X86Builtin) ShiftR(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 1 {
		return nil // 参数数量错误
	}
	// 实现SHIFTR指令编码
	if IsReg(i.Args[0], 32) {
		return []byte{0xD1, 0xE8} // SHR EAX, 1
	}
	return nil
}

// Store 实现STORE指令
func (b *X86Builtin) Store(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 2 {
		return nil // 参数数量错误
	}
	// 实现STORE指令编码
	if IsReg(i.Args[0], 32) && IsMem(i.Args[1]) {
		return []byte{0x89} // MOV [mem], EAX
	}
	return nil
}

// Sub 实现SUB指令
func (b *X86Builtin) Sub(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, subOpMap)
}

// Xor 实现XOR指令
func (b *X86Builtin) Xor(i *parser.Instruction) types.OpBytes {
	return opMapHandler(i, xorOpMap)
}

// Xchg 实现XCHG指令
func (b *X86Builtin) Xchg(i *parser.Instruction) types.OpBytes {
	if len(i.Args) != 2 {
		return nil // 参数数量错误
	}
	// 实现XCHG指令编码
	if IsReg(i.Args[0], 32) && IsReg(i.Args[1], 32) {
		return []byte{0x87, 0xC0} // XCHG EAX, ECX
	}
	return nil
}

// MOV指令操作码映射 (优化版)
var movOpMap = types.OpcodeMap{
	// ======================
	// 通用寄存器传输指令
	// ======================
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x88}, // MOV reg8, reg8
	{OpReg16, OpReg16}: {0x89}, // MOV reg16, reg16
	{OpReg32, OpReg32}: {0x89}, // MOV reg32, reg32
	{OpReg64, OpReg64}: {0x89}, // MOV reg64, reg64

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x8A}, // MOV reg8, mem8
	{OpReg16, OpMem16}: {0x8B}, // MOV reg16, mem16
	{OpReg32, OpMem32}: {0x8B}, // MOV reg32, mem32
	{OpReg64, OpMem64}: {0x8B}, // MOV reg64, mem64

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x88}, // MOV mem8, reg8
	{OpMem16, OpReg16}: {0x89}, // MOV mem16, reg16
	{OpMem32, OpReg32}: {0x89}, // MOV mem32, reg32
	{OpMem64, OpReg64}: {0x89}, // MOV mem64, reg64

	// 立即数到寄存器 (需动态添加寄存器编码)
	{OpReg8, OpImm8}:   {0xB0}, // MOV reg8, imm8 (实际编码: B0 + r)
	{OpReg16, OpImm16}: {0xB8}, // MOV reg16, imm16 (实际编码: B8 + r)
	{OpReg32, OpImm32}: {0xB8}, // MOV reg32, imm32 (实际编码: B8 + r)
	{OpReg64, OpImm64}: {0xB8}, // MOV reg64, imm64 (实际编码: B8 + r)

	// 立即数到内存 (需要后续字节)
	{OpMem8, OpImm8}:   {0xC6}, // MOV mem8, imm8
	{OpMem16, OpImm16}: {0xC7}, // MOV mem16, imm16
	{OpMem32, OpImm32}: {0xC7}, // MOV mem32, imm32
	{OpMem64, OpImm32}: {0xC7}, // MOV mem64, imm32

	// ======================
	// 段寄存器传输指令
	// ======================
	// MOV sreg, mem
	{OpSegES, OpMem}: {0x8E}, // MOV ES, mem
	{OpSegCS, OpMem}: {0x8E}, // MOV CS, mem
	{OpSegSS, OpMem}: {0x8E}, // MOV SS, mem
	{OpSegDS, OpMem}: {0x8E}, // MOV DS, mem
	{OpSegFS, OpMem}: {0x8E}, // MOV FS, mem
	{OpSegGS, OpMem}: {0x8E}, // MOV GS, mem

	// MOV mem, sreg
	{OpMem, OpSegES}: {0x8C}, // MOV mem, ES
	{OpMem, OpSegCS}: {0x8C}, // MOV mem, CS
	{OpMem, OpSegSS}: {0x8C}, // MOV mem, SS
	{OpMem, OpSegDS}: {0x8C}, // MOV mem, DS
	{OpMem, OpSegFS}: {0x8C}, // MOV mem, FS
	{OpMem, OpSegGS}: {0x8C}, // MOV mem, GS

	// ======================
	// 控制寄存器传输指令
	// ======================
	// MOV CRn, reg (需要后续字节)
	{OpCR0, OpReg}: {0x0F, 0x22}, // MOV CR0, reg
	{OpCR2, OpReg}: {0x0F, 0x22}, // MOV CR2, reg
	{OpCR3, OpReg}: {0x0F, 0x22}, // MOV CR3, reg
	{OpCR4, OpReg}: {0x0F, 0x22}, // MOV CR4, reg
	{OpCR8, OpReg}: {0x0F, 0x22}, // MOV CR8, reg

	// MOV reg, CRn (需要后续字节)
	{OpReg, OpCR0}: {0x0F, 0x20}, // MOV reg, CR0
	{OpReg, OpCR2}: {0x0F, 0x20}, // MOV reg, CR2
	{OpReg, OpCR3}: {0x0F, 0x20}, // MOV reg, CR3
	{OpReg, OpCR4}: {0x0F, 0x20}, // MOV reg, CR4
	{OpReg, OpCR8}: {0x0F, 0x20}, // MOV reg, CR8

	// ======================
	// 调试寄存器传输指令
	// ======================
	// MOV DRn, reg (需要后续字节)
	{OpDR0, OpReg}: {0x0F, 0x23}, // MOV DR0, reg
	{OpDR1, OpReg}: {0x0F, 0x23}, // MOV DR1, reg
	{OpDR2, OpReg}: {0x0F, 0x23}, // MOV DR2, reg
	{OpDR3, OpReg}: {0x0F, 0x23}, // MOV DR3, reg
	{OpDR6, OpReg}: {0x0F, 0x23}, // MOV DR6, reg
	{OpDR7, OpReg}: {0x0F, 0x23}, // MOV DR7, reg

	// MOV reg, DRn (需要后续字节)
	{OpReg, OpDR0}: {0x0F, 0x21}, // MOV reg, DR0
	{OpReg, OpDR1}: {0x0F, 0x21}, // MOV reg, DR1
	{OpReg, OpDR2}: {0x0F, 0x21}, // MOV reg, DR2
	{OpReg, OpDR3}: {0x0F, 0x21}, // MOV reg, DR3
	{OpReg, OpDR6}: {0x0F, 0x21}, // MOV reg, DR6
	{OpReg, OpDR7}: {0x0F, 0x21}, // MOV reg, DR7

	// ======================
	// 向量寄存器传输指令
	// ======================
	// XMM
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0x6F}, // MOVDQA xmm, mem
	{OpMem128, OpRegXMM}: {0x66, 0x0F, 0x7F}, // MOVDQA mem, xmm

	// YMM
	{OpRegYMM, OpMem256}: {0xC5, 0xFD, 0x6F}, // VMOVDQA ymm, mem
	{OpMem256, OpRegYMM}: {0xC5, 0xFD, 0x7F}, // VMOVDQA mem, ymm

	// ZMM
	{OpRegZMM, OpMem512}: {0x62, 0xF1, 0x7D, 0x48}, // VMOVDQA32 zmm, mem
	{OpMem512, OpRegZMM}: {0x62, 0xF1, 0x7D, 0x48}, // VMOVDQA32 mem, zmm

	// 寄存器到寄存器
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0x6F},       // MOVDQA xmm, xmm
	{OpRegYMM, OpRegYMM}: {0xC5, 0xFD, 0x6F},       // VMOVDQA ymm, ymm
	{OpRegZMM, OpRegZMM}: {0x62, 0xF1, 0x7D, 0x48}, // VMOVDQA32 zmm, zmm

	// ======================
	// 其他特殊传输指令
	// ======================
	// 标志寄存器
	{OpFlags, OpReg}: {0x9C}, // PUSHF
	{OpReg, OpFlags}: {0x9D}, // POPF

	// 标签/偏移
	{OpReg, OpLabel}:  {0xE8}, // CALL rel32
	{OpReg, OpOffset}: {0x8D}, // LEA reg, mem
}

// ======================
// PUSH 指令映射表 (优化版)
// ======================
var pushOpMap = types.OpcodeMap{
	// 通用寄存器 (操作码 50+r, 实际编码需添加寄存器索引)
	{OpReg8}:  {0x50}, // PUSH r8
	{OpReg16}: {0x50}, // PUSH r16
	{OpReg32}: {0x50}, // PUSH r32
	{OpReg64}: {0x50}, // PUSH r64

	// 段寄存器
	{OpSegES}: {0x06},       // PUSH ES
	{OpSegCS}: {0x0E},       // PUSH CS
	{OpSegSS}: {0x16},       // PUSH SS
	{OpSegDS}: {0x1E},       // PUSH DS
	{OpSegFS}: {0x0F, 0xA0}, // PUSH FS
	{OpSegGS}: {0x0F, 0xA8}, // PUSH GS

	// 立即数
	{OpImm8}:  {0x6A}, // PUSH imm8
	{OpImm16}: {0x68}, // PUSH imm16
	{OpImm32}: {0x68}, // PUSH imm32

	// 特殊指令
	{OpFlags}: {0x9C}, // PUSHF/D/Q
	{OpReg}:   {0x60}, // PUSHA/D
}

// ======================
// POP 指令映射表 (优化版)
// ======================
var popOpMap = types.OpcodeMap{
	// 通用寄存器 (操作码 58+r, 实际编码需添加寄存器索引)
	{OpReg8}:  {0x58}, // POP r8
	{OpReg16}: {0x58}, // POP r16
	{OpReg32}: {0x58}, // POP r32
	{OpReg64}: {0x58}, // POP r64

	// 段寄存器
	{OpSegDS}: {0x1F},       // POP DS
	{OpSegES}: {0x07},       // POP ES
	{OpSegSS}: {0x17},       // POP SS
	{OpSegFS}: {0x0F, 0xA1}, // POP FS
	{OpSegGS}: {0x0F, 0xA9}, // POP GS

	// 特殊指令
	{OpFlags}: {0x9D}, // POPF/D/Q
	{OpReg}:   {0x61}, // POPA/D
}

// ======================
// ADD 指令操作码映射 (优化版)
// ======================
var addOpMap = types.OpcodeMap{
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x00}, // ADD reg8, reg8
	{OpReg16, OpReg16}: {0x01}, // ADD reg16, reg16
	{OpReg32, OpReg32}: {0x01}, // ADD reg32, reg32
	{OpReg64, OpReg64}: {0x01}, // ADD reg64, reg64

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x02}, // ADD reg8, mem8
	{OpReg16, OpMem16}: {0x03}, // ADD reg16, mem16
	{OpReg32, OpMem32}: {0x03}, // ADD reg32, mem32
	{OpReg64, OpMem64}: {0x03}, // ADD reg64, mem64

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x00}, // ADD mem8, reg8
	{OpMem16, OpReg16}: {0x01}, // ADD mem16, reg16
	{OpMem32, OpReg32}: {0x01}, // ADD mem32, reg32
	{OpMem64, OpReg64}: {0x01}, // ADD mem64, reg64

	// 立即数到寄存器
	{OpReg8, OpImm8}:   {0x80}, // ADD reg8, imm8 (80 /0)
	{OpReg16, OpImm16}: {0x81}, // ADD reg16, imm16 (81 /0)
	{OpReg32, OpImm32}: {0x81}, // ADD reg32, imm32 (81 /0)
	{OpReg64, OpImm32}: {0x81}, // ADD reg64, imm32 (81 /0)

	// 立即数到内存
	{OpMem8, OpImm8}:   {0x80}, // ADD mem8, imm8 (80 /0)
	{OpMem16, OpImm16}: {0x81}, // ADD mem16, imm16 (81 /0)
	{OpMem32, OpImm32}: {0x81}, // ADD mem32, imm32 (81 /0)
	{OpMem64, OpImm32}: {0x81}, // ADD mem64, imm32 (81 /0)

	// 小立即数到寄存器 (符号扩展)
	{OpReg16, OpImm8}: {0x83}, // ADD reg16, imm8 (83 /0)
	{OpReg32, OpImm8}: {0x83}, // ADD reg32, imm8 (83 /0)
	{OpReg64, OpImm8}: {0x83}, // ADD reg64, imm8 (83 /0)

	// 小立即数到内存 (符号扩展)
	{OpMem16, OpImm8}: {0x83}, // ADD mem16, imm8 (83 /0)
	{OpMem32, OpImm8}: {0x83}, // ADD mem32, imm8 (83 /0)
	{OpMem64, OpImm8}: {0x83}, // ADD mem64, imm8 (83 /0)

	// 累加器特殊编码
	{OpReg8, OpImm8}:   {0x04}, // ADD AL, imm8
	{OpReg16, OpImm16}: {0x05}, // ADD AX, imm16
	{OpReg32, OpImm32}: {0x05}, // ADD EAX, imm32
	{OpReg64, OpImm32}: {0x05}, // ADD RAX, imm32

	// 向量寄存器
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0xFC}, // PADDB xmm, xmm
	{OpRegYMM, OpRegYMM}: {0xC5, 0xFD, 0xFC}, // VPADDB ymm, ymm
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0xFC}, // PADDB xmm, mem
	{OpRegYMM, OpMem256}: {0xC5, 0xFD, 0xFC}, // VPADDB ymm, mem
}

// ======================
// DIV 指令操作码映射 (优化版)
// ======================
var divOpMap = types.OpcodeMap{
	// 8位除法
	{OpReg8, OpNone}: {0xF6, 6}, // DIV reg8
	{OpMem8, OpNone}: {0xF6, 6}, // DIV mem8

	// 16位除法
	{OpReg16, OpNone}: {0xF7, 6}, // DIV reg16
	{OpMem16, OpNone}: {0xF7, 6}, // DIV mem16

	// 32位除法
	{OpReg32, OpNone}: {0xF7, 6}, // DIV reg32
	{OpMem32, OpNone}: {0xF7, 6}, // DIV mem32

	// 64位除法
	{OpReg64, OpNone}: {0xF7, 6}, // DIV reg64
	{OpMem64, OpNone}: {0xF7, 6}, // DIV mem64

	// 向量除法
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0x5E}, // DIVPS xmm, xmm
	{OpRegYMM, OpRegYMM}: {0xC5, 0xFC, 0x5E}, // VDIVPS ymm, ymm
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0x5E}, // DIVPS xmm, mem
	{OpRegYMM, OpMem256}: {0xC5, 0xFC, 0x5E}, // VDIVPS ymm, mem

	// 浮点除法
	{OpFPU, OpFPU}: {0xDC, 0xF8}, // FDIV ST(0), ST(i)
	{OpFPU, OpMem}: {0xDC, 0x30}, // FDIV mem32
}

// ======================
// SUB 指令操作码映射 (优化版)
// ======================
var subOpMap = types.OpcodeMap{
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x28}, // SUB reg8, reg8
	{OpReg16, OpReg16}: {0x29}, // SUB reg16, reg16
	{OpReg32, OpReg32}: {0x29}, // SUB reg32, reg32
	{OpReg64, OpReg64}: {0x29}, // SUB reg64, reg64

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x2A}, // SUB reg8, mem8
	{OpReg16, OpMem16}: {0x2B}, // SUB reg16, mem16
	{OpReg32, OpMem32}: {0x2B}, // SUB reg32, mem32
	{OpReg64, OpMem64}: {0x2B}, // SUB reg64, mem64

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x28}, // SUB mem8, reg8
	{OpMem16, OpReg16}: {0x29}, // SUB mem16, reg16
	{OpMem32, OpReg32}: {0x29}, // SUB mem32, reg32
	{OpMem64, OpReg64}: {0x29}, // SUB mem64, reg64

	// 立即数到寄存器
	{OpReg8, OpImm8}:   {0x80, 5}, // SUB reg8, imm8 (80 /5)
	{OpReg16, OpImm16}: {0x81, 5}, // SUB reg16, imm16 (81 /5)
	{OpReg32, OpImm32}: {0x81, 5}, // SUB reg32, imm32 (81 /5)
	{OpReg64, OpImm32}: {0x81, 5}, // SUB reg64, imm32 (81 /5)

	// 立即数到内存
	{OpMem8, OpImm8}:   {0x80, 5}, // SUB mem8, imm8 (80 /5)
	{OpMem16, OpImm16}: {0x81, 5}, // SUB mem16, imm16 (81 /5)
	{OpMem32, OpImm32}: {0x81, 5}, // SUB mem32, imm32 (81 /5)
	{OpMem64, OpImm32}: {0x81, 5}, // SUB mem64, imm32 (81 /5)

	// 小立即数到寄存器 (符号扩展)
	{OpReg16, OpImm8}: {0x83, 5}, // SUB reg16, imm8 (83 /5)
	{OpReg32, OpImm8}: {0x83, 5}, // SUB reg32, imm8 (83 /5)
	{OpReg64, OpImm8}: {0x83, 5}, // SUB reg64, imm8 (83 /5)

	// 小立即数到内存 (符号扩展)
	{OpMem16, OpImm8}: {0x83, 5}, // SUB mem16, imm8 (83 /5)
	{OpMem32, OpImm8}: {0x83, 5}, // SUB mem32, imm8 (83 /5)
	{OpMem64, OpImm8}: {0x83, 5}, // SUB mem64, imm8 (83 /5)

	// 累加器特殊编码
	{OpReg8, OpImm8}:   {0x2C}, // SUB AL, imm8
	{OpReg16, OpImm16}: {0x2D}, // SUB AX, imm16
	{OpReg32, OpImm32}: {0x2D}, // SUB EAX, imm32
	{OpReg64, OpImm32}: {0x2D}, // SUB RAX, imm32

	// 向量寄存器
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0xF8}, // PSUBB xmm, xmm
	{OpRegYMM, OpRegYMM}: {0xC5, 0xFD, 0xF8}, // VPSUBB ymm, ymm
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0xF8}, // PSUBB xmm, mem
	{OpRegYMM, OpMem256}: {0xC5, 0xFD, 0xF8}, // VPSUBB ymm, mem
}

// ======================
// MUL 指令操作码映射 (优化版)
// ======================
var mulOpMap = types.OpcodeMap{
	// 整数乘法 (一个操作数)
	{OpReg8, OpNone}:  {0xF6, 4}, // MUL reg8
	{OpMem8, OpNone}:  {0xF6, 4}, // MUL mem8
	{OpReg16, OpNone}: {0xF7, 4}, // MUL reg16
	{OpMem16, OpNone}: {0xF7, 4}, // MUL mem16
	{OpReg32, OpNone}: {0xF7, 4}, // MUL reg32
	{OpMem32, OpNone}: {0xF7, 4}, // MUL mem32
	{OpReg64, OpNone}: {0xF7, 4}, // MUL reg64
	{OpMem64, OpNone}: {0xF7, 4}, // MUL mem64

	// 向量乘法
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0xF5}, // PMULHW xmm, xmm
	{OpRegYMM, OpRegYMM}: {0xC5, 0xFD, 0xF5}, // VPMULHW ymm, ymm
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0xF5}, // PMULHW xmm, mem
	{OpRegYMM, OpMem256}: {0xC5, 0xFD, 0xF5}, // VPMULHW ymm, mem

	// 浮点乘法
	{OpFPU, OpFPU}: {0xDC, 0xC8}, // FMUL ST(0), ST(i)
	{OpFPU, OpMem}: {0xDC, 0x0D}, // FMUL mem64
}

// ======================
// RET 指令操作码映射
// ======================
var retOpMap = types.OpcodeMap{
	{OpNone}:  {0xC3, 0, 0, 0}, // RET
	{OpImm16}: {0xC2, 0, 0, 0}, // RET imm16
}

var cmpOpMap = types.OpcodeMap{
	// ======================
	// 通用整数比较指令
	// ======================
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x38, 0, 0, 0}, // CMP reg8, reg8
	{OpReg16, OpReg16}: {0x39, 0, 0, 0}, // CMP reg16, reg16
	{OpReg32, OpReg32}: {0x39, 0, 0, 0}, // CMP reg32, reg32
	{OpReg64, OpReg64}: {0x39, 0, 0, 0}, // CMP reg64, reg64

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x3A, 0, 0, 0}, // CMP reg8, mem8
	{OpReg16, OpMem16}: {0x3B, 0, 0, 0}, // CMP reg16, mem16
	{OpReg32, OpMem32}: {0x3B, 0, 0, 0}, // CMP reg32, mem32
	{OpReg64, OpMem64}: {0x3B, 0, 0, 0}, // CMP reg64, mem64

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x38, 0, 0, 0}, // CMP mem8, reg8
	{OpMem16, OpReg16}: {0x39, 0, 0, 0}, // CMP mem16, reg16
	{OpMem32, OpReg32}: {0x39, 0, 0, 0}, // CMP mem32, reg32
	{OpMem64, OpReg64}: {0x39, 0, 0, 0}, // CMP mem64, reg64

	// 累加器特殊编码
	{OpReg8, OpImm8}:   {0x3C, 0, 0, 0}, // CMP AL, imm8
	{OpReg16, OpImm16}: {0x3D, 0, 0, 0}, // CMP AX, imm16
	{OpReg32, OpImm32}: {0x3D, 0, 0, 0}, // CMP EAX, imm32
	{OpReg64, OpImm32}: {0x3D, 0, 0, 0}, // CMP RAX, imm32

	// 通用立即数到寄存器
	{OpReg8, OpImm8}:   {0x80, 7, 0, 0}, // CMP reg8, imm8 (80 /7)
	{OpReg16, OpImm16}: {0x81, 7, 0, 0}, // CMP reg16, imm16 (81 /7)
	{OpReg32, OpImm32}: {0x81, 7, 0, 0}, // CMP reg32, imm32 (81 /7)
	{OpReg64, OpImm32}: {0x81, 7, 0, 0}, // CMP reg64, imm32 (81 /7)

	// 小立即数到寄存器 (符号扩展)
	{OpReg16, OpImm8}: {0x83, 7, 0, 0}, // CMP reg16, imm8 (83 /7)
	{OpReg32, OpImm8}: {0x83, 7, 0, 0}, // CMP reg32, imm8 (83 /7)
	{OpReg64, OpImm8}: {0x83, 7, 0, 0}, // CMP reg64, imm8 (83 /7)

	// 立即数到内存
	{OpMem8, OpImm8}:   {0x80, 7, 0, 0}, // CMP mem8, imm8 (80 /7)
	{OpMem16, OpImm16}: {0x81, 7, 0, 0}, // CMP mem16, imm16 (81 /7)
	{OpMem32, OpImm32}: {0x81, 7, 0, 0}, // CMP mem32, imm32 (81 /7)
	{OpMem64, OpImm32}: {0x81, 7, 0, 0}, // CMP mem64, imm32 (81 /7)

	// 小立即数到内存 (符号扩展)
	{OpMem16, OpImm8}: {0x83, 7, 0, 0}, // CMP mem16, imm8 (83 /7)
	{OpMem32, OpImm8}: {0x83, 7, 0, 0}, // CMP mem32, imm8 (83 /7)
	{OpMem64, OpImm8}: {0x83, 7, 0, 0}, // CMP mem64, imm8 (83 /7)

	// ======================
	// 字符串比较指令
	// ======================
	{OpMem8, OpReg8}:   {0xA6, 0, 0, 0}, // CMPSB
	{OpMem16, OpReg16}: {0xA7, 0, 0, 0}, // CMPSW
	{OpMem32, OpReg32}: {0xA7, 0, 0, 0}, // CMPSD
	{OpMem64, OpReg64}: {0xA7, 0, 0, 0}, // CMPSQ

	// ======================
	// 系统寄存器比较
	// ======================
	// 控制寄存器
	{OpCR0, OpReg}: {0x0F, 0x20, 0, 0}, // CMP CR0, reg
	{OpCR2, OpReg}: {0x0F, 0x20, 2, 0}, // CMP CR2, reg
	{OpCR3, OpReg}: {0x0F, 0x20, 3, 0}, // CMP CR3, reg
	{OpCR4, OpReg}: {0x0F, 0x20, 4, 0}, // CMP CR4, reg
	{OpCR8, OpReg}: {0x0F, 0x20, 8, 0}, // CMP CR8, reg

	{OpReg, OpCR0}: {0x0F, 0x22, 0, 0}, // CMP reg, CR0
	{OpReg, OpCR2}: {0x0F, 0x22, 2, 0}, // CMP reg, CR2
	{OpReg, OpCR3}: {0x0F, 0x22, 3, 0}, // CMP reg, CR3
	{OpReg, OpCR4}: {0x0F, 0x22, 4, 0}, // CMP reg, CR4
	{OpReg, OpCR8}: {0x0F, 0x22, 8, 0}, // CMP reg, CR8

	// 调试寄存器
	{OpDR0, OpReg}: {0x0F, 0x21, 0, 0}, // CMP DR0, reg
	{OpDR1, OpReg}: {0x0F, 0x21, 1, 0}, // CMP DR1, reg
	{OpDR2, OpReg}: {0x0F, 0x21, 2, 0}, // CMP DR2, reg
	{OpDR3, OpReg}: {0x0F, 0x21, 3, 0}, // CMP DR3, reg
	{OpDR6, OpReg}: {0x0F, 0x21, 6, 0}, // CMP DR6, reg
	{OpDR7, OpReg}: {0x0F, 0x21, 7, 0}, // CMP DR7, reg

	{OpReg, OpDR0}: {0x0F, 0x23, 0, 0}, // CMP reg, DR0
	{OpReg, OpDR1}: {0x0F, 0x23, 1, 0}, // CMP reg, DR1
	{OpReg, OpDR2}: {0x0F, 0x23, 2, 0}, // CMP reg, DR2
	{OpReg, OpDR3}: {0x0F, 0x23, 3, 0}, // CMP reg, DR3
	{OpReg, OpDR6}: {0x0F, 0x23, 6, 0}, // CMP reg, DR6
	{OpReg, OpDR7}: {0x0F, 0x23, 7, 0}, // CMP reg, DR7

	// 测试寄存器
	{OpTR6, OpMem}: {0x0F, 0x24, 6, 0}, // CMP TR6, mem
	{OpTR7, OpMem}: {0x0F, 0x24, 7, 0}, // CMP TR7, mem
	{OpMem, OpTR6}: {0x0F, 0x26, 6, 0}, // CMP mem, TR6
	{OpMem, OpTR7}: {0x0F, 0x26, 7, 0}, // CMP mem, TR7
}

// AND 指令操作码映射
var andOpMap = types.OpcodeMap{
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x20, 0, 0, 0},
	{OpReg16, OpReg16}: {0x21, 0, 0, 0},
	{OpReg32, OpReg32}: {0x21, 0, 0, 0},
	{OpReg64, OpReg64}: {0x21, 0, 0, 0},

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x22, 0, 0, 0},
	{OpReg16, OpMem16}: {0x23, 0, 0, 0},
	{OpReg32, OpMem32}: {0x23, 0, 0, 0},
	{OpReg64, OpMem64}: {0x23, 0, 0, 0},

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x20, 0, 0, 0},
	{OpMem16, OpReg16}: {0x21, 0, 0, 0},
	{OpMem32, OpReg32}: {0x21, 0, 0, 0},
	{OpMem64, OpReg64}: {0x21, 0, 0, 0},

	// 累加器立即数
	{OpReg8, OpImm8}:   {0x24, 0, 0, 0},
	{OpReg16, OpImm16}: {0x25, 0, 0, 0},
	{OpReg32, OpImm32}: {0x25, 0, 0, 0},
	{OpReg64, OpImm32}: {0x25, 0, 0, 0},

	// 通用立即数
	{OpReg8, OpImm8}:   {0x80, 4, 0, 0},
	{OpReg16, OpImm16}: {0x81, 4, 0, 0},
	{OpReg32, OpImm32}: {0x81, 4, 0, 0},
	{OpReg64, OpImm32}: {0x81, 4, 0, 0},

	// 小立即数
	{OpReg16, OpImm8}: {0x83, 4, 0, 0},
	{OpReg32, OpImm8}: {0x83, 4, 0, 0},
	{OpReg64, OpImm8}: {0x83, 4, 0, 0},

	// 内存立即数
	{OpMem8, OpImm8}:   {0x80, 4, 0, 0},
	{OpMem16, OpImm16}: {0x81, 4, 0, 0},
	{OpMem32, OpImm32}: {0x81, 4, 0, 0},
	{OpMem64, OpImm32}: {0x81, 4, 0, 0},

	// 内存小立即数
	{OpMem16, OpImm8}: {0x83, 4, 0, 0},
	{OpMem32, OpImm8}: {0x83, 4, 0, 0},
	{OpMem64, OpImm8}: {0x83, 4, 0, 0},

	// 向量指令
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0x54, 0},
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0x54, 0},
	{OpRegXMM, OpImm8}:   {0x66, 0x0F, 0x54, 0},

	{OpRegYMM, OpRegYMM}: {0xC5, 0xFC, 0x54, 0},
	{OpRegYMM, OpMem256}: {0xC5, 0xFC, 0x54, 0},
	{OpRegYMM, OpImm8}:   {0xC5, 0xFC, 0x54, 0},

	{OpRegZMM, OpRegZMM}: {0x62, 0xF1, 0x7C, 0x48},
	{OpRegZMM, OpMem512}: {0x62, 0xF1, 0x7C, 0x48},
	{OpRegZMM, OpImm8}:   {0x62, 0xF1, 0x7C, 0x48},
}

// OR 指令操作码映射
var orOpMap = types.OpcodeMap{
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x08, 0, 0, 0},
	{OpReg16, OpReg16}: {0x09, 0, 0, 0},
	{OpReg32, OpReg32}: {0x09, 0, 0, 0},
	{OpReg64, OpReg64}: {0x09, 0, 0, 0},

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x0A, 0, 0, 0},
	{OpReg16, OpMem16}: {0x0B, 0, 0, 0},
	{OpReg32, OpMem32}: {0x0B, 0, 0, 0},
	{OpReg64, OpMem64}: {0x0B, 0, 0, 0},

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x08, 0, 0, 0},
	{OpMem16, OpReg16}: {0x09, 0, 0, 0},
	{OpMem32, OpReg32}: {0x09, 0, 0, 0},
	{OpMem64, OpReg64}: {0x09, 0, 0, 0},

	// 累加器立即数
	{OpReg8, OpImm8}:   {0x0C, 0, 0, 0},
	{OpReg16, OpImm16}: {0x0D, 0, 0, 0},
	{OpReg32, OpImm32}: {0x0D, 0, 0, 0},
	{OpReg64, OpImm32}: {0x0D, 0, 0, 0},

	// 通用立即数
	{OpReg8, OpImm8}:   {0x80, 1, 0, 0},
	{OpReg16, OpImm16}: {0x81, 1, 0, 0},
	{OpReg32, OpImm32}: {0x81, 1, 0, 0},
	{OpReg64, OpImm32}: {0x81, 1, 0, 0},

	// 小立即数
	{OpReg16, OpImm8}: {0x83, 1, 0, 0},
	{OpReg32, OpImm8}: {0x83, 1, 0, 0},
	{OpReg64, OpImm8}: {0x83, 1, 0, 0},

	// 内存立即数
	{OpMem8, OpImm8}:   {0x80, 1, 0, 0},
	{OpMem16, OpImm16}: {0x81, 1, 0, 0},
	{OpMem32, OpImm32}: {0x81, 1, 0, 0},
	{OpMem64, OpImm32}: {0x81, 1, 0, 0},

	// 内存小立即数
	{OpMem16, OpImm8}: {0x83, 1, 0, 0},
	{OpMem32, OpImm8}: {0x83, 1, 0, 0},
	{OpMem64, OpImm8}: {0x83, 1, 0, 0},

	// 向量指令
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0x56, 0},
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0x56, 0},
	{OpRegXMM, OpImm8}:   {0x66, 0x0F, 0x56, 0},

	{OpRegYMM, OpRegYMM}: {0xC5, 0xFC, 0x56, 0},
	{OpRegYMM, OpMem256}: {0xC5, 0xFC, 0x56, 0},
	{OpRegYMM, OpImm8}:   {0xC5, 0xFC, 0x56, 0},

	{OpRegZMM, OpRegZMM}: {0x62, 0xF1, 0x7C, 0x48},
	{OpRegZMM, OpMem512}: {0x62, 0xF1, 0x7C, 0x48},
	{OpRegZMM, OpImm8}:   {0x62, 0xF1, 0x7C, 0x48},
}

// XOR 指令操作码映射
var xorOpMap = types.OpcodeMap{
	// 寄存器到寄存器
	{OpReg8, OpReg8}:   {0x30, 0, 0, 0},
	{OpReg16, OpReg16}: {0x31, 0, 0, 0},
	{OpReg32, OpReg32}: {0x31, 0, 0, 0},
	{OpReg64, OpReg64}: {0x31, 0, 0, 0},

	// 内存到寄存器
	{OpReg8, OpMem8}:   {0x32, 0, 0, 0},
	{OpReg16, OpMem16}: {0x33, 0, 0, 0},
	{OpReg32, OpMem32}: {0x33, 0, 0, 0},
	{OpReg64, OpMem64}: {0x33, 0, 0, 0},

	// 寄存器到内存
	{OpMem8, OpReg8}:   {0x30, 0, 0, 0},
	{OpMem16, OpReg16}: {0x31, 0, 0, 0},
	{OpMem32, OpReg32}: {0x31, 0, 0, 0},
	{OpMem64, OpReg64}: {0x31, 0, 0, 0},

	// 累加器立即数
	{OpReg8, OpImm8}:   {0x34, 0, 0, 0},
	{OpReg16, OpImm16}: {0x35, 0, 0, 0},
	{OpReg32, OpImm32}: {0x35, 0, 0, 0},
	{OpReg64, OpImm32}: {0x35, 0, 0, 0},

	// 通用立即数
	{OpReg8, OpImm8}:   {0x80, 6, 0, 0},
	{OpReg16, OpImm16}: {0x81, 6, 0, 0},
	{OpReg32, OpImm32}: {0x81, 6, 0, 0},
	{OpReg64, OpImm32}: {0x81, 6, 0, 0},

	// 小立即数
	{OpReg16, OpImm8}: {0x83, 6, 0, 0},
	{OpReg32, OpImm8}: {0x83, 6, 0, 0},
	{OpReg64, OpImm8}: {0x83, 6, 0, 0},

	// 内存立即数
	{OpMem8, OpImm8}:   {0x80, 6, 0, 0},
	{OpMem16, OpImm16}: {0x81, 6, 0, 0},
	{OpMem32, OpImm32}: {0x81, 6, 0, 0},
	{OpMem64, OpImm32}: {0x81, 6, 0, 0},

	// 内存小立即数
	{OpMem16, OpImm8}: {0x83, 6, 0, 0},
	{OpMem32, OpImm8}: {0x83, 6, 0, 0},
	{OpMem64, OpImm8}: {0x83, 6, 0, 0},

	// 向量指令
	{OpRegXMM, OpRegXMM}: {0x66, 0x0F, 0x57, 0},
	{OpRegXMM, OpMem128}: {0x66, 0x0F, 0x57, 0},
	{OpRegXMM, OpImm8}:   {0x66, 0x0F, 0x57, 0},

	{OpRegYMM, OpRegYMM}: {0xC5, 0xFC, 0x57, 0},
	{OpRegYMM, OpMem256}: {0xC5, 0xFC, 0x57, 0},
	{OpRegYMM, OpImm8}:   {0xC5, 0xFC, 0x57, 0},

	{OpRegZMM, OpRegZMM}: {0x62, 0xF1, 0x7C, 0x48},
	{OpRegZMM, OpMem512}: {0x62, 0xF1, 0x7C, 0x48},
	{OpRegZMM, OpImm8}:   {0x62, 0xF1, 0x7C, 0x48},
}

// CALL 指令操作码映射
var callOpMap = types.OpcodeMap{
	// ======================
	// 直接调用
	// ======================
	// 标签调用
	{OpLabel}: {0xE8}, // CALL label

	// 相对地址调用
	{OpRel}: {0xE8}, // CALL rel32

	// 绝对地址调用
	{OpImm}: {0x9A}, // CALL far absolute

	// ======================
	// 间接调用
	// ======================
	// 寄存器间接调用
	{OpReg}: {0xFF, 0x10}, // CALL reg (FF /2)

	// 内存间接调用
	{OpMem}: {0xFF, 0x10}, // CALL mem (FF /2)
}

func DoASM(i *parser.Instruction, arch *types.Architecture) types.OpBytes {
	builtin := NewX86Builtin(arch)
	opcode := types.OpBytes{}
	if i.IsBuiltin() {
		switch i.Instruction {
		case "ADD":
			// 处理ADD指令：根据操作数类型生成不同机器码
			return builtin.Add(i)
		case "MOV":
			// 处理MOV指令：根据操作数类型生成不同机器码
			return builtin.Mov(i)
		case "PUSH":
			// 处理PUSH指令：将操作数压入栈
			return builtin.Push(i)
		case "POP":
			// 处理POP指令：从栈中弹出数据
			return builtin.Pop(i)
		case "SUB":
			// 处理SUB指令：减法运算
			return builtin.Sub(i)
		case "DIV":
			// 处理ADD指令：加法运算
			return builtin.Div(i)
		case "MUL":
			// 处理MUL指令：乘法运算
			return builtin.Mul(i)
		case "CMP":
			// 处理CMP指令：比较操作数
			return builtin.Cmp(i)
		case "CALL":
			// 处理CALL指令：调用子程序
			return builtin.Call(i)
		case "RET":
			// 处理RET指令：返回调用者
			return builtin.Ret(i)
		case "JMP":
			// 处理JMP指令：无条件跳转
			return builtin.Jmp(i)
		case "JE", "JZ":
			// 处理JE/JZ指令：相等/零跳转
			return builtin.JmpZero(i)
		case "XOR":
			// 处理XOR指令：异或运算
			return builtin.Xor(i)
		case "AND":
			// 处理AND指令：逻辑与运算
			return builtin.And(i)
		case "NOT":
			// 处理NOT指令：按位取反
			return builtin.Not(i)
		case "OR":
			// 处理OR指令：逻辑或运算
			return builtin.Or(i)
		default:
			// 未知指令返回空代码
			return []byte{}
		}
	}
	// 查询指令映射表
	op := arch.Instructions[i.Instruction]
	aop := types.Operands{}
	for i, arg := range i.Args {
		aop[i].Set(ValueToOperand(arg))
	}
	ok := false
	for i := 0; i < len(op.Operands); i++ {
		if op.Operands[i].Has(aop) || op.Operands[i].Is(aop) {
			ok = true
			break
		}
	}
	if !ok {
		panic("invalid operand")
	}
	for e := 0; e < len(aop); e++ {
		tmp := getOperandsSize(aop[e])
		if tmp > i.OpSize {
			i.OpSize = tmp
		}
	}
	opdBytes, err := NewOperandsEncoder(i).EncodeOperands()
	if err != nil {
		panic(err)
	}
	opcode = append(opcode, op.Opcode...)
	opcode = append(opcode, opdBytes...)
	return opcode
}

func opMapHandler(i *parser.Instruction, tab types.OpcodeMap, tryA ...bool) types.OpBytes {
	try := false
	if len(tryA) > 0 {
		try = tryA[0]
	}
	aop := types.Operands{}
	for e := 0; e < len(i.Args); e++ {
		aop[e] = ValueToOperand(i.Args[e])
	}
	opb := types.OpBytes{}
	for operand, op := range tab {
		if operand.Has(aop) || operand.Is(aop) {
			aop = operand
			opb = op
			break
		}
	}
	if len(opb) == 0 {
		if !try {
			panic("")
		}
		return nil
	}
	for e := 0; e < len(aop); e++ {
		tmp := getOperandsSize(aop[e])
		if tmp > i.OpSize {
			i.OpSize = tmp
		}
	}
	opdBytes, err := NewOperandsEncoder(i).EncodeOperands()
	if err != nil {
		if !try {
			panic(err)
		}
		return nil
	}
	opb = append(opb, opdBytes...)
	return opb
}
