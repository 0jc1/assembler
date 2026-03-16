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

var registers = map[string]int{
	"x0":0, "zero":0,
	"x1":1, "ra":1,
	"x2":2, "sp":2,
	"x3":3, "gp":3,
	"x4":4, "tp":4,

	"x5":5, "t0":5,
	"x6":6, "t1":6,
	"x7":7, "t2":7,

	"x8":8, "s0":8, "fp":8,
	"x9":9, "s1":9,

	"x10":10, "a0":10,
	"x11":11, "a1":11,
	"x12":12, "a2":12,
	"x13":13, "a3":13,
	"x14":14, "a4":14,
	"x15":15, "a5":15,
	"x16":16, "a6":16,
	"x17":17, "a7":17,

	"x18":18, "s2":18,
	"x19":19, "s3":19,
	"x20":20, "s4":20,
	"x21":21, "s5":21,
	"x22":22, "s6":22,
	"x23":23, "s7":23,
	"x24":24, "s8":24,
	"x25":25, "s9":25,
	"x26":26, "s10":26,
	"x27":27, "s11":27,

	"x28":28, "t3":28,
	"x29":29, "t4":29,
	"x30":30, "t5":30,
	"x31":31, "t6":31,
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

func CreateToken(tType TokenType, literal string) Token {
	return Token{Type: TokenType(tType), Literal: literal}
}

func isRegister(token string) bool {
	token = strings.TrimSuffix(token, ",")
	_, ok := registers[token]
	return ok
}

func (l *Lexer) ScanToken(token string) { 

	switch {
	case isRegister(token):
		fmt.Println(token)
		l.tokens = append(l.tokens, CreateToken(REGISTER, token))
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