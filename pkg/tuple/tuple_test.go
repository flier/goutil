package tuple_test

import (
	"fmt"
	"io"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/tuple"
)

func ExampleNew0() {
	t := New0()

	fmt.Println(t)
	fmt.Println(t.Len())

	// Output:
	// ()
	// 0
}

func ExampleNew1() {
	t := New1("hello")

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(0))
	fmt.Println(t.Put(0, "foobar"))
	fmt.Println(t.Del(0))

	// Output:
	// (hello)
	// hello
	// hello ()
	// () hello
	// 1 hello
	// (foobar) hello
	// ()
}

func ExampleNew2() {
	t := New2("hello", 42)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(1))
	fmt.Println(t.Put(1, 123))
	fmt.Println(t.Del(1))

	// Output:
	// (hello, 42)
	// hello 42
	// hello (42)
	// (hello) 42
	// 2 42
	// (hello, 123) 42
	// (hello)
}

func ExampleNew3() {
	t := New3("hello", 42, 3.14)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(2))
	fmt.Println(t.Put(2, math.Pi))
	fmt.Println(t.Del(2))

	// Output:
	// (hello, 42, 3.14)
	// hello 42 3.14
	// hello (42, 3.14)
	// (hello, 42) 3.14
	// 3 3.14
	// (hello, 42, 3.141592653589793) 3.14
	// (hello, 42)
}

func ExampleNew4() {
	t := New4("hello", 42, 3.14, io.EOF)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(3))
	fmt.Println(t.Put(3, io.ErrUnexpectedEOF))
	fmt.Println(t.Del(3))

	// Output:
	// (hello, 42, 3.14, EOF)
	// hello 42 3.14 EOF
	// hello (42, 3.14, EOF)
	// (hello, 42, 3.14) EOF
	// 4 EOF
	// (hello, 42, 3.14, unexpected EOF) EOF
	// (hello, 42, 3.14)
}

func ExampleNew5() {
	t := New5("hello", 42, 3.14, io.EOF, true)

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(4))
	fmt.Println(t.Put(4, false))
	fmt.Println(t.Del(4))

	// Output:
	// (hello, 42, 3.14, EOF, true)
	// hello 42 3.14 EOF true
	// hello (42, 3.14, EOF, true)
	// (hello, 42, 3.14, EOF) true
	// 5 true
	// (hello, 42, 3.14, EOF, false) true
	// (hello, 42, 3.14, EOF)
}

func ExampleNew6() {
	t := New6("hello", 42, 3.14, io.EOF, true, 'c')

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(5))
	fmt.Println(t.Put(5, 'b'))
	fmt.Println(t.Del(5))

	// Output:
	// (hello, 42, 3.14, EOF, true, 99)
	// hello 42 3.14 EOF true 99
	// hello (42, 3.14, EOF, true, 99)
	// (hello, 42, 3.14, EOF, true) 99
	// 6 99
	// (hello, 42, 3.14, EOF, true, 98) 99
	// (hello, 42, 3.14, EOF, true)
}

func ExampleNew7() {
	t := New7("hello", 42, 3.14, io.EOF, true, 'c', []string{"foo", "bar"})

	fmt.Println(t)
	fmt.Println(t.Unpack())
	fmt.Println(t.Head())
	fmt.Println(t.Tail())
	fmt.Println(t.Len(), t.Get(6))
	fmt.Println(t.Put(6, []string{"foobar"}))
	fmt.Println(t.Del(6))

	// Output:
	// (hello, 42, 3.14, EOF, true, 99, [foo bar])
	// hello 42 3.14 EOF true 99 [foo bar]
	// hello (42, 3.14, EOF, true, 99, [foo bar])
	// (hello, 42, 3.14, EOF, true, 99) [foo bar]
	// 7 [foo bar]
	// (hello, 42, 3.14, EOF, true, 99, [foobar]) [foo bar]
	// (hello, 42, 3.14, EOF, true, 99)
}

