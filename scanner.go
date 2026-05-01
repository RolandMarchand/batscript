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
	if text == nil {
		log.Fatal("null reference")
	}

	if !utf8.Valid(text) {
		tokens = nil
		err = errors.New("invalid encoding (requires UTF-8)")
		return
	}

	var numberPattern = regexp.MustCompile(`^\d+`)
	var symbolPattern = regexp.MustCompile(`^(\p{Letter}|_)(\p{Letter}|_|\d)*`)
	var ellipsisPattern = regexp.MustCompile(`^\.\.\.`)
	var commentPattern = regexp.MustCompile(`^#.*`)
	var spacePattern = regexp.MustCompile(`^[\f\t ]+`)
	var newlinePattern = regexp.MustCompile(`^(\r\n|\r|\n)`)
	var equalPattern = regexp.MustCompile(`^[=!<>]=`)
	var punctPattern = regexp.MustCompile(`^[=!:,.(){}<>*/%+-]`)
	var stringPattern = regexp.MustCompile(`^"(?:[^"\\]|\\.)*"`)

	var offset = 0
	var line = 1
	for len(text[offset:]) > 0 {
		if match := newlinePattern.Find(text[offset:]); match != nil {
			offset += len(match)
			line++
			continue
		}

		if match := commentPattern.Find(text[offset:]); match != nil {
			offset += len(match)
			continue
		}

		if match := spacePattern.Find(text[offset:]); match != nil {
			offset += len(match)
			continue
		}

		if match := ellipsisPattern.Find(text[offset:]); match != nil {
			var tok = Token{ELLIPSIS, "...", offset, line}
			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := numberPattern.Find(text[offset:]); match != nil {
			var tok = Token{NUMBER, string(match), offset, line}
			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := equalPattern.Find(text[offset:]); match != nil {
			var tok Token

			switch match[0] {
			case '=':
				tok = Token{DOUBLE_EQUAL, "==", offset, line}
			case '!':
				tok = Token{NOT_EQUAL, "!=", offset, line}
			case '<':
				tok = Token{LESSER_OR_EQUAL, "<=", offset, line}
			case '>':
				tok = Token{GREATER_OR_EQUAL, ">=", offset, line}
			default:
				log.Fatal("regex error")
			}

			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := punctPattern.Find(text[offset:]); match != nil {
			var tok Token

			switch match[0] {
			case '=':
				tok = Token{EQUAL, "=", offset, line}
			case '!':
				tok = Token{BANG, "!", offset, line}
			case ':':
				tok = Token{COLON, ":", offset, line}
			case ',':
				tok = Token{COMMA, ",", offset, line}
			case '.':
				tok = Token{DOT, ".", offset, line}
			case '(':
				tok = Token{LPAREN, "(", offset, line}
			case ')':
				tok = Token{RPAREN, ")", offset, line}
			case '{':
				tok = Token{LBRACE, "{", offset, line}
			case '}':
				tok = Token{RBRACE, "}", offset, line}
			case '<':
				tok = Token{LESSER_THAN, "<", offset, line}
			case '>':
				tok = Token{GREATER_THAN, ">", offset, line}
			case '+':
				tok = Token{PLUS, "+", offset, line}
			case '-':
				tok = Token{MINUS, "-", offset, line}
			case '*':
				tok = Token{STAR, "*", offset, line}
			case '/':
				tok = Token{SLASH, "/", offset, line}
			case '%':
				tok = Token{PERCENT, "%", offset, line}
			default:
				log.Fatal("regex error")
			}

			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := stringPattern.Find(text[offset:]); match != nil {
			var tok = Token{STRING, string(match), offset, line}
			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		if match := symbolPattern.Find(text[offset:]); match != nil {
			var symbol = string(match)
			var tok Token

			switch symbol {
			case "and":
				tok = Token{AND, symbol, offset, line}
			case "break":
				tok = Token{BREAK, symbol, offset, line}
			case "elif":
				tok = Token{ELIF, symbol, offset, line}
			case "else":
				tok = Token{ELSE, symbol, offset, line}
			case "fun":
				tok = Token{FUN, symbol, offset, line}
			case "if":
				tok = Token{IF, symbol, offset, line}
			case "int":
				tok = Token{INT, symbol, offset, line}
			case "nil":
				tok = Token{NIL, symbol, offset, line}
			case "obj":
				tok = Token{OBJ, symbol, offset, line}
			case "or":
				tok = Token{OR, symbol, offset, line}
			case "return":
				tok = Token{RETURN, symbol, offset, line}
			case "str":
				tok = Token{STR, symbol, offset, line}
			case "var":
				tok = Token{VAR, symbol, offset, line}
			case "while":
				tok = Token{WHILE, symbol, offset, line}
			default:
				tok = Token{SYMBOL, symbol, offset, line}
			}

			offset += len(match)
			tokens = append(tokens, tok)
			continue
		}

		var rune, size = utf8.DecodeRune(text[offset:])
		var tok = Token{ILLEGAL, string(rune), offset, line}
		offset += size
		tokens = append(tokens, tok)
	}

	var eof = Token{EOF, "", offset, line}
	tokens = append(tokens, eof)

	return
}
