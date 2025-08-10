package tuple_test

import (
	"fmt"
	"io"

	. "github.com/flier/goutil/pkg/tuple"
)

func ExampleNew2() {
	t := New2("hello", 42)

	fmt.Println(t)
	fmt.Println(t.Unpack())

	// Output:
	// (hello, 42)
	// hello 42
}

func ExampleNew3() {
	t := New3("hello", 42, 3.14)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())

	// Output:
	// (hello, 42, 3.14)
	// hello 42 3.14
	// hello (42, 3.14)
	// (hello, 42) 3.14
}

func ExampleNew4() {
	t := New4("hello", 42, 3.14, io.EOF)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())

	// Output:
	// (hello, 42, 3.14, EOF)
	// hello 42 3.14 EOF
	// hello (42, 3.14, EOF)
	// (hello, 42, 3.14) EOF
}

func ExampleNew5() {
	t := New5("hello", 42, 3.14, io.EOF, true)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())

	// Output:
	// (hello, 42, 3.14, EOF, true)
	// hello 42 3.14 EOF true
	// hello (42, 3.14, EOF, true)
	// (hello, 42, 3.14, EOF) true
}

func ExampleNew6() {
	t := New6("hello", 42, 3.14, io.EOF, true, 'c')

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())

	// Output:
	// (hello, 42, 3.14, EOF, true, 99)
	// hello 42 3.14 EOF true 99
	// hello (42, 3.14, EOF, true, 99)
	// (hello, 42, 3.14, EOF, true) 99
}

func ExampleNew7() {
	t := New7("hello", 42, 3.14, io.EOF, true, 'c', []string{"foo", "bar"})

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())

	// Output:
	// (hello, 42, 3.14, EOF, true, 99, [foo bar])
	// hello 42 3.14 EOF true 99 [foo bar]
	// hello (42, 3.14, EOF, true, 99, [foo bar])
	// (hello, 42, 3.14, EOF, true, 99) [foo bar]
}
