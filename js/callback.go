package js

import "strconv"

type EventCallbackFlag int

const (
	// PreventDefault can be used with NewEventCallback to call event.preventDefault synchronously.
	PreventDefault EventCallbackFlag = 1 << iota
	// StopPropagation can be used with NewEventCallback to call event.stopPropagation synchronously.
	StopPropagation
	// StopImmediatePropagation can be used with NewEventCallback to call event.stopImmediatePropagation synchronously.
	StopImmediatePropagation
)

var counter = 0

var fm = make(map[int]func([]Value))

//go:export __calledback
func callbacks(f int, ptr int, size int) {
	fun := fm[f]
	if fun == nil {
		return
	}
	if ptr == 0 {
		fun(nil)
		return
	}
	args := make([]Value, size)
	for i := range args {
		args[i].ref = uintptr(ptr + i)
	}
	fun(args)
	for _, v := range args {
		refRelease(v.ref)
	}
}

type Callback struct {
	Value
	id int
}

func NewCallback(f func(args []Value)) Callback {
	jsf := Global().Get("Function")
	defer jsf.release()
	cb := jsf.New(
		`var id = ` + strconv.Itoa(counter) + `;
		if (arguments.length == 0) {
			wasm.exports.__calledback(id,0,0);
			return;	
		}
		rc++
		ptr = rc;
		for (var i = 0 ; i < arguments.length ; i++) {
			refs.set(rc,arguments[i])
			rc++;
		}
		wasm.exports.__calledback(id,ptr,arguments.length);
	`)
	fm[counter] = f
	counter++
	return Callback{Value: cb, id: counter - 1}
}

func (c Callback) Release() {
	c.Value.release()
}
