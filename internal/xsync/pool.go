package xsync

import "sync"

// Pool is like sync.Pool, but strongly typed to make the interface a bit
// less messy.
type Pool[T any] struct {
	New   func() *T // Called to construct new values.
	Reset func(*T)  // Called to reset values before re-use.

	impl sync.Pool
}

// Get returns a cached value of type T.
//
//go:nosplit
func (p *Pool[T]) Get() *T {
	v, _ := p.impl.Get().(*T)
	if v == nil {
		switch p.New {
		case nil:
			v = new(T)
		default:
			v = p.New()
		}
	}
	return v
}

// Put returns a cached value of type T.
//
//go:nosplit
func (p *Pool[T]) Put(v *T) {
	if p.Reset != nil {
		p.Reset(v)
	}
	p.impl.Put(v)
}
