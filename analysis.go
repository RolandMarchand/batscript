package main

import (
	"errors"
	"fmt"
)

// Rules:
// 1. Inside of functions, symbols need to be defined either before its declaration or in global scope
// 2. The types of parameters and variables must be the same as the type of the assigned expression
// 3. Calling a global function (not a fun variable) must respect arg count and arg types
// 4. Only variables of type obj can have fields accessed
// 5. break is only used in while loops
// 6. return must return an expression of the same type as the current function, and must return nothing if the function returns nothing
// 7. No local variable can share a name with an earlier variable, a parameter, a function, or a global variable
// 8. No parameter can share a name with an earlier parameter, a function, or global variable
// 9. No global variable can share a name with an earlier function or global variable
// 10. No function can share a name with an earlier function or global variable
// 11. Global variables cannot be initialized with global variables declared later
// 12. A function with a non-void return type must return a value on all paths
// 13. All binary operators can only be done with members of the same type, but and and or do not have this restriction, and != and == work with fun/nil and obj/nil
// 14. Binary + is only allowed with strings and integers and symbols resolving to integers
// 15. Binary - * / % <= >= < > are only allowed with integers and symbols resolving to integers
// 16. Unary + and - are only allowed with integers and symbols resolving to integers
// 17. for =, on lvalues can only be field access or symbols
// 18. for =, if the lvalue is a symbol, the type of the rvalue must be the same as the variable of the lvalue
// 19. for =, the lvalue cannot be the symbol of a function
// 20. Initial field accesses can only be on object type variables
// 21. Two fields on the same scope in an object cannot have the same name
// 22. When an lvalue is called and it is not a field access, it must be a symbol resolving to a function
// 23. There must be one main function with no argument and no return type


type varTable []map[string]Type
type funcSig struct {
	params []Type
	ret Type
}
type funcTable map[string]funcSig

type symbolTable struct {
	vars varTable
	funcs funcTable
}

func analyze(ast Ast) {
	var t symbolTable
	var t
	var err = registerGlobalScope(ast, &t)
	if err != nil {
		panic(err)
	}

	analyzeGlobalVars(ast.vars, &t)
}

func registerGlobalScope(ast Ast, t *symbolTable, f *funcTable) (errs error) {
	var scope = make(map[string]Type, len(ast.funcs) + len(ast.vars))

	for _, v := range ast.vars {
		if _, exists := scope[v.name]; exists {
			errs = errors.Join(errs, fmt.Errorf(
				"Variable %s is already defined",
				v.name,
			))
			continue
		}

		scope[v.name] = v.kind
	}

	for _, f := range ast.funcs {
		if _, exists := scope[f.name]; exists {
			errs = errors.Join(errs, fmt.Errorf(
				"Function %s is already defined",
				f.name,
			))
			continue
		}

		scope[f.name] = TYPE_FUN
	}

	*t = append(*t, scope)

	return errs
}

func analyzeGlobalVars(vars []Var, t *symbolTable) error {
	var scope = make(map[string]Type, len(vars))

	for _, v := range vars {
		var exprType, err = getExprType(&v.assigned, *t)
		if err != nil {
			return err
		}

		if exprType != v.kind {
			return fmt.Errorf(
				"Variable %s is type %s but was assigned %s",
				v.name,
				v.kind,
				exprType,
			)
		}
	}

	*t = append(*t, scope)

	return nil
}

func getExprType(expr *Expr, t symbolTable) (Type, error) {
	return TYPE_VOID, nil
}
