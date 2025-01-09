package parser

import "CuteASM/lexer"

type LabelBlock struct {
	IsFunc bool
	Args   []*ArgBlock
	//Class      typeSys.Type
	//Return     []typeSys.Type
	Name string
	//BuildFlags []*Build
}

type ArgBlock struct {
	Name string
	//Type    typeSys.Type
	Defind *ArgBlock
	Offset int
}

func (l *LabelBlock) Parse(p *Parser) {
	if p.ThisBlock.Father != nil {
		p.Back(1)
	}
	beforeCursor := p.Lexer.Cursor
	code := p.Lexer.Next()
	if code.Type == lexer.LexTokenType["SEPARATOR"] && code.Value == "(" {
		l.IsFunc = true
		for {
			code1 := p.Lexer.Next()
			if (code1.Type != lexer.LexTokenType[""] && code1.Value == "BB" || code1.Value=="WW"||code1.Value == "DW"||code1.Value=="QW") {
				p.Error.MissError("", p.Lexer.Cursor, "")
			}
			
		}
	} else {
		p.Lexer.Cursor = beforeCursor
	}
	p.ThisBlock.AddChild(&Node{Value: l})
}
