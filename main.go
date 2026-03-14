package main

import (
    "fmt"
    "os"
    "assembler/internal/lexer"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <source.asm>")
        return
    }
    source := os.Args[1]

    file, err := os.Open(source)

    if err != nil {
        fmt.Println("error opening file ", source)
    }
    defer file.Close()

    l := lexer.New(file)
    //p := parser.New()
    //e := encoder.New()

    l.ScanTokens()

    //tokens := lexer.AllTokens(file)

    //ir := parser.Parse(tokens)
    //machineCode := encoder.Encode(ir)

    // write machine code to file
    //encoder.WriteBinary(machineCode, "output.bin")
}