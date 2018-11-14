package js

import "strings"

var b = new(strings.Builder)

func toString(ref uintptr) string {
	b.Reset()
	l := lengthOf(ref)
	for i := 0; i < l; i++ {
		b.WriteRune(charCodeAt(ref, i))
	}
	return b.String()
}

func getBool(ref uintptr) bool {
	switch getInt(ref) {
	case 1:
		return true
	}
	return false
}

func setBool(ref uintptr, val bool) {
	if val {
		setInt(ref, 1)
		return
	}
	setInt(ref, 0)
}
