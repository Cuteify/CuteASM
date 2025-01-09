package parser

import (
	"CuteASM/lexer"
)

type SECTION struct {
	Name string
	Desc string
}

func (s *SECTION) Parse(p *Parser) {
	code := p.Lexer.Next()
	if code.IsEmpty() || code.Type != lexer.LexTokenType["NAME"] {
		p.Lexer.Error.MissError("Syntax Error", p.Lexer.Cursor, "Need section Name")
	}
	if code.Value[0] != '.' {
		s.Name = "." + code.Value
	} else {
		s.Name = code.Value
	}
	p.ThisBlock.AddChild(&Node{
		Value: s,
	})
}
