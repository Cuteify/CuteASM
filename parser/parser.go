package parser

import (
	"CuteASM/arch/types"
	errorUtil "CuteASM/error"
	"CuteASM/lexer"

	"strings"
)

type Parser struct {
	Block       *Node // block
	ThisBlock   *Node // 当前block
	Lexer       *lexer.Lexer
	BracketsNum int
	Error       *errorUtil.Error
	Funcs       map[string]*Node
	Vars        map[string]*Node
	Types       map[string]*Node
	DontBack    int
	IsInFunc    bool
	line        int // 添加行号计数器
	arch        *types.Architecture
}

func (p *Parser) Next() (finish bool) {
	//beforeCursor := p.Lexer.Cursor
	code := p.Lexer.Next()
	if code.IsEmpty() {
		finish = true
		return
	}
	//fmt.Println(code)
	switch code.Type {
	case lexer.PSEUDO:
		switch code.Value {
		case "SECTION":
			section := &SECTION{}
			section.Parse(p)
		case "VAR":
			//varBlock := &VarBlock{}
			//varBlock.Parse(p)
		}
	case lexer.NAME:
		oldCursor := p.Lexer.Cursor
		code2 := p.Lexer.Next()
		if code2.IsEmpty() {
			finish = true
			return
		}
		if code2.Type == lexer.SEPARATOR && code2.Value == ":" {
			/*label := &LabelBlock{
				Name: code.Value,
			}
			label.Parse(p)*/
		} else {
			p.Lexer.Cursor = oldCursor
		}
	case lexer.INSTRUCTION:
		//instruction := &Instruction{Instruction: code.Value}
		//instruction.Parse(p)
		//if code.Value == "CALL" {
		//	fmt.Println("CALL", instruction)
		//}
	default:
		if code.Type == lexer.SEPARATOR && code.Value != ";" && code.Value != "\n" && code.Value != "\r" {
			p.Lexer.Error.MissError("Syntax Error", p.Lexer.Cursor, "Miss "+code.Value)
		}
	}
	return
}

func (p *Parser) AddChild(node *Node) {
	p.ThisBlock.AddChild(node)
}

func (p *Parser) Back(num int) error {
	if num == 0 {
		return nil
	}
	if p.ThisBlock.Father == nil {
		p.Error.MissError("Internal Compiler Errors", p.Lexer.Cursor, "Back at root")
	}
	if p.DontBack != 0 {
		p.DontBack--
		return p.Back(num - 1)
	}
	p.ThisBlock = p.ThisBlock.Father
	if num < 0 {
		num = -num
	}
	return p.Back(num - 1)
}

func (p *Parser) Need(value string) []lexer.Token {
	tmp2 := []lexer.Token{}
	for {
		tmp := p.Lexer.Next()
		if tmp.IsEmpty() {
			p.Error.MissError("Syntax Error", p.Lexer.Cursor, "need '"+value+"'")
		}
		if tmp.Value == "\n" {
			p.Error.MissError("Syntax Error", p.Lexer.Cursor, "need '"+value+"'")
		}
		tmp2 = append(tmp2, tmp)
		if tmp.Value == value && tmp.Type != lexer.STRING {
			return tmp2
		}
	}
}

func (p *Parser) FindEndCursor() int {
	tmp := strings.Index(p.Lexer.Text[p.Lexer.Cursor:], p.Lexer.LineFeed)
	if tmp == -1 {
		return len(p.Lexer.Text) - 1
	}
	return tmp + p.Lexer.Cursor
}

func (p *Parser) Wait(value string) int {
	return len(p.Need(value))
}

func (p *Parser) Has(token lexer.Token, stopCursor int) int {
	startCursor := p.Lexer.Cursor
	for stopCursor > p.Lexer.Cursor {
		code := p.Lexer.Next()
		if code.IsEmpty() {
			p.Lexer.Error.MissError("Invalid expression", p.Lexer.Cursor, "Incomplete expression")
		}
		if code.Value == token.Value && code.Type == token.Type {
			cursorTmp := p.Lexer.Cursor
			p.Lexer.Cursor = startCursor
			return cursorTmp
		}
	}
	p.Lexer.Cursor = startCursor
	return -1
}

func NewParser(lexer *lexer.Lexer, Arch *types.Architecture) *Parser {
	p := &Parser{
		Lexer: lexer,
		Error: lexer.Error,
		arch:  Arch,
	}
	p.Block = &Node{}
	p.ThisBlock = p.Block
	return p
}

func (p *Parser) Parse() *Node {
	for {
		if !p.Line() {
			break
		}
	}
	return p.Block
}

func (p *Parser) Line() bool {
	// 等待到行尾
	tokens := []lexer.Token{}
	for {
		token := p.Lexer.Next()
		if token.IsEmpty() && len(tokens) == 0 {
			return false
		}
		if token.IsEmpty() || ((token.Value == "\n" || token.Value == ";") && token.Type == lexer.SEPARATOR) {
			p.line++ // 遇到换行符时增加行号
			break
		}
		tokens = append(tokens, token)
	}
	if len(tokens) == 0 {
		return true
	}
	if p.isPseudo(tokens[0]) {
		p.ParsePseudo(tokens)
		return true
	}
	if p.isInstructions(tokens[0]) {
		i := &Instruction{}
		i.Parse(tokens, p)
	}
	if p.isLabel(tokens) {
		l := &LabelBlock{}
		l.Parse(tokens, p)
	}
	return true
}

func (p *Parser) isInstructions(token lexer.Token) bool {
	instruction := types.Instruction(strings.ToUpper(token.Value))
	_, ok := p.arch.Instructions[instruction]
	if !ok {
		for i := 0; i < len(types.BuiltinI); i++ {
			if types.BuiltinI[i] == instruction {
				ok = true
				break
			}
		}
	}
	if token.IsEmpty() {
		return false
	}
	if token.Type == lexer.NAME && ok {
		token.Type = lexer.INSTRUCTION
		return true
	}
	return false
}

func (p *Parser) isLabel(tokens []lexer.Token) bool {
	if len(tokens) < 2 {
		return false
	}
	if tokens[0].Type == lexer.NAME && tokens[1].Type == lexer.SEPARATOR && tokens[1].Value == ":" {
		return true
	}
	return false
}

func (p *Parser) isPseudo(token lexer.Token) bool {
	if token.Type == lexer.PSEUDO {
		return true
	}
	return false
}

func (p *Parser) ParsePseudo(tokens []lexer.Token) {
	instruction := &Instruction{}
	instruction.ParseInstruction(tokens, p)
	switch instruction.Instruction {
	case "VAR":
		v := &VarBlock{}
		v.Parse(instruction, p)
	case "SECTION":
		//s := &SECTION{}
		//s.Parse(instruction, p)
	}
}
