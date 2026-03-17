// Instruction to binary

 package encoder

type Encoder struct {
	
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

func New() *Encoder {
	p := &Encoder{}
	return p
}