package x86

import (
	"CuteASM/arch/types"
	"CuteASM/parser"
	"encoding/binary"
	"fmt"
	"strings"
)

// 操作数类型位标志 (使用int64)
const (
	// ======================
	// 基础寄存器类型
	// ======================
	OpReg8   types.Operand = 1 << iota // 8位通用寄存器 (AL, CL, ...)
	OpReg16                            // 16位通用寄存器 (AX, CX, ...)
	OpReg32                            // 32位通用寄存器 (EAX, ECX, ...)
	OpReg64                            // 64位通用寄存器 (RAX, RCX, ...)
	OpRegXMM                           // XMM寄存器 (XMM0-XMM15)
	OpRegYMM                           // YMM寄存器 (YMM0-YMM15)
	OpRegZMM                           // ZMM寄存器 (ZMM0-ZMM31)
	OpRegK                             // 掩码寄存器 (K0-K7)

	// ======================
	// 系统寄存器类型
	// ======================
	OpSegES // ES段寄存器
	OpSegCS // CS段寄存器
	OpSegSS // SS段寄存器
	OpSegDS // DS段寄存器
	OpSegFS // FS段寄存器
	OpSegGS // GS段寄存器

	OpCR0 // 控制寄存器0
	OpCR2 // 控制寄存器2
	OpCR3 // 控制寄存器3
	OpCR4 // 控制寄存器4
	OpCR8 // 控制寄存器8 (仅64位)

	OpDR0 // 调试寄存器0
	OpDR1 // 调试寄存器1
	OpDR2 // 调试寄存器2
	OpDR3 // 调试寄存器3
	OpDR6 // 调试寄存器6
	OpDR7 // 调试寄存器7

	OpTR6 // 测试寄存器6
	OpTR7 // 测试寄存器7

	// ======================
	// 立即数类型
	// ======================
	OpImm8   // 8位立即数
	OpImm16  // 16位立即数
	OpImm32  // 32位立即数
	OpImm64  // 64位立即数
	OpImmF32 // 32位浮点立即数
	OpImmF64 // 64位浮点立即数

	// ======================
	// 内存类型
	// ======================
	OpMem8   // 8位内存操作数
	OpMem16  // 16位内存操作数
	OpMem32  // 32位内存操作数
	OpMem64  // 64位内存操作数
	OpMem128 // 128位内存操作数
	OpMem256 // 256位内存操作数
	OpMem512 // 512位内存操作数

	// ======================
	// 特殊类型
	// ======================
	OpRel8   // 8位相对跳转
	OpRel16  // 16位相对跳转
	OpRel32  // 32位相对跳转
	OpRel64  // 64位相对跳转
	OpLabel  // 标签引用
	OpOffset // 偏移量

	// ======================
	// 特殊功能类型
	// ======================
	OpFPU // FPU栈寄存器 (ST0-ST7)
	OpMMX // MMX寄存器 (MM0-MM7)
	OpBND // 边界寄存器 (BND0-BND3)
	OpTMM // Tile矩阵寄存器 (TMM0-TMM7)

	// 标志寄存器
	OpFlags

	// ======================
	// 组合类型
	// ======================
	// 寄存器组合
	OpReg types.Operand = OpReg8 | OpReg16 | OpReg32 | OpReg64 |
		OpRegXMM | OpRegYMM | OpRegZMM | OpRegK |
		OpFPU | OpMMX | OpBND | OpTMM | OpFlags

	// 系统寄存器组合
	OpSysReg types.Operand = OpSegES | OpSegCS | OpSegSS | OpSegDS | OpSegFS | OpSegGS |
		OpCR0 | OpCR2 | OpCR3 | OpCR4 | OpCR8 |
		OpDR0 | OpDR1 | OpDR2 | OpDR3 | OpDR6 | OpDR7 |
		OpTR6 | OpTR7

	OpDR  types.Operand = OpDR0 | OpDR1 | OpDR2 | OpDR3 | OpDR6 | OpDR7
	OpTR  types.Operand = OpTR6 | OpTR7
	OpCR  types.Operand = OpCR0 | OpCR2 | OpCR3 | OpCR4 | OpCR8
	OpSeg types.Operand = OpSegES | OpSegCS | OpSegSS | OpSegDS | OpSegFS | OpSegGS

	// 立即数组合
	OpImm types.Operand = OpImm8 | OpImm16 | OpImm32 | OpImm64 | OpImmF32 | OpImmF64

	// 内存组合
	OpMem types.Operand = OpMem8 | OpMem16 | OpMem32 | OpMem64 |
		OpMem128 | OpMem256 | OpMem512

	// 相对跳转组合
	OpRel types.Operand = OpRel8 | OpRel16 | OpRel32 | OpRel64

	// 所有类型
	OpAll types.Operand = OpMem | OpReg | OpSysReg | OpImm | OpRel | OpLabel | OpOffset

	// 空类型
	OpNone types.Operand = 0
	TEST                 = OpImm8 | OpImm16 | OpImm32 | OpImm64
)

