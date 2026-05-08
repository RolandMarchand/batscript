package main

type Ast struct {
	funcs []Func
	vars  []Var
}

type Func struct {
	name   string
	params Params
	ret    Type
	body   []Stmt
}

type Params struct {
	vars     []Var
	variadic bool
}

type Var struct {
	name     string
	kind     Type
	assigned Expr
}

type Type int

const (
	TYPE_INT Type = iota
	TYPE_STR
	TYPE_FUN
	TYPE_OBJ
	TYPE_VOID
)

func (t Type) String() string {
	switch t {
	case TYPE_INT: return "int"
	case TYPE_STR: return "string"
	case TYPE_FUN: return "function"
	case TYPE_OBJ: return "object"
	case TYPE_VOID: return "void"
	default: return "unknown"
	}
}

type Stmt interface {
	stmtNode()
}

type StmtVar Var

func (s StmtVar) stmtNode() {}

type StmtExpr struct {
	expr Expr
}

func (s StmtExpr) stmtNode() {}

type StmtIf struct {
	test     Expr
	body     []Stmt
	elifs    []StmtElif
	elseBody []Stmt
}

func (s StmtIf) stmtNode() {}

type StmtElif struct {
	test Expr
	body []Stmt
}

type StmtWhile struct {
	test Expr
	body []Stmt
}

func (s StmtWhile) stmtNode() {}

type StmtReturn struct {
	value Expr
}

func (s StmtReturn) stmtNode() {}

type StmtBreak struct{}

func (s StmtBreak) stmtNode() {}

type Expr interface {
	exprNode()
}

type ExprBinary struct {
	left  Expr
	right Expr
	op    Op
}

func (e ExprBinary) exprNode() {}

type ExprUnary struct {
	value Expr
	op    Op
}

func (e ExprUnary) exprNode() {}

type Op int

const (
	OP_ADD Op = iota
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_AND
	OP_OR
	OP_EQ
	OP_NEQ
	OP_LE
	OP_GE
	OP_LT
	OP_GT
	OP_NOT
	OP_ASS
)

type ExprCall struct {
	callee Expr
	args   []Expr
}

func (e ExprCall) exprNode() {}

type ExprField struct {
	object Expr
	field  string
}

func (e ExprField) exprNode() {}

type ExprNumber struct {
	value int64
}

func (e ExprNumber) exprNode() {}

type ExprString struct {
	value string
}

func (e ExprString) exprNode() {}

type ExprSymbol struct {
	value string
}

func (e ExprSymbol) exprNode() {}

type ExprObject struct {
	fields []Var
}

func (e ExprObject) exprNode() {}

type ExprNil struct{}

func (e ExprNil) exprNode() {}
