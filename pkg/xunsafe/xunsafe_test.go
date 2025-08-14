package xunsafe_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestIndirect(t *testing.T) {
	t.Parallel()

	assert.False(t, xunsafe.IsDirect[int]())
	assert.False(t, xunsafe.IsDirect[string]())
	assert.False(t, xunsafe.IsDirect[[]byte]())

	assert.True(t, xunsafe.IsDirect[*int]())
	assert.True(t, xunsafe.IsDirect[[1]*int]())
	assert.True(t, xunsafe.IsDirect[any]())
	assert.True(t, xunsafe.IsDirect[map[int]int]())
	assert.True(t, xunsafe.IsDirect[chan int]())
	assert.True(t, xunsafe.IsDirect[unsafe.Pointer]())
	assert.True(t, xunsafe.IsDirect[struct{ _ *int }]())
	assert.True(t, xunsafe.IsDirect[*struct{ _ *int }]())
}

func TestAnyBytes(t *testing.T) {
	t.Parallel()

	i := 0xaaaa
	p := &i
	assert.False(t, xunsafe.IsDirectAny(i))
	assert.True(t, xunsafe.IsDirectAny(p))

	assert.Equal(t, xunsafe.Bytes(&i), xunsafe.AnyBytes(i))
	assert.Equal(t, xunsafe.Bytes(&p), xunsafe.AnyBytes(p))

	p2 := struct{ p *int }{p}
	assert.Equal(t, xunsafe.Bytes(&p2), xunsafe.AnyBytes(p2))
}

func TestPC(t *testing.T) {
	t.Parallel()

	f := func() int { return 42 }
	pc := xunsafe.NewPC(f)

	t.Logf("%#x\n", pc)
	assert.Equal(t, 42, pc.Get()())
}
