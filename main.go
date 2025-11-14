package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"log"
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

// Global state for labels, operands, and registers
var labels []map[string]string
var operands []string
var registers []map[string]string

var debug = 0

// createBinaryFile appends binary string to binary.obj file
func createBinaryFile(binary string) error {
	f, err := os.OpenFile("binary.obj", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(binary)
	return err
}

// parseBranch parses B instruction
func parseBranch(line string, ins string, lineNumber int) {
	cond := ins[1:]
	idx := 0
	for _, x := range conditions {
		if x == cond {
			idx = findConditionIndex(cond)
			break
		}
	}
	condBin := fmt.Sprintf("%04b", idx)
	binary := condBin + "1010"
	createBinaryFile(binary)
}

// parseBranchWithLink parses BL instruction
func parseBranchWithLink(line string, ins string, lineNumber int) {
	cond := ins[1:]
	idx := findConditionIndex(cond)
	condBin := fmt.Sprintf("%04b", idx)
	binary := condBin + "1011"
	createBinaryFile(binary)
}

// parseBranchAndExchange parses BX instruction
func parseBranchAndExchange(line string, ins string, lineNumber int) {
	cond := ins[2:]
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	if len(parts) < 2 {
		fmt.Printf("Syntax Error: BX instruction requires a register\n")
		return
	}
	reg := parts[1]
	idx := findConditionIndex(cond)
	condBin := fmt.Sprintf("%04b", idx)
	
	regBin := findRegisterBinary(reg)
	binary := condBin + "000100101111111111110001" + regBin
	createBinaryFile(binary)
}

// parseSWP parses SWP (single word swap) instruction
func parseSWP(line string, ins string, lineNumber int) {
	line = strings.TrimSpace(line)
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == ','
	})
	
	if len(parts) < 4 {
		fmt.Printf("Syntax Error: SWP instruction format incorrect\n")
		return
	}
	
	destReg := parts[1]
	sourceReg := parts[2]
	baseReg := parts[3]
	
	destRegBin := findRegisterBinary(destReg)
	sourceRegBin := findRegisterBinary(sourceReg)
	baseRegBin := findRegisterBinary(baseReg)
	
	binary := "000000010000" + baseRegBin + destRegBin + "00001001" + sourceRegBin
	createBinaryFile(binary)
}

// parseSWI parses SWI (software interrupt) instruction
func parseSWI(line string, lineNumber int) {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	if len(parts) < 2 {
		fmt.Printf("Syntax Error: SWI instruction requires a value\n")
		return
	}
	val := parts[1]
	val = strings.TrimPrefix(val, "&")
	val = strings.ToLower(val)
	
	hexVal, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Printf("Syntax Error: Invalid hex value: %s\n", val)
		return
	}
	
	valBin := fmt.Sprintf("%024b", hexVal)
	binary := "0000" + "1111" + valBin
	createBinaryFile(binary)
}

// checkIfLabel checks if a label exists and returns its value
func checkIfLabel(reg string) string {
	for _, labelMap := range labels {
		for k, v := range labelMap {
			if k == reg {
				return v
			}
		}
	}
	return ""
}

// parseSDT parses single data transfer instruction (LDR/STR)
func parseSDT(line string, lineNumber int) {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	
	if len(parts) < 3 {
		fmt.Printf("Syntax Error: SDT instruction format incorrect at line %d\n", lineNumber)
		return
	}
	
	sdt := parts[0]
	bit := singleDataTransfer[parts[0]]
	
	if sdt == "LDR" {
		sourceReg := strings.TrimSuffix(parts[1], ",")
		baseReg := parts[2]
		
		newBaseReg := checkIfLabel(baseReg)
		if newBaseReg != "" {
			parts[2] = newBaseReg
		}
		
		hexVal := strings.TrimPrefix(newBaseReg, "&")
		hexVal = strings.ToLower(hexVal)
		intVal, err := strconv.ParseUint(hexVal, 16, 32)
		if err == nil {
			newBaseRegBin := fmt.Sprintf("%b", intVal)
			binary := "0000" + "01" + "0" + "0" + "0" + "1" + newBaseRegBin + "00000000"
			createBinaryFile(binary)
			registers = append(registers, map[string]string{sourceReg: newBaseRegBin})
		}
	} else {
		sourceReg := strings.TrimSuffix(parts[1], ",")
		sourceRegBinary := ""
		baseReg := parts[2]
		
		sourceRegBinary = findRegisterBinary(sourceReg)
		registers = append(registers, map[string]string{baseReg: sourceRegBinary})
		binary := "0000010001" + sourceRegBinary + "00000000"
		createBinaryFile(binary)
	}
}

