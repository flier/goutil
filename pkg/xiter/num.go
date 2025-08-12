package xiter

// Number is an interface that represents any numeric type in Go, including integers,
// unsigned integers, and floating-point numbers.
type Number interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~int | ~uint | ~uintptr | ~float32 | ~float64
}

// Integer is an interface that represents any integer type in Go, including signed and unsigned integers of various bit widths.
type Integer interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~int | ~uint | ~uintptr
}

// Signed is an interface that represents any signed integer type in Go, including signed integers of various bit widths.
type Signed interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// Unsigned is an interface that represents any unsigned integer type in Go, including unsigned integers of various bit widths.
type Unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr
}
