// Instruction to binary

 package encoder

import (
    "assembler/internal/parser"
    "assembler/internal/lexer"
    "os"
    "fmt"
    "encoding/binary"
)

type Encoder struct {}

type InstrEncoding struct {
	Opcode uint32
	Funct3 uint32
	Funct7 uint32 // only used for R-type, shift I-types
}

var InstrEncodings = map[string]InstrEncoding{
	// R-type (opcode 0x33)
	"add":  {0x33, 0x0, 0x00},
	"sub":  {0x33, 0x0, 0x20},
	"and":  {0x33, 0x7, 0x00},
	"or":   {0x33, 0x6, 0x00},
	"xor":  {0x33, 0x4, 0x00},
	"sll":  {0x33, 0x1, 0x00},
	"srl":  {0x33, 0x5, 0x00},
	"sra":  {0x33, 0x5, 0x20},
	"slt":  {0x33, 0x2, 0x00},
	"sltu": {0x33, 0x3, 0x00},

	// M-extension (opcode 0x33, funct7 0x01)
	"mul": {0x33, 0x0, 0x01},
	"div": {0x33, 0x4, 0x01},
	"rem": {0x33, 0x6, 0x01},

	// I-type arithmetic (opcode 0x13)
	"addi":  {0x13, 0x0, 0x00},
	"andi":  {0x13, 0x7, 0x00},
	"ori":   {0x13, 0x6, 0x00},
	"xori":  {0x13, 0x4, 0x00},
	"slti":  {0x13, 0x2, 0x00},
	"sltiu": {0x13, 0x3, 0x00},

	// Shifts are I-type but use funct7 in the upper imm bits
	"slli": {0x13, 0x1, 0x00},
	"srli": {0x13, 0x5, 0x00},
	"srai": {0x13, 0x5, 0x20},

	// Loads (opcode 0x03)
	"lw":  {0x03, 0x2, 0x00},
	"lh":  {0x03, 0x1, 0x00},
	"lhu": {0x03, 0x5, 0x00},
	"lb":  {0x03, 0x0, 0x00},
	"lbu": {0x03, 0x4, 0x00},

	// jalr (opcode 0x67)
	"jalr": {0x67, 0x0, 0x00},

	// ecall/ebreak (opcode 0x73)
	"ecall":  {0x73, 0x0, 0x00},
	"ebreak": {0x73, 0x0, 0x00}, // distinguished by imm=1 vs imm=0

	// S-type (opcode 0x23)
	"sw": {0x23, 0x2, 0x00},
	"sh": {0x23, 0x1, 0x00},
	"sb": {0x23, 0x0, 0x00},

	// B-type (opcode 0x63)
	"beq":  {0x63, 0x0, 0x00},
	"bne":  {0x63, 0x1, 0x00},
	"blt":  {0x63, 0x4, 0x00},
	"bltu": {0x63, 0x6, 0x00},
	"bge":  {0x63, 0x5, 0x00},
	"bgeu": {0x63, 0x7, 0x00},

	// U-type
	"lui":   {0x37, 0x0, 0x00},
	"auipc": {0x17, 0x0, 0x00},

	// J-type
	"jal": {0x6F, 0x0, 0x00},
}


func encodeR(funct7, rs2, rs1, funct3, rd, opcode uint32) uint32 {
    return (funct7 << 25) |
           (rs2    << 20) |
           (rs1    << 15) |
           (funct3 << 12) |
           (rd     <<  7) |
           opcode
}

func encodeI(imm, rs1, funct3, rd, opcode uint32) uint32 {
    return ((imm & 0xFFF) << 20) |
           (rs1    << 15) |
           (funct3 << 12) |
           (rd     <<  7) |
           opcode
}

func encodeS(imm, rs2, rs1, funct3, opcode uint32) uint32 {
    immLo := (imm & 0x1F)
    immHi := (imm >> 5) & 0x7F

    return (immHi  << 25) |
           (rs2    << 20) |
           (rs1    << 15) |
           (funct3 << 12) |
           (immLo  <<  7) |
           opcode
}

