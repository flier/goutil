package inspect

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Inspector struct {
	strings.Builder
	options
	dump  func() func(any) string
	dump2 func() func(any, any) string
	i, w  int
}

func New(x []Option) *Inspector {
	var o options

	for _, opt := range x {
		opt(&o)
	}

	d := sync.OnceValue(func() (d func(any) string) {
		d = o.dump
		if d != nil {
			return
		}

		f := "%v"
		if o.format != "" {
			f = o.format
		} else if o.pretty {
			f = "%#v"
		}

		d = func(v any) string { return fmt.Sprintf(f, v) }

		return
	})

	d2 := sync.OnceValue(func() (d func(any, any) string) {
		d = o.dump2
		if d != nil {
			return
		}

		f := "%v:%v"
		if o.format != "" {
			f = o.format
		} else if o.pretty {
			f = "%#v:%#v"
		}

		d = func(k, v any) string { return fmt.Sprintf(f, k, v) }

		return
	})

	return &Inspector{options: o, dump: d, dump2: d2}
}

func (i *Inspector) Start() {
	if i.label != "" {
		i.WriteString(i.label)
		i.WriteString(": ")
	}

	i.WriteByte('[')
	i.w = 1
}

func (i *Inspector) Stop() {
	i.WriteByte(']')

	w := i.writer
	if w == nil {
		w = os.Stdout
	}

	_, _ = fmt.Fprintln(w, i.String())
}

func (i *Inspector) Inspect(v any) {
	i.inspect(func() string {
		d := i.dump()

		return d(v)
	})
}

func (i *Inspector) Inspect2(k, v any) {
	i.inspect(func() string {
		d := i.dump2()

		return d(k, v)
	})
}

func (i *Inspector) inspect(f func() string) {
	if i.limit == 0 || i.i < i.limit {
		s := f()

		if i.w += len(s); i.width != 0 && i.w > i.width {
			i.WriteByte('\n')
			i.w = len(s)
		}

		if i.i > 0 {
			i.WriteByte(' ')
			i.w += 1
		}

		i.WriteString(s)
	} else if i.i > 0 && i.i == i.limit {
		if i.width != 0 && i.w+4 > i.width {
			i.WriteByte('\n')
		}
		i.WriteString(" ...")
	}

	i.i += 1
}
