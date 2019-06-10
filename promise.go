package main

type Promise struct {
	ch     chan interface{}
	err_ch chan interface{}
}

type PromiseFunc func(arg interface{}) interface{}

func newPromise() Promise {
	var p Promise
	p.ch = make(chan interface{})
	p.err_ch = make(chan interface{})
	return p
}

func passChannel(new_ch chan interface{}, old_ch chan interface{}) {
	go func() { new_ch <- (<-old_ch) }()
}

func passChannelValue(ch chan interface{}, val interface{}) {
	go func() { ch <- val }()
}

func passChannels(p_new Promise, p_old Promise) {
	passChannel(p_new.ch, p_old.ch)
	passChannel(p_new.err_ch, p_old.err_ch)
}

//Resolve returns a promise that resolves to the give input value
func Resolve(arg interface{}) Promise {
	if _, ok := arg.(Promise); ok {
		return arg.(Promise)
	}
	var p = newPromise()
	passChannelValue(p.ch, arg)
	return p
}

func Reject(arg interface{}) Promise {
	var p = newPromise()
	if _, ok := arg.(Promise); ok {
		passChannel(p.err_ch, arg.(Promise).ch)
	}
	passChannelValue(p.err_ch, arg)
	return p
}

func handlePromise(p Promise, val interface{}) {
	if _, ok := val.(Promise); ok {
		passChannels(p, val.(Promise))
	} else {
		passChannelValue(p.ch, val)
	}
}

func apply(p, p_new Promise, args []PromiseFunc) {
	defer func() {
		if r := recover(); r != nil {
			passChannelValue(p_new.err_ch, r)
		}
	}()
	select {
	case rejectValue := <-p.err_ch:
		if len(args) < 2 {
			passChannelValue(p_new.err_ch, rejectValue)
			return
		}
		handlePromise(p_new, args[1](rejectValue))
	case resolveVlaue := <-p.ch:
		handlePromise(p_new, args[0](resolveVlaue))
	}
}

func (p Promise) Then(args ...PromiseFunc) Promise {
	var p_new = newPromise()
	go apply(p, p_new, args)
	return p_new
}

func (p Promise) Catch(onRejected PromiseFunc) Promise {
	var identityFunc = func(arg interface{}) interface{} { return arg }
	return p.Then(identityFunc, onRejected)
}

func (p Promise) Finally(onFinally PromiseFunc) Promise {
	return p.Then(onFinally, onFinally)
}