func getOperandsSize(opd types.Operand) int {
	lengthMap := map[types.Operand]int{
		OpImm8 | OpMem8 | OpReg8 | OpRel8:     8,
		OpImm16 | OpMem16 | OpReg16 | OpRel16: 16,
		OpImm32 | OpMem32 | OpReg32 | OpRel32: 32,
		OpImm64 | OpMem64 | OpReg64 | OpRel64: 64,
	}
	for k, v := range lengthMap {
		if opd.Has(k) {
			return v
		}
	}
	return 0
}

// 操作数类型常量
const (
	REGISTER = iota
	MEMORY
	IMMEDIATE
	LABEL
	PSEUDO
)

// 操作数编码器
type OperandsEncoder struct {
	inst     *parser.Instruction
	argTypes []types.Operand    // 缓存操作数类型
	memAddr  *parser.MemoryAddr // 缓存内存操作数地址
}

// 创建新的操作数编码器
func NewOperandsEncoder(inst *parser.Instruction) *OperandsEncoder {
	encoder := &OperandsEncoder{
		inst: inst,
	}

	// 预计算操作数类型
	encoder.argTypes = make([]types.Operand, len(inst.Args))
	for i, arg := range inst.Args {
		encoder.argTypes[i] = ValueToOperand(arg)
	}

	// 查找内存操作数
	for _, arg := range inst.Args {
		if arg.Type == parser.ADDR {
			encoder.memAddr = arg.Addr
			break
		}
	}

	return encoder
}

// 编码操作数部分（ModR/M、SIB、位移、立即数）
func (e *OperandsEncoder) EncodeOperands() ([]byte, error) {
	var operandBytes []byte

	// 处理ModR/M字节
	modRM, err := e.generateModRM()
	if err != nil {
		return nil, err
	}
	if modRM != 0 {
		operandBytes = append(operandBytes, modRM)
	}

	// 处理SIB字节
	sib, err := e.generateSIB()
	if err != nil {
		return nil, err
	}
	if sib != 0 {
		operandBytes = append(operandBytes, sib)
	}

	// 处理位移
	disp, err := e.generateDisplacement()
	if err != nil {
		return nil, err
	}
	operandBytes = append(operandBytes, disp...)

	// 处理立即数
	imm, err := e.generateImmediate()
	if err != nil {
		return nil, err
	}
	operandBytes = append(operandBytes, imm...)

	return operandBytes, nil
}

// 生成ModR/M字节
func (e *OperandsEncoder) generateModRM() (byte, error) {
	// 单操作数指令处理
	if len(e.inst.Args) == 1 {
		return e.encodeSingleOperand()
	}

	// 双操作数指令处理
	if len(e.inst.Args) >= 2 {
		return e.encodeDualOperands()
	}

	return 0, nil // 无操作数指令
}

