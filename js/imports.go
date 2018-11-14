package js

func refRelease(p uintptr)

func newBool(b bool) uintptr

func newInt(i int) uintptr

func newFloat(f float32) uintptr

func newString() uintptr

func zero(ref uintptr, t Type)

func newArray() uintptr

func appendString(ref uintptr, r rune)

func intAt(ref uintptr, index int) int

func typeOf(ref uintptr) Type

func getInt(ref uintptr) int

func setInt(ref uintptr, i int)

func getFloat(ref uintptr) float32

func setFloat(ref uintptr, f float32)

func getRef(ref uintptr, key uintptr) uintptr

func getIndex(ref uintptr, index int) uintptr

func setIndex(ref uintptr, index int, val uintptr)

func pushRef(ref uintptr, v uintptr)

func instanceOf(r1, r2 uintptr) int

func charCodeAt(ref uintptr, index int) rune

func call(ref uintptr, key uintptr, args uintptr) uintptr

func invoke(ref uintptr, args uintptr) uintptr

func noo(ref uintptr, args uintptr) uintptr

func lengthOf(ref uintptr) int

func setVal(ref uintptr, key uintptr, val uintptr)
