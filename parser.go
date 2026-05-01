package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	offset uint
}

func (p *Parser) peek() Token {
	if len(p.tokens) == 0 {
		panic("invalid argument, token input length 0")
	}

	if p.offset >= uint(len(p.tokens)) {
		return p.tokens[len(p.tokens)-1]
	}

	return p.tokens[p.offset]
}

func (p *Parser) consume(kinds ...TokenKind) Token {
	if len(p.tokens) == 0 {
		panic("invalid argument, token input length 0")
	}

	if p.offset >= uint(len(p.tokens)) {
		p.offset = uint(len(p.tokens)) - 1
	}

	var token = p.tokens[p.offset]
	p.offset++

	for _, kind := range kinds {
		if token.kind == kind {
			return token
		}
	}

	var err = fmt.Errorf(
		`Line %d: unexpected token "%s", expected any: %s`,
		token.line,
		token.lexeme,
		kinds,
	)
	panic(err)
}

func getAst(tokens []Token) Ast {
	if tokens == nil {
		panic("null reference")
	}

	var parser Parser = Parser{tokens, 0}

	return parseAst(&parser)
}

func parseAst(p *Parser) Ast {
	var ast Ast

	for {
		switch p.peek().kind {
		case FUN:
			var f Func = parseFunc(p)
			ast.funcs = append(ast.funcs, f)
		case VAR:
			var v Var = parseVar(p)
			ast.vars = append(ast.vars, v)
		default:
			p.consume(EOF)
			return ast
		}
	}
}

func parseVar(p *Parser) Var {
	var v Var

	p.consume(VAR)
	v.name = p.consume(SYMBOL).lexeme
	p.consume(COLON)
	v.kind = parseType(p)
	p.consume(EQUAL)
	v.assigned = parseExpr(p)

	return v
}

func parseType(p *Parser) Type {
	switch p.consume(INT, STR, FUN, OBJ).kind {
	case INT:
		return TYPE_INT
	case STR:
		return TYPE_STR
	case FUN:
		return TYPE_FUN
	case OBJ:
		return TYPE_OBJ
	default:
		panic("invalid state")
	}
}

func parseFunc(p *Parser) Func {
	var f Func

	p.consume(FUN)
	f.name = p.consume(SYMBOL).lexeme
	p.consume(LPAREN)
	f.params = parseParams(p)
	p.consume(RPAREN)
	f.ret = parseReturnType(p)
	p.consume(LBRACE)
	f.body = parseStmts(p)
	p.consume(RBRACE)

	return f
}

func parseReturnType(p *Parser) Type {
	if p.peek().kind != INT &&
		p.peek().kind != STR &&
		p.peek().kind != FUN &&
		p.peek().kind != OBJ {

		return TYPE_VOID
	}
	switch p.consume(INT, STR, FUN, OBJ).kind {
	case INT:
		return TYPE_INT
	case STR:
		return TYPE_STR
	case FUN:
		return TYPE_FUN
	case OBJ:
		return TYPE_OBJ
	default:
		panic("invalid state")
	}
}

func parseParams(p *Parser) Params {
	var params Params

	for {
		if p.peek().kind == ELLIPSIS {
			p.consume(ELLIPSIS)
			params.variadic = true
			return params
		}

		if p.peek().kind != SYMBOL {
			return params
		}

		params.vars = append(params.vars, parseParam(p))

		if p.peek().kind == COMMA {
			p.consume(COMMA)
		}
	}
}

func parseParam(p *Parser) Var {
	var v Var

	v.name = p.consume(SYMBOL).lexeme
	p.consume(COLON)
	v.kind = parseType(p)
	p.consume(EQUAL)
	v.assigned = parseExpr(p)

	return v
}

func parseStmts(p *Parser) []Stmt {
	var stmts []Stmt

	for {
		switch p.peek().kind {
		case VAR:
			stmts = append(stmts, StmtVar(parseVar(p)))
		case IF:
			stmts = append(stmts, parseIf(p))
		case WHILE:
			stmts = append(stmts, parseWhile(p))
		case RETURN:
			stmts = append(stmts, parseReturn(p))
		case BREAK:
			stmts = append(stmts, parseBreak(p))
		case RBRACE:
			return stmts
		default:
			stmts = append(stmts, parseStmtExpr(p))
		}
	}
}

