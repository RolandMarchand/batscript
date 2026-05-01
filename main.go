package main

import ()

func main() {
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
		panic(err)
	}

	var ast = getAst(tokens)
	printAst(ast)
}
