package main

import (
	"fmt"
	"testing"
)

var fmtWrapper = func(x interface{}) interface{} {
	fmt.Println(x)
	return x
}

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

func errFunc(a interface{}) interface{} {
	panic("Error Function called")
	return 0
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

func TestFinally(t *testing.T) {
	a := Resolve(7).Then(errFunc).Finally(
		func(x interface{}) interface{} { return 8 },
	)
	assertValue(t, <-a.ch, 8, "TestFinally failed")
}

func TestMultipleTypes(t *testing.T) {
	var getDay = func(x interface{}) interface{} {
		if val, ok := x.(int); ok {
			switch val {
			case 1:
				return "Monday"
			case 2:
				return "Tuesday"
			case 3:
				return "Wednesday"
			case 4:
				return "Thrusday"
			case 5:
				return "Friday"
			case 6:
				return "Saturday"
			case 7:
				return "Sunday"
			}
		}
		panic("input to getDay should be an integet between 1 and 7")
	}
	var getLength = func(x interface{}) interface{} {
		if val, ok := x.([]interface{}); ok {
			return len(val)
		}
		if val, ok := x.(string); ok {
			return len(val)
		}
		panic("input to getLength should be of type []interface{} or string")
		return 0
	}

	a := Resolve(7).Then(getDay).Then(getLength)
	assertValue(t, <-a.ch, 6, "TestMultipleTypes failed")
	arr := []interface{}{0, 1, 2, 3}
	a = Resolve(arr).Then(getLength)
	assertValue(t, <-a.ch, 4, "TestMultipleTypes failed")
}
