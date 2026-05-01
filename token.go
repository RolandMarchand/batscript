package main

import (
	"fmt"
)

type TokenKind int

type Token struct {
	kind   TokenKind
	lexeme string
	pos    int
}

const (
	ILLEGAL TokenKind = iota
	EOF
	EQUAL    // =
	BANG     // !
	COLON    // :
	COMMA    // ,
	DOT      // .
	ELLIPSIS // ...

	LPAREN // (
	RPAREN // )

	LBRACE // {
	RBRACE // }

	DOUBLE_EQUAL     // ==
	NOT_EQUAL        // !=
	LESSER_OR_EQUAL  // <=
	GREATER_OR_EQUAL // >=
	LESSER_THAN      // <
	GREATER_THAN     // >

	PLUS  // +
	MINUS // -

	STAR    // *
	SLASH   // /
	PERCENT // %

	NUMBER // 123
	STRING // "hello"
	SYMBOL // main

	// Keywords
	AND
	BREAK
	ELIF
	ELSE
	FUN
	IF
	INT
	NIL
	OBJ
	OR
	RETURN
	STR
	VAR
	WHILE
)

var tokens = [...]string{
	ILLEGAL:          "ILLEGAL",
	EOF:              "EOF",
	EQUAL:            "=",
	BANG:             "!",
	COLON:            ":",
	COMMA:            ",",
	DOT:              ".",
	ELLIPSIS:         "...",
	LPAREN:           "(",
	RPAREN:           ")",
	LBRACE:           "{",
	RBRACE:           "}",
	DOUBLE_EQUAL:     "==",
	NOT_EQUAL:        "!=",
	LESSER_OR_EQUAL:  "<=",
	GREATER_OR_EQUAL: ">=",
	LESSER_THAN:      "<",
	GREATER_THAN:     ">",
	PLUS:             "+",
	MINUS:            "-",
	STAR:             "*",
	SLASH:            "/",
	PERCENT:          "%",
	NUMBER:           "NUMBER",
	STRING:           "STRING",
	SYMBOL:           "SYMBOL",
	AND:              "and",
	BREAK:            "break",
	ELIF:             "elif",
	ELSE:             "else",
	FUN:              "fun",
	IF:               "if",
	INT:              "int",
	NIL:              "nil",
	OBJ:              "obj",
	OR:               "or",
	RETURN:           "return",
	STR:              "str",
	VAR:              "var",
	WHILE:            "while",
}

func (t TokenKind) String() string {
	var i = int(t)
	if i >= 0 && i < len(tokens) {
		return tokens[t]
	}
	return "UNKNOWN"
}

func (t Token) String() string {
	return fmt.Sprintf(
		"{kind=%s, lexeme=%s, pos=%d}",
		t.kind,
		t.lexeme,
		t.pos,
	)
}
