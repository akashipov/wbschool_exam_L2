package main

import (
	"fmt"
)

type A struct {
	a string
}

func (a *A) Print() {
	fmt.Println("Hello")
}

type B interface {
	Print()
}

func main() {
	var i [3]int
	var b [3]int
	fmt.Println(i == b)
}
