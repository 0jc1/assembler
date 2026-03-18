package main

import (
    "fmt"
    "os"
    "assembler/internal/lexer"
    "assembler/internal/parser"
    "assembler/internal/encoder"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <source.asm>")
        os.Exit(1)
    }
    source := os.Args[1]
    file, err := os.Open(source)

    if err != nil {
        fmt.Println("Error opening file", source)
        os.Exit(1) 
    }
    defer file.Close()

    l := lexer.New(file)
    p := parser.New()
    e := encoder.New()

    tokens := l.ScanTokens()
    ir := p.Parse(tokens)
    fmt.Println(ir)
    machineCode := e.Encode(ir)
    fmt.Println(machineCode)

    // write machine code to file
    fmt.Println("Created output file output.bin")
    e.WriteBinary(machineCode, "output.bin")
}