func TestTuple(t *testing.T) {
	Convey("Given some tuples", t, func() {
		Convey("When create Tuple0", func() {
			t := New0()

			So(t.String(), ShouldEqual, "()")
			So(t.Len(), ShouldEqual, 0)
			So(func() { t.Get(0) }, ShouldPanicWith, ErrOutOfRange)
			So(func() { t.Put(0, 123) }, ShouldPanicWith, ErrOutOfRange)
			So(func() { t.Del(0) }, ShouldPanicWith, ErrOutOfRange)
		})

		Convey("When create Tuple1", func() {
			t := New1("hello")

			So(t.String(), ShouldEqual, "(hello)")
			So(t.Len(), ShouldEqual, 1)

			Convey("Then unpack the tuple", func() {
				So(t.Unpack(), ShouldEqual, "hello")
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar)")
					So(old, ShouldEqual, "hello")
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New0())
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})

		Convey("When create Tuple2", func() {
			t := New2("hello", 42)

			So(t.String(), ShouldEqual, "(hello, 42)")
			So(t.Len(), ShouldEqual, 2)

			Convey("Then unpack the tuple", func() {
				v0, v1 := t.Unpack()
				So(v0, ShouldEqual, "hello")
				So(v1, ShouldEqual, 42)
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
					So(t.Get(1), ShouldEqual, 42)
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar, 42)")
					So(old, ShouldEqual, "hello")
				})

				Convey("at index 1", func() {
					new, old := t.Put(1, 123)
					So(new.String(), ShouldEqual, "(hello, 123)")
					So(old, ShouldEqual, 42)
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
					So(func() { t.Put(1, false) }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New1(42))
				})

				Convey("at index 1", func() {
					So(t.Del(1), ShouldEqual, New1("hello"))
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})

		Convey("When create Tuple3", func() {
			t := New3("hello", 42, 3.14)

			So(t.String(), ShouldEqual, "(hello, 42, 3.14)")
			So(t.Len(), ShouldEqual, 3)

			Convey("Then unpack the tuple", func() {
				v0, v1, v3 := t.Unpack()
				So(v0, ShouldEqual, "hello")
				So(v1, ShouldEqual, 42)
				So(v3, ShouldEqual, 3.14)
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
					So(t.Get(1), ShouldEqual, 42)
					So(t.Get(2), ShouldEqual, 3.14)
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar, 42, 3.14)")
					So(old, ShouldEqual, "hello")
				})

				Convey("at index 1", func() {
					new, old := t.Put(1, 123)
					So(new.String(), ShouldEqual, "(hello, 123, 3.14)")
					So(old, ShouldEqual, 42)
				})

				Convey("at index 2", func() {
					new, old := t.Put(2, 1.23)
					So(new.String(), ShouldEqual, "(hello, 42, 1.23)")
					So(old, ShouldEqual, 3.14)
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
					So(func() { t.Put(1, false) }, ShouldPanic)
					So(func() { t.Put(2, 'c') }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New2(42, 3.14))
				})

				Convey("at index 1", func() {
					So(t.Del(1), ShouldEqual, New2("hello", 3.14))
				})

				Convey("at index2", func() {
					So(t.Del(2), ShouldEqual, New2("hello", 42))
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})

		Convey("When create Tuple4", func() {
			t := New4("hello", 42, 3.14, io.EOF)

			So(t.String(), ShouldEqual, "(hello, 42, 3.14, EOF)")
			So(t.Len(), ShouldEqual, 4)

			Convey("Then unpack the tuple", func() {
				v0, v1, v3, v4 := t.Unpack()
				So(v0, ShouldEqual, "hello")
				So(v1, ShouldEqual, 42)
				So(v3, ShouldEqual, 3.14)
				So(v4, ShouldEqual, io.EOF)
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
					So(t.Get(1), ShouldEqual, 42)
					So(t.Get(2), ShouldEqual, 3.14)
					So(t.Get(3), ShouldEqual, io.EOF)
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar, 42, 3.14, EOF)")
					So(old, ShouldEqual, "hello")
				})

				Convey("at index 1", func() {
					new, old := t.Put(1, 123)
					So(new.String(), ShouldEqual, "(hello, 123, 3.14, EOF)")
					So(old, ShouldEqual, 42)
				})

				Convey("at index 2", func() {
					new, old := t.Put(2, 1.23)
					So(new.String(), ShouldEqual, "(hello, 42, 1.23, EOF)")
					So(old, ShouldEqual, 3.14)
				})

				Convey("at index 3", func() {
					new, old := t.Put(3, io.ErrUnexpectedEOF)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, unexpected EOF)")
					So(old, ShouldEqual, io.EOF)
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
					So(func() { t.Put(1, false) }, ShouldPanic)
					So(func() { t.Put(2, 'c') }, ShouldPanic)
					So(func() { t.Put(3, 3.14) }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New3(42, 3.14, io.EOF))
				})

				Convey("at index 1", func() {
					So(t.Del(1), ShouldEqual, New3("hello", 3.14, io.EOF))
				})

				Convey("at index2", func() {
					So(t.Del(2), ShouldEqual, New3("hello", 42, io.EOF))
				})

				Convey("at index3", func() {
					So(t.Del(3), ShouldEqual, New3("hello", 42, 3.14))
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})

		Convey("When create Tuple5", func() {
			t := New5("hello", 42, 3.14, io.EOF, true)

			So(t.String(), ShouldEqual, "(hello, 42, 3.14, EOF, true)")
			So(t.Len(), ShouldEqual, 5)

			Convey("Then unpack the tuple", func() {
				v0, v1, v3, v4, v5 := t.Unpack()
				So(v0, ShouldEqual, "hello")
				So(v1, ShouldEqual, 42)
				So(v3, ShouldEqual, 3.14)
				So(v4, ShouldEqual, io.EOF)
				So(v5, ShouldBeTrue)
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
					So(t.Get(1), ShouldEqual, 42)
					So(t.Get(2), ShouldEqual, 3.14)
					So(t.Get(3), ShouldEqual, io.EOF)
					So(t.Get(4), ShouldEqual, true)
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar, 42, 3.14, EOF, true)")
					So(old, ShouldEqual, "hello")
				})

				Convey("at index 1", func() {
					new, old := t.Put(1, 123)
					So(new.String(), ShouldEqual, "(hello, 123, 3.14, EOF, true)")
					So(old, ShouldEqual, 42)
				})

				Convey("at index 2", func() {
					new, old := t.Put(2, 1.23)
					So(new.String(), ShouldEqual, "(hello, 42, 1.23, EOF, true)")
					So(old, ShouldEqual, 3.14)
				})

				Convey("at index 3", func() {
					new, old := t.Put(3, io.ErrUnexpectedEOF)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, unexpected EOF, true)")
					So(old, ShouldEqual, io.EOF)
				})

				Convey("at index 4", func() {
					new, old := t.Put(4, false)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, EOF, false)")
					So(old, ShouldBeTrue)
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
					So(func() { t.Put(1, false) }, ShouldPanic)
					So(func() { t.Put(2, 'c') }, ShouldPanic)
					So(func() { t.Put(3, 3.14) }, ShouldPanic)
					So(func() { t.Put(4, 'c') }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New4(42, 3.14, io.EOF, true))
				})

				Convey("at index 1", func() {
					So(t.Del(1), ShouldEqual, New4("hello", 3.14, io.EOF, true))
				})

				Convey("at index2", func() {
					So(t.Del(2), ShouldEqual, New4("hello", 42, io.EOF, true))
				})

				Convey("at index3", func() {
					So(t.Del(3), ShouldEqual, New4("hello", 42, 3.14, true))
				})

				Convey("at index4", func() {
					So(t.Del(4), ShouldEqual, New4("hello", 42, 3.14, io.EOF))
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})

		Convey("When create Tuple6", func() {
			t := New6("hello", 42, 3.14, io.EOF, true, 'c')

			So(t.String(), ShouldEqual, "(hello, 42, 3.14, EOF, true, 99)")
			So(t.Len(), ShouldEqual, 6)

			Convey("Then unpack the tuple", func() {
				v0, v1, v3, v4, v5, v6 := t.Unpack()
				So(v0, ShouldEqual, "hello")
				So(v1, ShouldEqual, 42)
				So(v3, ShouldEqual, 3.14)
				So(v4, ShouldEqual, io.EOF)
				So(v5, ShouldBeTrue)
				So(v6, ShouldAlmostEqual, 'c')
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
					So(t.Get(1), ShouldEqual, 42)
					So(t.Get(2), ShouldEqual, 3.14)
					So(t.Get(3), ShouldEqual, io.EOF)
					So(t.Get(4), ShouldEqual, true)
					So(t.Get(5), ShouldEqual, 'c')
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar, 42, 3.14, EOF, true, 99)")
					So(old, ShouldEqual, "hello")
				})

				Convey("at index 1", func() {
					new, old := t.Put(1, 123)
					So(new.String(), ShouldEqual, "(hello, 123, 3.14, EOF, true, 99)")
					So(old, ShouldEqual, 42)
				})

				Convey("at index 2", func() {
					new, old := t.Put(2, 1.23)
					So(new.String(), ShouldEqual, "(hello, 42, 1.23, EOF, true, 99)")
					So(old, ShouldEqual, 3.14)
				})

				Convey("at index 3", func() {
					new, old := t.Put(3, io.ErrUnexpectedEOF)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, unexpected EOF, true, 99)")
					So(old, ShouldEqual, io.EOF)
				})

				Convey("at index 4", func() {
					new, old := t.Put(4, false)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, EOF, false, 99)")
					So(old, ShouldBeTrue)
				})

				Convey("at index 5", func() {
					new, old := t.Put(5, 'b')
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, EOF, true, 98)")
					So(old, ShouldEqual, 'c')
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
					So(func() { t.Put(1, false) }, ShouldPanic)
					So(func() { t.Put(2, 'c') }, ShouldPanic)
					So(func() { t.Put(3, 3.14) }, ShouldPanic)
					So(func() { t.Put(4, 'c') }, ShouldPanic)
					So(func() { t.Put(5, 123) }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New5(42, 3.14, io.EOF, true, 'c'))
				})

				Convey("at index 1", func() {
					So(t.Del(1), ShouldEqual, New5("hello", 3.14, io.EOF, true, 'c'))
				})

				Convey("at index2", func() {
					So(t.Del(2), ShouldEqual, New5("hello", 42, io.EOF, true, 'c'))
				})

				Convey("at index3", func() {
					So(t.Del(3), ShouldEqual, New5("hello", 42, 3.14, true, 'c'))
				})

				Convey("at index4", func() {
					So(t.Del(4), ShouldEqual, New5("hello", 42, 3.14, io.EOF, 'c'))
				})

				Convey("at index5", func() {
					So(t.Del(5), ShouldEqual, New5("hello", 42, 3.14, io.EOF, true))
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})

		Convey("When create Tuple7", func() {
			t := New7("hello", 42, 3.14, io.EOF, true, 'c', []string{"foo", "bar"})

			So(t.String(), ShouldEqual, "(hello, 42, 3.14, EOF, true, 99, [foo bar])")
			So(t.Len(), ShouldEqual, 7)

			Convey("Then unpack the tuple", func() {
				v0, v1, v3, v4, v5, v6, v7 := t.Unpack()
				So(v0, ShouldEqual, "hello")
				So(v1, ShouldEqual, 42)
				So(v3, ShouldEqual, 3.14)
				So(v4, ShouldEqual, io.EOF)
				So(v5, ShouldBeTrue)
				So(v6, ShouldAlmostEqual, 'c')
				So(v7, ShouldEqual, []string{"foo", "bar"})
			})

			Convey("Then get value", func() {
				Convey("at index", func() {
					So(t.Get(0), ShouldEqual, "hello")
					So(t.Get(1), ShouldEqual, 42)
					So(t.Get(2), ShouldEqual, 3.14)
					So(t.Get(3), ShouldEqual, io.EOF)
					So(t.Get(4), ShouldEqual, true)
					So(t.Get(5), ShouldEqual, 'c')
					So(t.Get(6), ShouldEqual, []string{"foo", "bar"})
				})

				Convey("index out of range", func() {
					So(func() { t.Get(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then put a new value", func() {
				Convey("at index 0", func() {
					new, old := t.Put(0, "foobar")
					So(new.String(), ShouldEqual, "(foobar, 42, 3.14, EOF, true, 99, [foo bar])")
					So(old, ShouldEqual, "hello")
				})

				Convey("at index 1", func() {
					new, old := t.Put(1, 123)
					So(new.String(), ShouldEqual, "(hello, 123, 3.14, EOF, true, 99, [foo bar])")
					So(old, ShouldEqual, 42)
				})

				Convey("at index 2", func() {
					new, old := t.Put(2, 1.23)
					So(new.String(), ShouldEqual, "(hello, 42, 1.23, EOF, true, 99, [foo bar])")
					So(old, ShouldEqual, 3.14)
				})

				Convey("at index 3", func() {
					new, old := t.Put(3, io.ErrUnexpectedEOF)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, unexpected EOF, true, 99, [foo bar])")
					So(old, ShouldEqual, io.EOF)
				})

				Convey("at index 4", func() {
					new, old := t.Put(4, false)
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, EOF, false, 99, [foo bar])")
					So(old, ShouldBeTrue)
				})

				Convey("at index 5", func() {
					new, old := t.Put(5, 'b')
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, EOF, true, 98, [foo bar])")
					So(old, ShouldEqual, 'c')
				})

				Convey("at index 6", func() {
					new, old := t.Put(6, []string{"foobar"})
					So(new.String(), ShouldEqual, "(hello, 42, 3.14, EOF, true, 99, [foobar])")
					So(old, ShouldEqual, []string{"foo", "bar"})
				})

				Convey("with wrong type", func() {
					So(func() { t.Put(0, 123) }, ShouldPanic)
					So(func() { t.Put(1, false) }, ShouldPanic)
					So(func() { t.Put(2, 'c') }, ShouldPanic)
					So(func() { t.Put(3, 3.14) }, ShouldPanic)
					So(func() { t.Put(4, 'c') }, ShouldPanic)
					So(func() { t.Put(5, 123) }, ShouldPanic)
					So(func() { t.Put(6, 123) }, ShouldPanic)
				})

				Convey("index out of range", func() {
					So(func() { t.Put(t.Len(), "foobar") }, ShouldPanicWith, ErrOutOfRange)
				})
			})

			Convey("Then delete a value", func() {
				Convey("at index 0", func() {
					So(t.Del(0), ShouldEqual, New6(42, 3.14, io.EOF, true, 'c', []string{"foo", "bar"}))
				})

				Convey("at index 1", func() {
					So(t.Del(1), ShouldEqual, New6("hello", 3.14, io.EOF, true, 'c', []string{"foo", "bar"}))
				})

				Convey("at index2", func() {
					So(t.Del(2), ShouldEqual, New6("hello", 42, io.EOF, true, 'c', []string{"foo", "bar"}))
				})

				Convey("at index3", func() {
					So(t.Del(3), ShouldEqual, New6("hello", 42, 3.14, true, 'c', []string{"foo", "bar"}))
				})

				Convey("at index4", func() {
					So(t.Del(4), ShouldEqual, New6("hello", 42, 3.14, io.EOF, 'c', []string{"foo", "bar"}))
				})

				Convey("at index5", func() {
					So(t.Del(5), ShouldEqual, New6("hello", 42, 3.14, io.EOF, true, []string{"foo", "bar"}))
				})

				Convey("at index6", func() {
					So(t.Del(6), ShouldEqual, New6("hello", 42, 3.14, io.EOF, true, 'c'))
				})

				Convey("index out of range", func() {
					So(func() { t.Del(t.Len()) }, ShouldPanicWith, ErrOutOfRange)
				})
			})
		})
	})
}
