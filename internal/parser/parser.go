// Converts tokens to IR

package parser

import (
	"assembler/internal/lexer"
	"fmt"
	"strings"
	"strconv"
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
	// R-type
	"add":  R,
	"sub":  R,
	"and":  R,
	"or":   R,
	"xor":  R,
	"sll":  R,
	"srl":  R,
	"sra":  R,
	"slt":  R,
	"sltu": R,
	"mul":  R,
	"div":  R,
	"rem":  R,

	// I-type
	"addi":   I,
	"andi":   I,
	"ori":    I,
	"xori":   I,
	"slli":   I,
	"srli":   I,
	"srai":   I,
	"slti":   I,
	"sltiu":  I,
	"lw":     I,
	"lh":     I,
	"lhu":    I,
	"lb":     I,
	"lbu":    I,
	"jalr":   I,
	"ecall":  I,
	"ebreak": I,

	// S-type
	"sw": S,
	"sh": S,
	"sb": S,

	// B-type
	"beq":  B,
	"bne":  B,
	"blt":  B,
	"bltu": B,
	"bge":  B,
	"bgeu": B,

	// U-type
	"lui":   U,
	"auipc": U,

	// J-type
	"jal": J,

	// Pseudoinstructions
	// these are mapped to the format of
	// the real instructions they expand to
	"li":   I,
	"mv":   R,
	"nop":  I,
	"neg":  R,
	"not":  I,
	"j":    J,
	"jr":   I,
	"ret":  I,
	"call": J,
	"tail": J,
	"la":   U,
	"bgt":  B,
	"ble":  B,
	"bgtu": B,
	"bleu": B,
	"beqz": B,
	"bnez": B,
	"bltz": B,
	"bgez": B,
	"blez": B,
	"bgtz": B,
	"seqz": I,
	"snez": R,
	"sltz": R,
	"sgtz": R,
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
		if strings.HasSuffix(token.Literal, ":") {
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

		var rd lexer.Token
		var rs1 lexer.Token
		var imm lexer.Token

		p.Next()
		rd = p.GetToken()

		if rd.Type != lexer.REGISTER {
			panic("Parser failed. rd is not a register")
		}

		p.Next()

		token := p.GetToken()
		if token.Type == lexer.NUMBER {
			imm = token
			p.Next() // skip paren
			p.Next()
			rs1 = p.GetToken()
			if rs1.Type != lexer.REGISTER {
				panic("rs1 is not a register")
			}
			p.Next() // skip paren
		} else if token.Type == lexer.REGISTER {
			rs1 = token
			p.Next()
			imm = p.GetToken()
			if imm.Type != lexer.NUMBER {
				panic("imm is not a number")
			}
		}

		val, err := strconv.ParseInt(imm.Literal, 10, 32)
		if err != nil {
			panic(fmt.Sprintf("invalid immediate: %s", imm.Literal))
		}

		val2 := int32(val)

		args = append(args,
			Register{
				Name: rd.Literal,
			},
			Register{
				Name: rs1.Literal,
			},
			Immediate{
				Value: val2,
			},
		)

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
		Op:     op,
		Args:   args,
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
	p.pos++
}

func (p *Parser) Prev() {
	p.pos--
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
			p.line++
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
