package main

type Promise struct{ ch chan int }

type PromiseFunc func(arg interface{}) int

func newPromise() Promise {
	var p Promise
	p.ch = make(chan int)
	return p
}

//Resolve returns a promise that resolves to the give input value
func Resolve(arg int) Promise {
	var p = newPromise()
	go func() { p.ch <- arg }()
	return p
}

func apply(p, p_new Promise, args []PromiseFunc) {
	output := <-p.ch
	val := args[0](output)
	p_new.ch <- val
}

func (p Promise) Then(args ...PromiseFunc) Promise {
	var p_new = newPromise()
	go apply(p, p_new, args)
	return p_new
}
