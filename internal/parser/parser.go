// Converts tokens to IR

package parser

import (
	"assembler/internal/lexer"
	"fmt"
	"strings"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
	line   int 
}

type Program struct {
	Instructions []Instruction
	Labels       map[string]int
}

type FormatType int

const (
	R FormatType = iota
	I
	S
	B
	U
	J
)

var InstrFormat map[string]FormatType = map[string]FormatType{
	"add": R,
	"sub": R,
	"and": R,

	"addi": I, 
	"lw": I, 
	"jalr": I, 
	"li": I,

	"sw": S, 
	"sb": S, 

	"beq": B,
	"bne": B,
	"lui": U,
	"auipc": U,

	"jal": J,
}

type Instruction struct {
	Op     string
	Args   []Operand
	Format FormatType
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

var LocationCounter int = 0

func isLabel(token lexer.Token) bool {
	if token.Type == lexer.IDENT {
		if strings.HasSuffix(token.Literal,":") {
			return true
		}
	}
	return false
}

func (p *Parser) ParseOperands(format FormatType) []Operand {

	var args []Operand

	switch format {
	case R:
		for range 3 {
			p.Next()
			token := p.GetToken()

			if token.Type != lexer.REGISTER {
				panic(fmt.Sprintf("Parser failed. R-type needs 3 regs, line %d", p.line))
			}

			args = append(args, Register{
				Name: token.Literal,
			})
		}
	case I: 
	case S: 
	case B:
	case U: 
	case J:

	}

	return args
}

func (p *Parser) ParseInstruction(token lexer.Token) (Instruction, error) {
	var op string 
	var format FormatType
	var args []Operand

	if token.Type == lexer.TokenType(lexer.IDENT) {
		op = token.Literal
		format = InstrFormat[op]
		args = p.ParseOperands(format)
	} else {
		return Instruction{}, fmt.Errorf("not ident")
	}

	return Instruction{
		Op: op,
		Args: args,
		Format: format,
	}, nil
}

func New() *Parser {
	p := &Parser{}
	return p
}

func (p *Parser) GetToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) Next() {
	p.pos += 1
}

func (p *Parser) Prev() {
	p.pos -= 1
}


func (p *Parser) Parse(tokens []lexer.Token) *Program {
	p.tokens = tokens 

	prog := &Program{
		Instructions: []Instruction{},
		Labels:       make(map[string]int),
	}

	for p.pos < len(tokens) {
		token := p.GetToken()
		if isLabel(token) {
			prog.Labels[token.Literal] = LocationCounter

		} else if token.Type == lexer.NEWLINE {
			p.line += 1 
		} else {
			instr, err := p.ParseInstruction(token)

			if err == nil {
				prog.Instructions = append(prog.Instructions, instr)
			}

			//fmt.Println(instr)
		}
		p.Next()
	}

	// run the backpatch list
	// check for undefined symbols

	return prog
}
