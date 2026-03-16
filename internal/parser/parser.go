// Converts tokens to IR

package parser

import (
	"assembler/internal/lexer"
	"fmt"
)

type Parser struct {
    tokens []lexer.Token
    pos    int	
}

type Program struct {
    Instructions []Instruction
    Labels       map[string]int
}

type Instruction struct {
    Op   string
    Args []Operand
}

type Operand interface{}

type Register struct {
    Name string
}

type Immediate struct {
    Value int32
}

type Memory struct {
    Offset int32
    Base   string
}

type LabelRef struct {
    Name string
}

func isLabel(token lexer.Token) bool {
	return true 
}	

func parseInstruction() string {
	return ""
}

func New() *Parser {
	p := &Parser{}
	return p
}

func (p *Parser) Parse(tokens []lexer.Token) *Program {
    prog := &Program{}


	for _, token := range tokens {
		if isLabel(token) {
			//
		} else {
			instr := parseInstruction()
			fmt.Println(instr)

		}
	}

    return prog
}