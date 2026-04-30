package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"unicode/utf8"
)

func getTokensFromFile(filename string) ([]Token, error) {
	var text, err = os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return getTokens(text)
}

func getTokens(text []byte) (tokens []Token, err error) {

	if !utf8.Valid(text) {
		tokens = nil
		err = errors.New("invalid encoding (requires UTF-8)")
		return
	}

	var numberPattern = regexp.MustCompile(`^\d+`)
	var symbolPattern = regexp.MustCompile(`^\p{Letter}(\p{Letter}|\d)*`)
	var ellipsisPattern = regexp.MustCompile(`^\.\.\.`)
	var commentPattern = regexp.MustCompile(`^#.*`)
	var spacePattern = regexp.MustCompile(`^\s+`)
	var equalPattern = regexp.MustCompile(`^[=!<>]=`)
	var punctPattern = regexp.MustCompile(`^[=!:,.(){}<>*/%+-]`)
	var stringPattern = regexp.MustCompile(`^"(?:[^"\\]|\\.)*"`)

	var offset = 0
	for len(text[offset:]) > 0 {
		if match := commentPattern.Find(text[offset:]); match != nil {
			offset += len(match)
			continue
		}

		if match := spacePattern.Find(text[offset:]); match != nil {
			offset += len(match)
			continue
		}

		if match := ellipsisPattern.Find(text[offset:]); match != nil {
			var tok = Token{ELLIPSIS, "...", offset}
			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := numberPattern.Find(text[offset:]); match != nil {
			var tok = Token{NUMBER, string(match), offset}
			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := equalPattern.Find(text[offset:]); match != nil {
			var tok Token

			switch match[0] {
			case '=': tok = Token{DOUBLE_EQUAL, "==", offset}
			case '!': tok = Token{NOT_EQUAL, "!=", offset}
			case '<': tok = Token{LESSER_THAN_EQUAL, "<=", offset}
			case '>': tok = Token{GREATER_THAN_EQUAL, ">=", offset}
			default: log.Fatal("regex error")
			}

			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := punctPattern.Find(text[offset:]); match != nil {
			var tok Token

			switch match[0] {
			case '=': tok = Token{EQUAL, "=", offset}
			case '!': tok = Token{BANG, "!", offset}
			case ':': tok = Token{COLON, ":", offset}
			case ',': tok = Token{COMMA, ",", offset}
			case '.': tok = Token{DOT, ".", offset}
			case '(': tok = Token{LPAREN, "(", offset}
			case ')': tok = Token{RPAREN, ")", offset}
			case '{': tok = Token{LBRACE, "{", offset}
			case '}': tok = Token{RBRACE, "}", offset}
			case '<': tok = Token{LESSER_THAN, "<", offset}
			case '>': tok = Token{GREATER_THAN, ">", offset}
			case '+': tok = Token{PLUS, "+", offset}
			case '-': tok = Token{MINUS, "-", offset}
			case '*': tok = Token{STAR, "*", offset}
			case '/': tok = Token{SLASH, "/", offset}
			case '%': tok = Token{PERCENT, "%", offset}
			default: log.Fatal("regex error")
			}

			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := stringPattern.Find(text[offset:]); match != nil {
			var tok = Token{STRING, string(match), offset}
			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := symbolPattern.Find(text[offset:]); match != nil {
			var symbol = string(match)
			var tok Token

			switch symbol {
			case "and": tok = Token{AND, symbol, offset}
			case "break": tok = Token{BREAK, symbol, offset}
			case "elif": tok = Token{ELIF, symbol, offset}
			case "else": tok = Token{ELSE, symbol, offset}
			case "fun": tok = Token{FUN, symbol, offset}
			case "if": tok = Token{IF, symbol, offset}
			case "int": tok = Token{INT, symbol, offset}
			case "nil": tok = Token{NIL, symbol, offset}
			case "non": tok = Token{NON, symbol, offset}
			case "obj": tok = Token{OBJ, symbol, offset}
			case "or": tok = Token{OR, symbol, offset}
			case "return": tok = Token{RETURN, symbol, offset}
			case "str": tok = Token{STR, symbol, offset}
			case "var": tok = Token{VAR, symbol, offset}
			case "while": tok = Token{WHILE, symbol, offset}
			default: tok = Token{SYMBOL, symbol, offset}
			}

			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		var rune, size = utf8.DecodeRune(text[offset:])
		var tok = Token{ILLEGAL, string(rune), offset}
		offset += size		
		tokens = append(tokens, tok)
	}

	return
}
