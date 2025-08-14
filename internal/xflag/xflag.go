//go:build go1.21

package xflag

import (
	"flag"
	"sync"
)

var parsed = sync.OnceValue(func() map[string]struct{} {
	m := make(map[string]struct{})
	flag.Visit(func(f *flag.Flag) { m[f.Name] = struct{}{} })
	return m
})

// Func is like [flags.Func], but avoids the need for an init func by allocating
// its own storage for the return value.
func Func[T any](name, usage string, fn func(string) (T, error)) *T {
	v := new(T)
	flag.Func(name, usage, func(s string) (err error) {
		*v, err = fn(s)
		return err
	})
	return v
}

// Lookup looks up a flag by name of the given type.
//
// Panics if this flag is of the wrong type, or if the flag value is not a
// [flag.Getter].
func Lookup[T any](name string) T {
	return flag.Lookup(name).Value.(flag.Getter).Get().(T) //nolint:errcheck
}

// Parsed returns whether the given flag was parsed.
func Parsed(name string) bool {
	if !flag.Parsed() {
		return false
	}
	_, ok := parsed()[name]
	return ok
}
