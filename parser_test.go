package main

import (
	"reflect"
	"testing"
)

func TestAst(t *testing.T) {
	var tokens, err = getTokens([]byte(`
var global: str = "hello"

fun f(arg1: int = 10, arg2: str = "hello" ...) int {
	"standalone expression"
	var my_var: obj = obj {a_field: fun = f}
	if (1 and 0) {
		print("hello")
		print("hi" + "yo")
	} elif nil { # not truthy
		print("never prints")
		print("nothing either")
	} elif "truthy" {
		print("should print")
	} else {
		# Doesn't return
		return 10
	}

	while (!nil) {
	      break
	      break # never evaluates
	}

	my_var.a_field = nil

	
}

fun f2() {
}
`))

	if err != nil {
		t.Fatal(err)
	}

	var input = getAst(tokens)
	var expected = Ast{
		[]Func{{
			"f",
			Params{
				[]Var{
					{"arg1", TYPE_INT, ExprNumber{10}},
					{"arg2", TYPE_STR, ExprString{"hello"}}},
				true},
			TYPE_INT,
			[]Stmt{
				StmtExpr{ExprString{"standalone expression"}},
				StmtVar{"my_var", TYPE_OBJ, ExprObject{
					[]Var{
						{"a_field", TYPE_FUN, ExprSymbol{"f"}},
					},
				}},
				StmtIf{
					// Test
					ExprBinary{ExprNumber{1}, ExprNumber{0}, OP_AND},
					// Body
					[]Stmt{
						StmtExpr{
							ExprCall{
								ExprSymbol{"print"},
								[]Expr{ExprString{"hello"}},
							},
						},
						StmtExpr{
							ExprCall{
								ExprSymbol{"print"},
								[]Expr{
									ExprBinary{
										ExprString{"hi"},
										ExprString{"yo"},
										OP_ADD,
									},
								},
							},
						},
					},
					// Elifs
					[]StmtElif{
						{
							ExprNil{},
							[]Stmt{
								StmtExpr{
									ExprCall{
										ExprSymbol{"print"},
										[]Expr{ExprString{"never prints"}},
									},
								},
								StmtExpr{
									ExprCall{
										ExprSymbol{"print"},
										[]Expr{ExprString{"nothing either"}},
									},
								},
							},
						},
						{
							ExprString{"truthy"},
							[]Stmt{
								StmtExpr{
									ExprCall{
										ExprSymbol{"print"},
										[]Expr{ExprString{"should print"}},
									},
								},
							},
						},
					},
					// Else
					[]Stmt{StmtReturn{ExprNumber{10}}},
				},
				// While
				StmtWhile{
					ExprUnary{ExprNil{}, OP_NOT},
					[]Stmt{
						StmtBreak{},
						StmtBreak{},
					},
				},
				StmtExpr{
					ExprBinary{
						ExprField{ExprSymbol{"my_var"}, "a_field"},
						ExprNil{},
						OP_ASS,
					},
				},
			}},
			{"f2", Params{}, TYPE_VOID, nil},
		},
		[]Var{
			{"global", TYPE_STR, ExprString{"hello"}},
		},
	}

	if !reflect.DeepEqual(input, expected) {
		t.Fatal("Input did not match expected")
	}
}
