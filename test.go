//build +tinygo

package main

import (
	"github.com/j7b/syscall/js"
)

type value struct {
	ref uintptr
}

func round(f float32) int {
	return int(f + 0.5)
}

var cb js.Callback

func main() {
	g := js.Global()
	a := g.Get("Array")
	a1 := a.New(1, 2, 3)
	a2 := a.New(3, 2, 1)
	if !a1.InstanceOf(a) {
		println("want a1 to be array")
	}
	if !a2.InstanceOf(a) {
		println("want a2 to be array")
	}
	b := js.ValueOf(true)
	if b.Bool() != true {
		println("wanted b true")
	}
	lio := a2.Call("lastIndexOf", 2)
	if lio.Int() != 1 {
		println("want lio 1 got", lio.Int())
	}
	str := g.Get("String")
	s1 := str.New("abc")
	if s1.String() != "abc" {
		println("wanted abc, got", s1.String())
	}
	number := g.Get("Number")
	n1 := number.New(3.9)
	rounded := round(n1.Float())
	if rounded != 4 {
		println(`wanted 4`)
	}
	g.Set("hello", "sailor")
	if s := g.Get("hello"); s.String() != "sailor" {
		println("want sailor")
	}
	three := a1.Index(2)
	if three.Int() != 3 {
		println("wanted 3")
	}
	function := g.Get("Function")
	f1 := function.New(`a`, `b`, `return a+b;`)
	c := f1.Invoke(3, 4)
	if c.Int() != 7 {
		println("wanted 7")
	}
	if a1.Length() != 3 {
		println("want a1 to be 3")
	}
	a1.SetIndex(1, 43)
	if a1.Index(1).Int() != 43 {
		println("wanted 43")
	}
	raf := g.Get("requestAnimationFrame")
	var then float32
	cb = js.NewCallback(func(vals []js.Value) {
		now := vals[0].Float()
		then = now
		_ = then
		// raf.Invoke(cb)
	})
	raf.Invoke(cb.Value)
	println("done")
}
