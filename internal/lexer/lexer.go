// Tokenizes assembly text

package lexer

import (
	"io"
	"io/ioutil"
	"fmt"
	"strings"
)

type Lexer struct {
    input        string
	tokens 		 []Token
    position     int
    readPosition int
    ch           byte
    line         int
    column       int
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT      // labels, mnemonics
	REGISTER   // R0, R1, R15
	IMMEDIATE  // #123, #0xFF
	NUMBER     // 123, 0xFF
	COMMA      // ,
	COLON      // :
	LBRACKET   // [
	RBRACKET   // ]
	DIRECTIVE  // .text, .data
	NEWLINE
)

type Token struct {
    Type     TokenType
    Literal  string
  //  Position Position
}

func New(r io.Reader) *Lexer {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	fmt.Println("new")

	l := &Lexer{
		input: string(data),
		line:  1,
		column: 0,
	}

	//l.readChar() // initialize first character
	return l
}

// func (l *Lexer) NextToken() Token {

// }

// func (l *Lexer) AllTokens() []Token {

// }

func (l *Lexer) ScanToken(token string) { 

	switch token {
	case "hello":
		fmt.Println("HELLLOOooo")
		l.tokens = append(l.tokens, Token{Type: TokenType(EOF), Literal: ""})
	default:
		fmt.Println("unknown token")
	}
}

func (l *Lexer) ScanTokens() []Token {
	words := strings.Fields(l.input)

	for _, word := range words {
		//fmt.Println(i, word)
		l.ScanToken(word)
	}

	return nil
}