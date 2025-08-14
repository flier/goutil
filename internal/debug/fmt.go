package debug

import (
	"fmt"
	"reflect"
	"runtime"
)

// Formatter is a fmt.Formatter implementation that just calls a function.
type Formatter func(s fmt.State)

func (f Formatter) Format(s fmt.State, verb rune) {
	if verb != 'v' {
		_, _ = fmt.Fprintf(s, "%%%c(%T=%v)", verb, f, Func(f))
		return
	}
	f(s)
}

func (f Formatter) String() string { return fmt.Sprint(f) }

// Fprintf is like Fprintf, but the printing is delayed until the returned value
// is formatted with %v.
func Fprintf(format string, args ...any) Formatter {
	return Formatter(func(s fmt.State) { _, _ = fmt.Fprintf(s, format, args...) })
}

// Func pretty-prints a function value.
func Func(f any) Formatter {
	return Formatter(func(s fmt.State) {
		v := reflect.ValueOf(f)

		var pc uintptr
		switch v.Kind() {
		case reflect.Func:
			pc = uintptr(v.UnsafePointer())
		case reflect.Uintptr:
			pc = uintptr(v.Uint())
		default:
			_, _ = fmt.Fprintf(s, "%%v(NONFUNC:%v)", v)
		}

		fn := runtime.FuncForPC(pc)
		name := fn.Name()
		if name == "" {
			name = "<unknown>"
		}

		_, _ = fmt.Fprintf(s, "%#x:%s", pc, name)
	})
}

// Dict pretty-prints the given entries as a dictionary, with an optional
// prefix.
func Dict(prefix any, kv ...any) Formatter {
	return Formatter(func(s fmt.State) {
		if len(kv)%2 != 0 {
			panic("dbg: length must be divisible by 2")
		}

		if prefix == nil {
			prefix = ""
		}

		first := true
		_, _ = fmt.Fprintf(s, "%v{", prefix)
		for i := 0; i < len(kv)/2; i++ {
			k := kv[2*i]
			v := kv[2*i+1]
			if v == nil {
				continue
			}

			if !first {
				_, _ = fmt.Fprint(s, ", ")
			}
			first = false
			_, _ = fmt.Fprintf(s, "%v: %v", k, v)
		}
		_, _ = fmt.Fprint(s, "}")
	})
}
