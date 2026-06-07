package main

type Node interface {
	node()
}

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
	Node
	stmtNode()
}

type StmtVar Var

type StmtExpr struct {
	expr Expr
}

type StmtIf struct {
	test     Expr
	body     []Stmt
	elifs    []StmtElif
	elseBody []Stmt
}

type StmtElif struct {
	test Expr
	body []Stmt
}

type StmtWhile struct {
	test Expr
	body []Stmt
}

type StmtReturn struct {
	value Expr
}

type StmtBreak struct{}

type Expr interface {
	Node
	exprNode()
}

type ExprBinary struct {
	left  Expr
	right Expr
	op    Op
}

type ExprUnary struct {
	value Expr
	op    Op
}

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

type ExprField struct {
	object Expr
	field  string
}

type ExprNumber struct {
	value int64
}

type ExprString struct {
	value string
}

type ExprSymbol struct {
	value string
}

type ExprObject struct {
	fields []Var
}

type ExprNil struct{}

func (a Ast) node(){}
func (f Func) node(){}
func (p Params) node(){}
func (v Var) node(){}
func (t Type) node(){}
func (s StmtVar) node() {}
func (s StmtExpr) node() {}
func (s StmtIf) node() {}
func (s StmtElif) node() {}
func (s StmtWhile) node() {}
func (s StmtReturn) node() {}
func (s StmtBreak) node() {}
func (e ExprBinary) node() {}
func (e ExprUnary) node() {}
func (e ExprCall) node() {}
func (e ExprField) node() {}
func (e ExprNumber) node() {}
func (e ExprString) node() {}
func (e ExprSymbol) node() {}
func (e ExprObject) node() {}
func (e ExprNil) node() {}

func (s StmtVar) stmtNode() {}
func (s StmtExpr) stmtNode() {}
func (s StmtIf) stmtNode() {}
func (s StmtElif) stmtNode() {}
func (s StmtWhile) stmtNode() {}
func (s StmtReturn) stmtNode() {}
func (s StmtBreak) stmtNode() {}
func (e ExprBinary) exprNode() {}
func (e ExprUnary) exprNode() {}
func (e ExprCall) exprNode() {}
func (e ExprField) exprNode() {}
func (e ExprNumber) exprNode() {}
func (e ExprString) exprNode() {}
func (e ExprSymbol) exprNode() {}
func (e ExprObject) exprNode() {}
func (e ExprNil) exprNode() {}

func findNodes[T Node](root Node) []T {
	if root == nil {
		return nil
	}

	var result []T

	if t, ok := root.(T); ok {
		result = append(result, t)
	}

	switch e := root.(type) {
	// Other
	case Ast:
		for _, n := range e.funcs {
			result = append(result, findNodes[T](n)...)
		}
		for _, n := range e.vars {
			result = append(result, findNodes[T](n)...)
		}
	case Func:
		result = append(result, findNodes[T](e.params)...)
		result = append(result, findNodes[T](e.ret)...)
		for _, n := range e.body {
			result = append(result, findNodes[T](n)...)
		}
	case Params:
		for _, n := range e.vars {
			result = append(result, findNodes[T](n)...)
		}
	case Var:
		result = append(result, findNodes[T](e.kind)...)
		result = append(result, findNodes[T](e.assigned)...)
	// Statements
	case StmtVar:
		result = append(result, findNodes[T](e.kind)...)
		result = append(result, findNodes[T](e.assigned)...)
	case StmtExpr:
		result = append(result, findNodes[T](e.expr)...)
	case StmtIf:
		result = append(result, findNodes[T](e.test)...)
		for _, n := range e.body {
			result = append(result, findNodes[T](n)...)
		}
		for _, n := range e.elifs {
			result = append(result, findNodes[T](n)...)
		}
		for _, n := range e.elseBody {
			result = append(result, findNodes[T](n)...)
		}
	case StmtElif:
		result = append(result, findNodes[T](e.test)...)
		for _, b := range e.body {
			result = append(result, findNodes[T](b)...)
		}
	case StmtWhile:
		result = append(result, findNodes[T](e.test)...)
		for _, n := range(e.body) {
			result = append(result, findNodes[T](n)...)
		}
	case StmtReturn:
		result = append(result, findNodes[T](e.value)...)
	// Expressions
	case ExprBinary:
		result = append(result, findNodes[T](e.left)...)
		result = append(result, findNodes[T](e.right)...)
	case ExprUnary:
		result = append(result, findNodes[T](e.value)...)
	case ExprCall:
		result = append(result, findNodes[T](e.callee)...)
		for _, arg := range e.args {
			result = append(result, findNodes[T](arg)...)
		}
	case ExprField:
		result = append(result, findNodes[T](e.object)...)
	case ExprObject:
		for _, field := range e.fields {
			result = append(result, findNodes[T](field)...)
		}
	}

	return result
}
