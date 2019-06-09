package main

import (
	"testing"
)

func addOne(a interface{}) interface{} {
	if val, ok := a.(int); ok {
		return val + 1
	} else {
		panic("Input should be interger")
	}
}

func getAddOnePromise(a interface{}) interface{} {
	return Resolve(addOne(a))
}

func subtractOne(a interface{}) interface{} {
	if val, ok := a.(int); ok {
		return val - 1
	} else {
		panic("Input should be interger")
	}
}

func getSubtractOnePromise(a interface{}) interface{} {
	return Resolve(subtractOne(a))
}

func assertValue(t *testing.T, testValue, expectedValue interface{}, msg string) {
	if testValue != expectedValue {
		t.Error(msg)
	}
}

func TestResolve(t *testing.T) {
	a := Resolve(7)
	out := <-a.ch
	if out != 7 {
		t.Error(`Resolve(7) does not resolve to 7`)
	}
}
func TestReject(t *testing.T) {
	a := Reject(7)
	out := <-a.err_ch
	if out != 7 {
		t.Error(`Reject(7) does not reject to 7`)
	}
}

func TestThenChain(t *testing.T) {
	a := Resolve(7).Then(addOne).Then(addOne)
	assertValue(t, <-a.ch, 9, "Then Failed")
}

func TestPromiseInsideThen(t *testing.T) {
	a := Resolve(7).Then(getAddOnePromise)
	assertValue(t, <-a.ch, 8, "TestPromiseInsideThen Failed")
}

func TestOnRejected(t *testing.T) {
	a := Reject(7).Then(getAddOnePromise, subtractOne)
	assertValue(t, <-a.ch, 6, "TestOnRejected Failed")
}

func TestOnRejectedWithPromise(t *testing.T) {
	a := Reject(7).Then(getAddOnePromise, getSubtractOnePromise)
	assertValue(t, <-a.ch, 6, "TestOnRejected Failed")
}

func TestPromiseInsideResolve(t *testing.T) {
	a := Resolve(getAddOnePromise(7)).Then(addOne)
	assertValue(t, <-a.ch, 9, "TestPromiseInsideResolve Failed")
}

func TestPanic(t *testing.T) {
	err_msg := "Error!"
	var getErrorMessage = func(msg interface{}) interface{} {
		return err_msg
	}
	a := Resolve("this will give error").Then(addOne).Then(addOne, getErrorMessage)
	assertValue(t, <-a.ch, err_msg, "TestPanic failed")
}

func TestCatch(t *testing.T) {
	err_msg := "Error!"
	var identityFunc = func(arg interface{}) interface{} { return arg }
	a := Reject(err_msg).Catch(identityFunc)
	assertValue(t, <-a.ch, err_msg, "TestCatch failed")
}