func encodeB(imm, rs2, rs1, funct3, opcode uint32) uint32 {
    imm12  := (imm >> 12) & 0x1
    imm11  := (imm >> 11) & 0x1
    imm105 := (imm >>  5) & 0x3F
    imm41  := (imm >>  1) & 0xF

    return (imm12  << 31) |
           (imm105 << 25) |
           (rs2    << 20) |
           (rs1    << 15) |
           (funct3 << 12) |
           (imm41  <<  8) |
           (imm11  <<  7) |
           opcode
}

func encodeU(imm, rd, opcode uint32) uint32 {
    return ((imm & 0xFFFFF) << 12) |
           (rd  <<  7) |
           opcode
}

func encodeJ(imm, rd, opcode uint32) uint32 {
    imm20   := (imm >> 20) & 0x1
    imm101  := (imm >>  1) & 0x3FF
    imm11   := (imm >> 11) & 0x1
    imm1912 := (imm >> 12) & 0xFF

    return (imm20   << 31) |
           (imm101  << 21) |
           (imm11   << 20) |
           (imm1912 << 12) |
           (rd      <<  7) |
           opcode
}

func (e *Encoder) Encode(prog *parser.Program) []uint32 { 

    var machineCode []uint32

    for _, instr := range prog.Instructions {

        i := InstrEncodings[instr.Op]
        r := lexer.Registers // map of register numbers
        
        var code uint32 
        var rd parser.Register
        var rs1 parser.Register
        var rs2 parser.Register
        var imm parser.Immediate
        var mem parser.Memory

        switch instr.Format {
        case parser.R:
            rd = instr.Args[0].(parser.Register)
            rs1 = instr.Args[1].(parser.Register)
            rs2 = instr.Args[2].(parser.Register)
            
            code = encodeR(i.Funct7, r[rs2.Name], r[rs1.Name], i.Funct3, r[rd.Name], i.Opcode)
        case parser.I:
            rd = instr.Args[0].(parser.Register)
            imm = instr.Args[1].(parser.Immediate)
            rs1 = instr.Args[2].(parser.Register)

            code = encodeI(uint32(imm.Value), r[rs1.Name], i.Funct3, r[rd.Name], i.Opcode)
        case parser.S:
            rs1 = instr.Args[0].(parser.Register)
            imm = instr.Args[1].(parser.Immediate)
            rs2 = instr.Args[2].(parser.Register)
            
            code = encodeS(uint32(imm.Value), r[rs2.Name], r[rs1.Name], i.Funct3, i.Opcode)
        case parser.B:
            rs1 = instr.Args[0].(parser.Register)
            rs2 = instr.Args[1].(parser.Register)
            mem = instr.Args[2].(parser.Memory)

            code = encodeB(uint32(mem.Offset), r[rs2.Name], r[rs1.Name], i.Funct3, i.Opcode)
        case parser.U:
            rd = instr.Args[0].(parser.Register)
            imm = instr.Args[1].(parser.Immediate)

            code = encodeU(uint32(imm.Value),r[rd.Name], i.Opcode)
        case parser.J:
            rd = instr.Args[0].(parser.Register)
            mem = instr.Args[1].(parser.Memory)

            code = encodeJ(uint32(mem.Offset),r[rd.Name], i.Opcode)
        }

        // add the line of code to machineCode
        machineCode = append(machineCode, code)
    }

    return machineCode
}

func (e *Encoder) WriteBinary(machineCode []uint32, output string) error {
    f, err := os.Create(output)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer f.Close()

    for i, word := range machineCode {
        fmt.Printf("%08x: %08x\n", i*4, word)
        err := binary.Write(f, binary.LittleEndian, word)
        if err != nil {
            return fmt.Errorf("failed to write word 0x%08X: %w", word, err)
        } 
    }

    return nil
}

func New() *Encoder {
	p := &Encoder{}
	return p
}