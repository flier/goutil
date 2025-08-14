package debug

import (
	"fmt"
	"runtime"
	"strings"
)

// Unsupported returns "unimplemented" error for the calling function.
func Unsupported() error {
	pc, _, _, _ := runtime.Caller(1)
	return &errUnsupported{pc}
}

// errUnsupported is the error returned by Unimplemented.
type errUnsupported struct{ pc uintptr }

func (e *errUnsupported) Error() string {
	name := runtime.FuncForPC(e.pc).Name()
	if name == "" {
		return "hyperpb: unsupported operation"
	}

	slash := strings.LastIndexByte(name, '/')
	name = name[slash+1:]
	return fmt.Sprintf("hyperpb: %s() is not supported", name)
}