// parseDPI parses data processing instruction
func parseDPI(line string, lineNumber int) {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	
	if len(parts) < 2 {
		fmt.Printf("Syntax Error: DPI instruction format incorrect at line %d\n", lineNumber)
		return
	}
	
	opcode := parts[0]
	opcodeBin := dataProc[opcode]
	
	var sourceReg, destReg, operandReg string
	
	if len(parts) > 3 {
		sourceReg = strings.TrimSuffix(parts[1], ",")
		destReg = strings.TrimSuffix(parts[2], ",")
		operandReg = strings.TrimSuffix(parts[3], ",")
	} else {
		sourceReg = strings.TrimSuffix(parts[1], ",")
		operandReg = strings.TrimSuffix(parts[2], ",")
		destReg = ""
	}
	
	destRegBin := findRegisterBinary(destReg)
	operandRegBin := findRegisterBinary(operandReg)
	
	binary := "0000" + "00" + "0" + opcodeBin + "1" + operandRegBin + destRegBin
	createBinaryFile(binary)
}

// parseLabel parses label declarations (DCW, DCD)
func parseLabel(line string, lineNumber int) {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	
	if len(parts) < 3 {
		fmt.Printf("Syntax Error! Cannot parse label due to spaces in line %d\n", lineNumber)
		fmt.Println("Tip: Replace spaces with tabs in source")
		return
	}
	
	labelName := parts[0]
	value := parts[2]
	labels = append(labels, map[string]string{labelName: value})
}

// findConditionIndex finds the index of a condition code
func findConditionIndex(cond string) int {
	for i, c := range conditions {
		if c == cond {
			return i
		}
	}
	return 0
}

// findRegisterBinary finds the binary representation of a register
func findRegisterBinary(reg string) string {
	for _, regMap := range registers {
		for k, v := range regMap {
			if k == reg {
				return v
			}
		}
	}
	return "0000"
}

// showDebug prints debug information
func showDebug() {
	fmt.Println("*** LABEL ***")
	for _, label := range labels {
		fmt.Println(label)
	}
	fmt.Println("*** REGISTERS ***")
	for _, reg := range registers {
		fmt.Println(reg)
	}
	fmt.Println("*** OPERANDS ***")
	for _, op := range operands {
		fmt.Println(op)
	}
}

// parseFileContent reads and parses the content of an assembly file
func parseFileContent(lines []string) {
	// First loop: parse labels
	for lineNumber, line := range lines {
		// Remove comments
		if idx := strings.Index(line, ";"); idx != -1 {
			line = line[:idx]
		}
		
		// Check for label declarations
		if strings.Contains(line, "DCW") || strings.Contains(line, "DCD") {
			parseLabel(line, lineNumber)
		}
	}
	
	// Second loop: parse instructions
	for lineNumber, line := range lines {
		// Remove comments
		if idx := strings.Index(line, ";"); idx != -1 {
			line = line[:idx]
		}
		
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check for branch instructions
		for i := 0; i < len(conditions); i++ {
			ins := "B" + conditions[i]
			if strings.Contains(line, ins) {
				parseBranch(line, ins, lineNumber)
				continue
			}
			ins = "BL" + conditions[i]
			if strings.Contains(line, ins) {
				parseBranchWithLink(line, ins, lineNumber)
				continue
			}
			ins = "BX" + conditions[i]
			if strings.Contains(line, ins) {
				parseBranchAndExchange(line, ins, lineNumber)
				continue
			}
		}
		
		// Check for data processing instructions
		for dpi := range dataProc {
			if strings.Contains(line, dpi) {
				parseDPI(line, lineNumber)
				break
			}
		}
		
		// Check for single data transfer instructions
		for sdt := range singleDataTransfer {
			if strings.Contains(line, sdt) {
				parseSDT(line, lineNumber)
				break
			}
		}
		
		// Check for software interrupts
		for swi := range softwareInterrupt {
			if strings.Contains(line, swi) {
				parseSWI(line, lineNumber)
				break
			}
		}
		
		// Check for SWP instruction
		if strings.Contains(line, "SWP") {
			parseSWP(line, "SWP", lineNumber)
		}
	}
}

// getFile processes an assembly file
func getFile(f string) {
	fmt.Println("Getting file:", f)
	
	// Check if file exists and is readable
	file, err := os.Open(f)
	if err != nil {
		fmt.Println("File Not Found!")
		return
	}
	defer file.Close()
	
	// Read file content
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	// Parse the file
	parseFileContent(lines)
	
	if debug == 1 {
		showDebug()
	}
	fmt.Println("Binary File 'binary.obj' has been created")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Type: go run main.go -f <fileName>.s")
		os.Exit(1)
	}
	
	flag := os.Args[1]
	fileName := os.Args[2]
	
	if flag == "-f" {
		getFile(fileName)
	} else {
		fmt.Println("Type: go run main.go -f <fileName>.s")
	}
}