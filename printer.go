package main

import (
	"fmt"
	"strings"
)

type Printer struct {
	depth int
}

func (p *Printer) indent() string {
	return strings.Repeat("  ", p.depth)
}

func (p *Printer) line(format string, args ...any) {
	fmt.Printf(p.indent()+format+"\n", args...)
}

func (p *Printer) enter(format string, args ...any) {
	p.line(format, args...)
	p.depth++
}

func (p *Printer) leave() {
	p.depth--
}

func printAst(ast Ast) {
	p := &Printer{}
	p.line("PROGRAM")
	p.depth++
	for _, v := range ast.vars {
		printVar(p, v)
	}
	for _, f := range ast.funcs {
		printFunc(p, f)
	}
}

func printFunc(p *Printer, f Func) {
	p.enter("FUNC %s", f.name)
	printParams(p, f.params)
	p.line("RET %s", printType(f.ret))
	p.enter("BODY")
	for _, s := range f.body {
		printStmt(p, s)
	}
	p.leave()
	p.leave()
}

func printParams(p *Printer, params Params) {
	p.enter("PARAMS variadic=%v", params.variadic)
	for _, v := range params.vars {
		printVar(p, v)
	}
	p.leave()
}

func printVar(p *Printer, v Var) {
	p.enter("VAR %s : %s", v.name, printType(v.kind))
	printExpr(p, v.assigned)
	p.leave()
}

func printType(t Type) string {
	switch t {
	case TYPE_INT:
		return "int"
	case TYPE_STR:
		return "str"
	case TYPE_FUN:
		return "fun"
	case TYPE_OBJ:
		return "obj"
	case TYPE_VOID:
		return "void"
	default:
		return "unknown"
	}
}

func printStmt(p *Printer, s Stmt) {
	switch s := s.(type) {
	case StmtVar:
		printVar(p, Var(s))
	case StmtExpr:
		p.enter("EXPR_STMT")
		printExpr(p, s.expr)
		p.leave()
	case StmtIf:
		printIf(p, s)
	case StmtWhile:
		p.enter("WHILE")
		p.enter("TEST")
		printExpr(p, s.test)
		p.leave()
		p.enter("BODY")
		for _, stmt := range s.body {
			printStmt(p, stmt)
		}
		p.leave()
		p.leave()
	case StmtReturn:
		p.enter("RETURN")
		printExpr(p, s.value)
		p.leave()
	case StmtBreak:
		p.line("BREAK")
	}
}

func printIf(p *Printer, s StmtIf) {
	p.enter("IF")
	p.enter("TEST")
	printExpr(p, s.test)
	p.leave()
	p.enter("BODY")
	for _, stmt := range s.body {
		printStmt(p, stmt)
	}
	p.leave()
	for _, elif := range s.elifs {
		p.enter("ELIF")
		p.enter("TEST")
		printExpr(p, elif.test)
		p.leave()
		p.enter("BODY")
		for _, stmt := range elif.body {
			printStmt(p, stmt)
		}
		p.leave()
		p.leave()
	}
	if len(s.elseBody) > 0 {
		p.enter("ELSE")
		for _, stmt := range s.elseBody {
			printStmt(p, stmt)
		}
		p.leave()
	}
	p.leave()
}

func printExpr(p *Printer, e Expr) {
	switch e := e.(type) {
	case ExprBinary:
		p.enter("BINARY %s", printOp(e.op))
		printExpr(p, e.left)
		printExpr(p, e.right)
		p.leave()
	case ExprUnary:
		p.enter("UNARY %s", printOp(e.op))
		printExpr(p, e.value)
		p.leave()
	case ExprCall:
		p.enter("CALL")
		p.enter("CALLEE")
		printExpr(p, e.callee)
		p.leave()
		p.enter("ARGS")
		for _, arg := range e.args {
			printExpr(p, arg)
		}
		p.leave()
		p.leave()
	case ExprField:
		p.enter("FIELD .%s", e.field)
		printExpr(p, e.object)
		p.leave()
	case ExprNumber:
		p.line("NUMBER %d", e.value)
	case ExprString:
		p.line("STRING %q", e.value)
	case ExprSymbol:
		p.line("SYMBOL %s", e.value)
	case ExprObject:
		p.enter("OBJECT")
		for _, f := range e.fields {
			printVar(p, f)
		}
		p.leave()
	case ExprNil:
		p.line("NIL")
	}
}

func printOp(op Op) string {
	switch op {
	case OP_ADD:
		return "+"
	case OP_SUB:
		return "-"
	case OP_MUL:
		return "*"
	case OP_DIV:
		return "/"
	case OP_MOD:
		return "%"
	case OP_AND:
		return "and"
	case OP_OR:
		return "or"
	case OP_EQ:
		return "=="
	case OP_NEQ:
		return "!="
	case OP_LE:
		return "<="
	case OP_GE:
		return ">="
	case OP_LT:
		return "<"
	case OP_GT:
		return ">"
	case OP_NOT:
		return "!"
	case OP_ASS:
		return "="
	default:
		return "unknown"
	}
}
