//go:build !debug

package debug

const Enabled = false

func Log([]any, string, string, ...any) {}
func Assert(bool, string, ...any)       {}

type Value[T any] struct {
	_ struct{}
}

func (v *Value[T]) Get() *T {
	panic("called Value.Get() when not in debug mode")
}
