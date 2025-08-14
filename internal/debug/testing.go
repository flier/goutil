package debug

import (
	"testing"

	"github.com/timandy/routine"
)

var tls = routine.NewThreadLocal[testing.TB]()

// WithTesting sets a testing pointer for debugging.
//
// This will cause t.Log() to be used to print debug traces instead of Debug.
func WithTesting(t testing.TB) func() {
	t.Helper()

	prev := tls.Get()
	tls.Set(t)
	return func() {
		tls.Set(prev)
	}
}