// 编码单操作数指令的操作数部分
func (e *OperandsEncoder) encodeSingleOperand() (byte, error) {
	arg := e.inst.Args[0]
	opType := e.argTypes[0]

	switch {
	case opType.Has(OpReg):
		// 寄存器操作数
		return byte(arg.Reg.Num & 0x7), nil

	case opType.Has(OpMem):
		// 内存操作数
		mod := e.getModValue(arg.Addr.Displacement, arg.Addr.LabelRef)
		rm := e.getRMValue(arg.Addr)
		return mod<<6 | rm, nil

	case opType.Has(OpImm):
		return 0, nil // 立即数单独处理

	case opType.Has(OpLabel):
		return 0, nil // 标签单独处理

	default:
		return 0, fmt.Errorf("unsupported single operand for %s: %s",
			e.inst.Instruction, operandTypeName(opType))
	}
}

// 编码双操作数指令的操作数部分
func (e *OperandsEncoder) encodeDualOperands() (byte, error) {
	fmt.Println("^_^")
	dst := e.inst.Args[0]
	src := e.inst.Args[1]
	dstType := e.argTypes[0]
	srcType := e.argTypes[1]

	// 寄存器到寄存器
	if dstType.Has(OpReg) && srcType.Has(OpReg) {
		return (byte(src.Reg.Num&7) << 3) | byte(dst.Reg.Num&7), nil
	}


	// 寄存器到内存
	if dstType.Has(OpReg) && srcType.Has(OpMem) {
		if src.Addr == nil {
			return 0, fmt.Errorf("memory operand is missing address information")
		}

		mod := e.getModValue(src.Addr.Displacement, src.Addr.LabelRef)
		rm := e.getRMValue(src.Addr)
		fmt.Println("寄存器到内存", dst.Reg.Num, mod, rm)
		return (byte(dst.Reg.Num & 7) << 3) | (mod << 6) | rm, nil
	}

	// 立即数到内存
	if dstType.Has(OpMem) && srcType.Has(OpImm) {
		if dst.Addr == nil {
			return 0, fmt.Errorf("memory operand is missing address information")
		}

		mod := e.getModValue(dst.Addr.Displacement, dst.Addr.LabelRef)
		rm := e.getRMValue(dst.Addr)
		return (mod << 6) | rm, nil
	}

	// 立即数到寄存器
	if dstType.Has(OpReg) && srcType.Has(OpImm) {
		// 对于MOV指令，有特殊处理（使用B0-B7操作码）
		if e.inst.Instruction == "MOV" {
			return 0, nil // MOV指令有特殊操作码，不需要ModR/M
		}

		// 其他指令（如ADD、SUB等）需要ModR/M
		return byte(dst.Reg.Num&7) << 3, nil
	}

	// 标签到寄存器
	if dstType.Has(OpReg) && srcType.Has(OpLabel) {
		return byte(dst.Reg.Num&7) << 3, nil
	}

	// 系统寄存器处理
	if dstType.Has(OpSysReg) || srcType.Has(OpSysReg) {
		return e.encodeSysReg()
	}

	return 0, fmt.Errorf("unsupported operand combination for %s: %s, %s",
		e.inst.Instruction, operandTypeName(dstType), operandTypeName(srcType))
}

// 编码系统寄存器操作数部分
func (e *OperandsEncoder) encodeSysReg() (byte, error) {
	dstType := e.argTypes[0]
	srcType := e.argTypes[1]

	// 控制寄存器
	if dstType.Has(OpCR) || srcType.Has(OpCR) {
		return e.encodeControlReg()
	}

	// 调试寄存器
	if dstType.Has(OpDR) || srcType.Has(OpDR) {
		return e.encodeDebugReg()
	}

	// 段寄存器
	if dstType.Has(OpSeg) || srcType.Has(OpSeg) {
		return e.encodeSegmentReg()
	}

	return 0, fmt.Errorf("unsupported system register for %s", e.inst.Instruction)
}