func parseIf(p *Parser) StmtIf {
	var stmt StmtIf

	p.consume(IF)
	stmt.test = parseExpr(p)
	p.consume(LBRACE)
	stmt.body = parseStmts(p)
	p.consume(RBRACE)
	stmt.elifs = parseStmtElif(p)
	stmt.elseBody = parseStmtElse(p)

	return stmt
}

func parseStmtElif(p *Parser) []StmtElif {
	var stmts []StmtElif

	for {
		if p.peek().kind != ELIF {
			return stmts
		}

		var stmt StmtElif

		p.consume(ELIF)
		stmt.test = parseExpr(p)
		p.consume(LBRACE)
		stmt.body = parseStmts(p)
		p.consume(RBRACE)

		stmts = append(stmts, stmt)
	}
}

func parseStmtElse(p *Parser) []Stmt {
	var stmts []Stmt

	if p.peek().kind != ELSE {
		return stmts
	}

	p.consume(ELSE)
	p.consume(LBRACE)
	stmts = parseStmts(p)
	p.consume(RBRACE)

	return stmts
}

func parseWhile(p *Parser) StmtWhile {
	var stmt StmtWhile

	p.consume(WHILE)
	stmt.test = parseExpr(p)
	p.consume(LBRACE)
	stmt.body = parseStmts(p)
	p.consume(RBRACE)

	return stmt
}

func parseReturn(p *Parser) StmtReturn {
	p.consume(RETURN)
	return StmtReturn{parseExpr(p)}
}

func parseBreak(p *Parser) StmtBreak {
	p.consume(BREAK)
	return StmtBreak{}
}

func parseStmtExpr(p *Parser) StmtExpr {
	return StmtExpr{parseExpr(p)}
}

func parseExpr(p *Parser) Expr {
	return parseExprAssign(p)
}

func parseExprAssign(p *Parser) Expr {
	var expr = parseExprLogic(p)

	if p.peek().kind == EQUAL {
		p.consume(EQUAL)
		expr = ExprBinary{expr, parseExprAssign(p), OP_ASS}
	}

	return expr
}

func parseExprLogic(p *Parser) Expr {
	var expr = parseExprCmp(p)

	for p.peek().kind == AND || p.peek().kind == OR {
		var op Op
		if p.consume(AND, OR).kind == AND {
			op = OP_AND
		} else {
			op = OP_OR
		}

		expr = ExprBinary{expr, parseExprCmp(p), op}
	}

	return expr
}

func parseExprCmp(p *Parser) Expr {
	var expr = parseExprAdd(p)

	for p.peek().kind == DOUBLE_EQUAL ||
		p.peek().kind == NOT_EQUAL ||
		p.peek().kind == LESSER_OR_EQUAL ||
		p.peek().kind == GREATER_OR_EQUAL ||
		p.peek().kind == LESSER_THAN ||
		p.peek().kind == GREATER_THAN {

		var op Op
		switch p.consume(
			DOUBLE_EQUAL,
			NOT_EQUAL,
			LESSER_OR_EQUAL,
			GREATER_OR_EQUAL,
			LESSER_THAN,
			GREATER_THAN,
		).kind {
		case DOUBLE_EQUAL:
			op = OP_EQ
		case NOT_EQUAL:
			op = OP_NEQ
		case LESSER_OR_EQUAL:
			op = OP_LE
		case GREATER_OR_EQUAL:
			op = OP_GE
		case LESSER_THAN:
			op = OP_LT
		case GREATER_THAN:
			op = OP_GT
		default:
			panic("invalid state")
		}

		expr = ExprBinary{expr, parseExprAdd(p), op}
	}

	return expr
}

func parseExprAdd(p *Parser) Expr {
	var expr = parseExprMult(p)

	for p.peek().kind == PLUS || p.peek().kind == MINUS {
		var op Op
		if p.consume(PLUS, MINUS).kind == PLUS {
			op = OP_ADD
		} else {
			op = OP_SUB
		}

		expr = ExprBinary{expr, parseExprMult(p), op}
	}

	return expr
}

