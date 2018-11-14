package js

type Error struct {
	Value
}

func (e Error) Error() string {
	return ""
}

type Type int

const (
	TypeUndefined Type = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeSymbol
	TypeObject
	TypeFunction
)

func Release(v Value) {
	v.release()
}

func (t Type) String() string {
	return ""
}

type TypedArray struct {
	Value
}

func (t TypedArray) Release() {
	t.Value.release()
}

func TypedArrayOf(slice interface{}) TypedArray {
	return TypedArray{}
}

type Value struct {
	ref uintptr
}

func (v Value) release() {
	refRelease(v.ref)
}

func Global() Value {
	return Value{ref: uintptr(8)}
}

func Null() Value {
	return Value{ref: uintptr(TypeNull)}
}

func Undefined() Value {
	return Value{ref: uintptr(TypeUndefined)}
}

func ValueOf(v interface{}) Value {
	switch t := v.(type) {
	case nil:
		return Null()
	case TypedArray:
		return t.Value
	case Callback:
		return t.Value
	case Value:
		return t
	case *Value:
		return Value{ref: t.ref}
	case bool:
		return Value{ref: newBool(t)}
	case int:
		return Value{ref: newInt(t)}
	case float32:
		return Value{ref: newFloat(t)}
	case float64:
		f := float32(t)
		return Value{ref: newFloat(f)}
	case string:
		ref := newString()
		for _, r := range t {
			appendString(ref, r)
		}
		return Value{ref: ref}
	case []interface{}:
		ref := newArray()
		for _, v := range t {
			val := ValueOf(v)
			pushRef(ref, val.ref)
		}
		return Value{ref: ref}
	}
	panic("unsupported type")
}

func (v Value) Bool() bool {
	if typeOf(v.ref) == TypeBoolean {
		return getBool(v.ref)
	}
	panic("not a bool")
}

func (v Value) Call(key string, args ...interface{}) Value {
	kv := ValueOf(key)
	defer kv.release()
	av := ValueOf(args)
	defer av.release()
	return Value{ref: call(v.ref, kv.ref, av.ref)}
}

func (v Value) Float() float32 {
	return getFloat(v.ref)
}

func (v Value) Get(key string) Value {
	kv := ValueOf(key)
	defer kv.release()
	return Value{ref: getRef(v.ref, kv.ref)}
}

func (v Value) Index(index int) Value {
	return Value{ref: getIndex(v.ref, index)}
}

func (v Value) InstanceOf(o Value) bool {
	if instanceOf(v.ref, o.ref) == 1 {
		return true
	}
	return false
}

func (v Value) Int() int {
	return getInt(v.ref)
}

func (v Value) Invoke(args ...interface{}) Value {
	av := ValueOf(args)
	defer av.release()
	return Value{ref: invoke(v.ref, av.ref)}
}

func (v Value) Length() int {
	return lengthOf(v.ref)
}

func (v Value) New(args ...interface{}) Value {
	av := ValueOf(args)
	defer av.release()
	return Value{ref: noo(v.ref, av.ref)}
}

func (v Value) Set(p string, i interface{}) {
	key := ValueOf(p)
	defer key.release()
	val, ok := i.(Value)
	if !ok {
		val = ValueOf(i)
		defer val.release()
	}
	setVal(v.ref, key.ref, val.ref)
}

func (v Value) SetIndex(index int, i interface{}) {
	val, ok := i.(Value)
	if !ok {
		val = ValueOf(i)
		defer val.release()
	}
	setIndex(v.ref, index, val.ref)
}

func (v Value) String() string {
	r := v.Call("toString")
	defer r.release()
	return toString(r.ref)
}

func (v Value) Type() Type {
	return typeOf(v.ref)
}