// 编码控制寄存器操作数部分
func (e *OperandsEncoder) encodeControlReg() (byte, error) {
	dstType := e.argTypes[0]

	var regField, crField int

	if dstType.Has(OpCR) {
		// dst是控制寄存器
		crField = getCRField(e.inst.Args[0].Reg.Name)
		regField = e.inst.Args[1].Reg.Num
	} else {
		// src是控制寄存器
		crField = getCRField(e.inst.Args[1].Reg.Name)
		regField = e.inst.Args[0].Reg.Num
	}

	return byte(crField<<3) | byte(regField&7), nil
}

// 编码调试寄存器操作数部分
func (e *OperandsEncoder) encodeDebugReg() (byte, error) {
	dstType := e.argTypes[0]

	var regField, drField int

	if dstType.Has(OpDR) {
		// dst是调试寄存器
		drField = getDRField(e.inst.Args[0].Reg.Name)
		regField = e.inst.Args[1].Reg.Num
	} else {
		// src是调试寄存器
		drField = getDRField(e.inst.Args[1].Reg.Name)
		regField = e.inst.Args[0].Reg.Num
	}

	return byte(drField<<3) | byte(regField&7), nil
}

// 编码段寄存器操作数部分
func (e *OperandsEncoder) encodeSegmentReg() (byte, error) {
	dstType := e.argTypes[0]

	var segField int

	if dstType.Has(OpSeg) {
		// dst是段寄存器
		segField = getSegField(e.inst.Args[0].Reg.Name)
	} else {
		// src是段寄存器
		segField = getSegField(e.inst.Args[1].Reg.Name)
	}

	return byte(segField << 3), nil
}

// 获取Mod字段值
func (e *OperandsEncoder) getModValue(disp float64, labelRef string) byte {
	if labelRef != "" {
		return 0b10 // 标签引用总是32位位移
	}
	if disp == 0 {
		return 0b00 // 无位移
	}
	if disp >= -128 && disp <= 127 {
		return 0b01 // 8位位移
	}
	return 0b10 // 32位位移
}

// 获取RM字段值
func (e *OperandsEncoder) getRMValue(addr *parser.MemoryAddr) byte {
	if addr == nil {
		return 0
	}

	// 如果有索引寄存器，使用SIB寻址
	if addr.IndexReg != nil {
		return 0b100 // SIB寻址
	}

	// 安全获取基址寄存器编号
	baseRegNum := 0
	if addr.BaseReg != nil {
		baseRegNum = addr.BaseReg.Num
	}

	return byte(baseRegNum & 0x7)
}

// 生成SIB字节
func (e *OperandsEncoder) generateSIB() (byte, error) {
	if e.memAddr == nil || e.memAddr.IndexReg == nil {
		return 0, nil
	}

	scaleCode := 0
	switch e.memAddr.Scale {
	case 2:
		scaleCode = 1
	case 4:
		scaleCode = 2
	case 8:
		scaleCode = 3
	}

	indexNum := e.memAddr.IndexReg.Num & 0x7
	baseNum := e.memAddr.BaseReg.Num & 0x7

	return byte((scaleCode << 6) | (indexNum << 3) | baseNum), nil
}

// 生成位移
func (e *OperandsEncoder) generateDisplacement() ([]byte, error) {
	if e.memAddr == nil {
		return nil, nil
	}

	dispValue := e.memAddr.Displacement
	labelRef := e.memAddr.LabelRef

	if dispValue == 0 && labelRef == "" {
		return nil, nil
	}

	if labelRef != "" {
		return make([]byte, 4), nil // 预留32位空间
	}

	if dispValue >= -128 && dispValue <= 127 {
		return []byte{byte(int8(dispValue))}, nil
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(int32(dispValue)))
	return buf, nil
}

