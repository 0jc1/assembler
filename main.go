package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
**************** SOME FORMATS *************
_dpi_format = _InstructionFormat("Cond 0 0 I Opcode S Rn Rd Operand2")
_branch_format = _InstructionFormat("Cond 1 0 1 L Offset")
_bx_format = _InstructionFormat("Cond 0 0 0 1 0 0 1 0 1 1 1 1 1 1 1 1 1 1 1 1 0 0 L 1 Rm")
_load_store_format = _InstructionFormat("Cond 0 1 I P U B W L Rn Rd Operand2")
_load_store_multi_format = _InstructionFormat("Cond 1 0 0 P U S W L Rn RegisterList")
_mul_format = _InstructionFormat("Cond 0 0 0 0 0 0 0 S Rd 0 0 0 0 Rs 1 0 0 1 Rm")
_mla_format = _InstructionFormat("Cond 0 0 0 0 0 0 1 S Rd Rn Rs 1 0 0 1 Rm")
_clz_format = _InstructionFormat("Cond 0 0 0 1 0 1 1 0 1 1 1 1 Rd 1 1 1 1 0 0 0 1 Rm")
_mrs_format = _InstructionFormat("Cond 0 0 0 1 0 R 0 0 1 1 1 1 Rd 0 0 0 0 0 0 0 0 0 0 0 0")
_msr_format_reg = _InstructionFormat("Cond 0 0 0 1 0 R 1 0 f s x c 1 1 1 1 0 0 0 0 0 0 0 0 Rm")
_msr_format_imm = _InstructionFormat("Cond 0 0 1 1 0 R 1 0 f s x c 1 1 1 1 Operand2")
_swi_format = _InstructionFormat("Cond 1 1 1 1 Imm24")
*/

// List of ARM condition codes
var conditions = []string{
	"EQ", // Equal: Z set (zero flag)
	"NE", // Not equal: Z clear
	"CS", // Carry set / unsigned higher or same
	"CC", // Carry clear / unsigned lower
	"MI", // Minus / negative: N set
	"PL", // Plus / positive or zero: N clear
	"VS", // Overflow set: V set
	"VC", // Overflow clear: V clear
	"HI", // Unsigned higher: C set and Z clear
	"LS", // Unsigned lower or same: C clear or Z set
	"GE", // Signed greater or equal: N == V
	"LT", // Signed less than: N != V
	"GT", // Signed greater than: Z clear and N == V
	"LE", // Signed less or equal: Z set or N != V
}

// Data processing opcodes with their binary codes as strings
var dataProc = map[string]string{
	"AND": "0000", // Bitwise AND: Rd = Rn & Operand2
	"EOR": "0001", // Bitwise Exclusive OR: Rd = Rn ^ Operand2
	"SUB": "0010", // Subtract: Rd = Rn - Operand2
	"RSB": "0011", // Reverse Subtract: Rd = Operand2 - Rn
	"ADD": "0100", // Add: Rd = Rn + Operand2
	"ADC": "0101", // Add with Carry: Rd = Rn + Operand2 + Carry
	"SBC": "0110", // Subtract with Carry: Rd = Rn - Operand2 + Carry - 1
	"RSC": "0111", // Reverse Subtract with Carry: Rd = Operand2 - Rn + Carry - 1
	"TST": "1000", // Test (AND sets flags only): sets condition flags based on Rn & Operand2
	"TEQ": "1001", // Test Equivalence (EOR sets flags only): sets flags based on Rn ^ Operand2
	"CMP": "1010", // Compare (SUB sets flags only): sets flags based on Rn - Operand2
	"CMN": "1011", // Compare Negative (ADD sets flags only): sets flags based on Rn + Operand2
	"ORR": "1100", // Bitwise OR: Rd = Rn | Operand2
	"MOV": "1101", // Move: Rd = Operand2
	"BIC": "1110", // Bit Clear: Rd = Rn & ~Operand2
	"MVN": "1111", // Move Not: Rd = ~Operand2
}

// Single data transfer opcodes
var singleDataTransfer = map[string]int{
	"LDR": 1,
	"STR": 0,
}

// Software interrupt opcodes
var softwareInterrupt = map[string]int{
	"SWI": 0,
	"SVC": 1,
}

func parseFile(file *os.File) {
	// Placeholder for file parsing logic
	fmt.Println("Parsing file:", file.Name())
	// Create a scanner to read the file line by line
    var scanner *bufio.Scanner
	scanner = bufio.NewScanner(file)

	// Read and print each line
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filepath as an argument")
		os.Exit(1)
	}

	// Get the filepath from command line arguments
	var filepath string = os.Args[1]
    var file *os.File
    var err error
    
	file, err = os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	parseFile(file);
}