func parseExprMult(p *Parser) Expr {
	var expr = parseExprUnary(p)

	for p.peek().kind == STAR ||
		p.peek().kind == SLASH ||
		p.peek().kind == PERCENT {

		var op Op
		switch p.consume(STAR, SLASH, PERCENT).kind {
		case STAR:
			op = OP_MUL
		case SLASH:
			op = OP_DIV
		case PERCENT:
			op = OP_MOD
		default:
			panic("invalid state")
		}

		expr = ExprBinary{expr, parseExprUnary(p), op}
	}

	return expr
}

func parseExprUnary(p *Parser) Expr {
	if p.peek().kind != PLUS &&
		p.peek().kind != MINUS &&
		p.peek().kind != BANG {

		return parseExprPostfix(p)
	}

	var op Op
	switch p.consume(PLUS, MINUS, BANG).kind {
	case PLUS:
		op = OP_ADD
	case MINUS:
		op = OP_SUB
	case BANG:
		op = OP_NOT
	default:
		panic("invalid state")
	}

	return ExprUnary{parseExprUnary(p), op}
}

func parseExprPostfix(p *Parser) Expr {
	var expr = parseExprPrimary(p)

	for {
		if p.peek().kind == DOT {
			p.consume(DOT)
			var field = p.consume(SYMBOL).lexeme
			expr = ExprField{expr, field}
		} else if p.peek().kind == LPAREN {
			p.consume(LPAREN)
			var args = parseArgs(p)
			p.consume(RPAREN)
			expr = ExprCall{expr, args}
		} else {
			break
		}
	}

	return expr
}

func parseArgs(p *Parser) []Expr {
	var args []Expr

	for {
		if p.peek().kind == RPAREN {
			return args
		}

		args = append(args, parseExpr(p))

		if p.peek().kind == COMMA {
			p.consume(COMMA)
		}
	}
}

func parseExprPrimary(p *Parser) Expr {
	var expr Expr

	if p.peek().kind == LPAREN {
		p.consume(LPAREN)
		expr = parseExpr(p)
		p.consume(RPAREN)
	} else {
		expr = parseValue(p)
	}

	return expr
}

func parseValue(p *Parser) Expr {
	switch p.peek().kind {
	case NUMBER:
		return parseExprNumber(p)
	case STRING:
		return parseExprString(p)
	case SYMBOL:
		return parseExprSymbol(p)
	case OBJ:
		return parseExprObj(p)
	case NIL:
		return parseExprNil(p)
	default:
		return ExprNil{}
	}
}

func parseExprNumber(p *Parser) ExprNumber {
	var expr ExprNumber

	var numText = p.consume(NUMBER).lexeme
	var number, err = strconv.ParseInt(numText, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("%s too big (max 9223372036854775807)",
			numText))
	}

	expr.value = number

	return expr
}

func parseExprString(p *Parser) ExprString {
	var expr ExprString

	var str = p.consume(STRING).lexeme
	var unqoted, _ = strconv.Unquote(str)
	expr.value = unqoted

	return expr
}

func parseExprSymbol(p *Parser) ExprSymbol {
	return ExprSymbol{p.consume(SYMBOL).lexeme}
}

func parseExprObj(p *Parser) ExprObject {
	var obj ExprObject

	p.consume(OBJ)
	p.consume(LBRACE)

	var vars []Var
	for {
		if p.peek().kind != SYMBOL {
			break
		}

		var v Var

		v.name = p.consume(SYMBOL).lexeme
		p.consume(COLON)
		v.kind = parseType(p)
		p.consume(EQUAL)
		v.assigned = parseExpr(p)

		vars = append(vars, v)

		if p.peek().kind == COMMA {
			p.consume(COMMA)
		}
	}
	obj.fields = vars

	p.consume(RBRACE)

	return obj
}

func parseExprNil(p *Parser) ExprNil {
	p.consume(NIL)
	return ExprNil{}
}
