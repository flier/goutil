package inspect

type options struct {
	dump   func(any) string
	dump2  func(any, any) string
	format string
	label  string
	limit  int
	pretty bool
	width  int
}

type Option func(*options)

var (
	// Pretty enables pretty printing
	Pretty Option = func(o *options) { o.pretty = true }
)

// Format give a format string instead of "%v" by default
func Format(f string) Option {
	return func(o *options) { o.format = f }
}

// Label decorate the output
func Label(l string) Option {
	return func(o *options) { o.label = l }
}

// Limits the number of items that are inspected
func Limit(n int) Option {
	return func(o *options) {
		o.limit = n
	}
}

// Width defines the number of characters per line used when pretty is true
func Width(w int) Option {
	return func(o *options) { o.width = w }
}

// Dump give a function to dump the value.
func Dump(d func(any) string) Option {
	return func(o *options) { o.dump = d }
}

// Dump2 give a function to dump the key-value.
func Dump2(d func(any, any) string) Option {
	return func(o *options) { o.dump2 = d }
}
