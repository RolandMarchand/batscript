package main

import "fmt"

type i interface {
	lmao()
}

type s struct {
	s string
}

func (S s) lmao() {
}

func Main() {
	var l i = s{"hey"}
	penis(l)
	fmt.Println(l)
}

func penis(I i) {
	var S = I.(s)
	S.s = "penis"
}