// 生成立即数
func (e *OperandsEncoder) generateImmediate() ([]byte, error) {
	// 查找立即数操作数
	var immValue float64
	var immSize int

	for _, arg := range e.inst.Args {
		if arg.Type == parser.NUMBER {
			immValue = arg.Num
			immSize = e.getImmSize()
			break
		}
		if arg.Type == parser.LABEL {
			// 标签作为立即数处理
			return make([]byte, 4), nil // 预留32位空间
		}
	}

	if immValue == 0 {
		return nil, nil
	}

	// 编码立即数
	switch immSize {
	case 8:
		return []byte{byte(int8(immValue))}, nil
	case 16:
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(int16(immValue)))
		return buf, nil
	case 32:
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(int32(immValue)))
		return buf, nil
	case 64:
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(int64(immValue)))
		return buf, nil
	default:
		return nil, fmt.Errorf("unsupported immediate size: %d", immSize)
	}
}

// 获取立即数大小（使用OpSize优先）
func (e *OperandsEncoder) getImmSize() int {
	// 优先使用指令的OpSize字段
	if e.inst.OpSize > 0 {
		return e.inst.OpSize
	}

	// 如果没有设置OpSize，则根据操作数类型确定
	for i, arg := range e.inst.Args {
		if arg.Type == parser.NUMBER {
			// 根据目标操作数类型确定大小
			if i == 0 && e.argTypes[0].Has(OpReg) {
				switch e.argTypes[0] {
				case OpReg8:
					return 8
				case OpReg16:
					return 16
				case OpReg32:
					return 32
				case OpReg64:
					return 64
				}
			}

			// 根据数值大小确定
			return getImmSizeByValue(arg.Num)
		}
	}

	return 32 // 默认32位
}

// 根据值获取立即数大小
func getImmSizeByValue(value float64) int {
	// 根据数值大小确定
	if value >= -128 && value <= 127 {
		return 8
	}
	if value >= -32768 && value <= 32767 {
		return 16
	}
	if value >= -2147483648 && value <= 2147483647 {
		return 32
	}
	return 64
}

// 获取控制寄存器字段
func getCRField(regName string) int {
	switch strings.ToUpper(regName) {
	case "CR0":
		return 0
	case "CR2":
		return 2
	case "CR3":
		return 3
	case "CR4":
		return 4
	case "CR8":
		return 8
	default:
		return 0
	}
}

// 获取调试寄存器字段
func getDRField(regName string) int {
	switch strings.ToUpper(regName) {
	case "DR0":
		return 0
	case "DR1":
		return 1
	case "DR2":
		return 2
	case "DR3":
		return 3
	case "DR6":
		return 6
	case "DR7":
		return 7
	default:
		return 0
	}
}

// 获取段寄存器字段
func getSegField(regName string) int {
	switch strings.ToUpper(regName) {
	case "ES":
		return 0
	case "CS":
		return 1
	case "SS":
		return 2
	case "DS":
		return 3
	case "FS":
		return 4
	case "GS":
		return 5
	default:
		return 0
	}
}

// 操作数类型名称
func operandTypeName(opType types.Operand) string {
	switch {
	case opType.Has(OpReg8):
		return "REG8"
	case opType.Has(OpReg16):
		return "REG16"
	case opType.Has(OpReg32):
		return "REG32"
	case opType.Has(OpReg64):
		return "REG64"
	case opType.Has(OpRegXMM):
		return "XMM"
	case opType.Has(OpRegYMM):
		return "YMM"
	case opType.Has(OpRegZMM):
		return "ZMM"
	case opType.Has(OpRegK):
		return "MASK"
	case opType.Has(OpSysReg):
		return "SYSREG"
	case opType.Has(OpImm8):
		return "IMM8"
	case opType.Has(OpImm16):
		return "IMM16"
	case opType.Has(OpImm32):
		return "IMM32"
	case opType.Has(OpImm64):
		return "IMM64"
	case opType.Has(OpMem8):
		return "MEM8"
	case opType.Has(OpMem16):
		return "MEM16"
	case opType.Has(OpMem32):
		return "MEM32"
	case opType.Has(OpMem64):
		return "MEM64"
	case opType.Has(OpMem128):
		return "MEM128"
	case opType.Has(OpMem256):
		return "MEM256"
	case opType.Has(OpMem512):
		return "MEM512"
	case opType.Has(OpLabel):
		return "LABEL"
	default:
		return "UNKNOWN"
	}
}
