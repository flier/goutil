package inspect

import (
	"fmt"
	"os"
	"strings"
)

type Inspector struct {
	strings.Builder
	options
	dump  func(any) string
	dump2 func(any, any) string
	i, w  int
}

func New(x []Option) *Inspector {
	var o options

	for _, opt := range x {
		opt(&o)
	}

	d := o.dump
	if d == nil {
		f := o.format
		if f == "" {
			if o.pretty {
				f = "%#v"
			} else {
				f = "%v"
			}
		}

		d = func(v any) string {
			return fmt.Sprintf(f, v)
		}
	}

	d2 := o.dump2
	if d2 == nil {
		f := o.format
		if f == "" {
			if o.pretty {
				f = "%#v:%#v"
			} else {
				f = "%v:%v"
			}
		}

		d2 = func(k, v any) string { return fmt.Sprintf(f, k, v) }
	}

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
		return i.dump(v)
	})
}

func (i *Inspector) Inspect2(k, v any) {
	i.inspect(func() string {
		return i.dump2(k, v)
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
