// Tokenizes assembly text

package lexer

import (
	"io"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"unicode"
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
	NUMBER     // 123, 0xFF
	COMMA      // ,
	COLON      // :

	LPAREN      // (
	RPAREN      // )
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

	"x5":5, "t0":5, //temp 
	"x6":6, "t1":6,
	"x7":7, "t2":7,

	"x8":8, "s0":8, "fp":8,
	"x9":9, "s1":9,

	"x10":10, "a0":10, // args 
	"x11":11, "a1":11,
	"x12":12, "a2":12,
	"x13":13, "a3":13,
	"x14":14, "a4":14,
	"x15":15, "a5":15,
	"x16":16, "a6":16,
	"x17":17, "a7":17,

	"x18":18, "s2":18, //saved 
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


// func (l *Lexer) NextToken() Token {

// }

// func (l *Lexer) AllTokens() []Token {

// }

func isRegister(token string) bool {
	token = strings.TrimSuffix(token, ",")
	_, ok := registers[token]
	return ok
}


func isNumber(s string) bool {
	_, err := strconv.ParseInt(s, 0, 32)
	return err == nil
}

func isDirective(s string) bool {
	if strings.HasPrefix(s, ".") {
		return true 
	}
	return false
}

func isIdent(s string) bool {
	if len(s) == 0 { 
		return false
	}

	runes := []rune(s)

	// first character
	if !(unicode.IsLetter(runes[0]) || runes[0] == '_') {
		return false
	}

	// remaining characters
	for _, r := range runes[1:] {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	return true
}

func New(r io.Reader) *Lexer {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	fmt.Println("new lexer")

	l := &Lexer{
		input: string(data),
		line:  1,
		column: 0,
	}

	//l.readChar() // initialize first character
	return l
}

func (l *Lexer) ScanToken(token string) { 

	switch {
	case isRegister(token):
		l.CreateToken(REGISTER, token)
	case isNumber(token): 
		l.CreateToken(NUMBER, token)
	case token == "(": 
		l.CreateToken(LPAREN, "(")
	case token == ")": 
		l.CreateToken(RPAREN,")")
	case isIdent(token):
		l.CreateToken(IDENT, token)
	case isDirective(token):
		l.CreateToken(DIRECTIVE, token)
	default:
		fmt.Println("unknown token")
	}
}

func (l *Lexer) ScanTokens() []Token {
	// lol
	l.input = strings.ReplaceAll(l.input, "(", " ( ")
	l.input = strings.ReplaceAll(l.input, ")", " ) ")

	lines := strings.Split(l.input, "\n") 

	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			word = strings.ToLower(word)
			word = strings.TrimSuffix(word, ",")

			// if comment then go to next line
			if word[0] == '#' { 
				break 
			}

			//fmt.Println(word)
			l.ScanToken(word)
		}
		l.CreateToken(NEWLINE, "newline")
	}

	fmt.Println(l.tokens)
	return nil
}

func (l *Lexer) CreateToken(tType TokenType, literal string) {
	var t Token 
	t = Token{Type: TokenType(tType), Literal: literal}
	l.tokens = append(l.tokens, t)
}