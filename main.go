package main

import (
	"fmt"
)

func addOne(a interface{}) int {
	if val, ok := a.(int); ok {
		return val + 1
	} else {
		return 0
	}
}

func main() {
	fmt.Println("####################################")
	p := Resolve(10).Then(addOne).Then(addOne)
	fmt.Println(<-p.ch)
}
