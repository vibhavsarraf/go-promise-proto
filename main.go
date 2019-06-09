package main

import (
	"fmt"
)

func check(args ...int) {
	fmt.Println(len(args))
}

func main() {
	fmt.Println("####################################")
	// p := Resolve(10).Then(addOne).Then(addOne)
	// fmt.Println(<-p.ch)
	check(10)
}
