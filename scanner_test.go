package main

import (
	"testing"
)

func TestOffset(t *testing.T) {
	var input = "one two three four"
	var expected = [...]int{0, 4, 8, 14, 18}
	var tokens, err = getTokens([]byte(input))

	if err != nil {
		t.Fatalf(`input "%s": %s`, input, err)
	}

	for i, token := range tokens {
		if i >= len(expected) {
			t.Fatalf("line %d: unexpected token %s",
				token.line,
				token.lexeme)
		}

		if token.pos != expected[i] {
			t.Errorf(`input "%s": expected offset %d, got %d`,
				input,
				expected[i],
				token.pos)
		}
	}
}

func TestLine(t *testing.T) {
	var input = `one
two three # four
five

six`
	var expected = [...]int{1, 2, 2, 3, 5, 5}
	var tokens, err = getTokens([]byte(input))

	if err != nil {
		t.Fatalf(`input "%s": %s`, input, err)
	}

	for i, token := range tokens {
		if i >= len(expected) {
			t.Fatalf("unexpected token %s", token.lexeme)
		}

		if token.line != expected[i] {
			t.Errorf(
				`input "%s": expected line %d, got %d`,
				input,
				expected[i],
				token.line,
			)
		}
	}
}

func TestWhitespace(t *testing.T) {
	var input = `hello # this should be ignored 123 😭
123 # there should be 2 tokens`
	var tokens, err = getTokens([]byte(input))

	if err != nil {
		t.Fatalf(`input "%s": %s`, input, err)
	}

	var expected = 3

	if len(tokens) != expected {
		t.Fatalf(
			`input "%s": expected %d tokens, got %d: %v`,
			input,
			expected,
			len(tokens),
			tokens,
		)
	}
}

func TestTokens(t *testing.T) {
	var tests = []struct {
		input    string
		expected Token
	}{
		{"😭", Token{ILLEGAL, "😭", 0, 1}},
		{"=", Token{EQUAL, "=", 0, 1}},
		{"!", Token{BANG, "!", 0, 1}},
		{":", Token{COLON, ":", 0, 1}},
		{",", Token{COMMA, ",", 0, 1}},
		{".", Token{DOT, ".", 0, 1}},
		{"...", Token{ELLIPSIS, "...", 0, 1}},
		{"(", Token{LPAREN, "(", 0, 1}},
		{")", Token{RPAREN, ")", 0, 1}},
		{"{", Token{LBRACE, "{", 0, 1}},
		{"}", Token{RBRACE, "}", 0, 1}},
		{"==", Token{DOUBLE_EQUAL, "==", 0, 1}},
		{"!=", Token{NOT_EQUAL, "!=", 0, 1}},
		{"<=", Token{LESSER_OR_EQUAL, "<=", 0, 1}},
		{">=", Token{GREATER_OR_EQUAL, ">=", 0, 1}},
		{"<", Token{LESSER_THAN, "<", 0, 1}},
		{">", Token{GREATER_THAN, ">", 0, 1}},
		{"+", Token{PLUS, "+", 0, 1}},
		{"-", Token{MINUS, "-", 0, 1}},
		{"*", Token{STAR, "*", 0, 1}},
		{"/", Token{SLASH, "/", 0, 1}},
		{"%", Token{PERCENT, "%", 0, 1}},
		{"0123456789", Token{NUMBER, "0123456789", 0, 1}},
		{`"hel\"o"`, Token{STRING, `"hel\"o"`, 0, 1}},
		{"a_symbol123", Token{SYMBOL, "a_symbol123", 0, 1}},
		{"and", Token{AND, "and", 0, 1}},
		{"break", Token{BREAK, "break", 0, 1}},
		{"elif", Token{ELIF, "elif", 0, 1}},
		{"else", Token{ELSE, "else", 0, 1}},
		{"fun", Token{FUN, "fun", 0, 1}},
		{"if", Token{IF, "if", 0, 1}},
		{"int", Token{INT, "int", 0, 1}},
		{"nil", Token{NIL, "nil", 0, 1}},
		{"obj", Token{OBJ, "obj", 0, 1}},
		{"or", Token{OR, "or", 0, 1}},
		{"return", Token{RETURN, "return", 0, 1}},
		{"str", Token{STR, "str", 0, 1}},
		{"var", Token{VAR, "var", 0, 1}},
		{"while", Token{WHILE, "while", 0, 1}},
	}

	for _, tt := range tests {
		var tokens, err = getTokens([]byte(tt.input))
		if err != nil {
			t.Errorf(
				"input %s (expected %s): %s",
				tt.input,
				tt.expected,
				err,
			)
			continue
		}

		var expected = 2

		if len(tokens) != expected {
			t.Errorf(
				"input %s: expected %d token, got %d: %v",
				tt.input,
				expected,
				len(tokens),
				tokens,
			)
			continue
		}

		if tokens[0] != tt.expected {
			t.Errorf(
				"input %s: expected %s",
				tt.input,
				tt.expected,
			)
		}
	}
}
