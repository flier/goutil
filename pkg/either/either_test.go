package either_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/either"
)

func TestEither(t *testing.T) {
	Convey("Given some Either values", t, func() {
		l := Left[string, int]("error")
		r := Right[string](42)

		Convey("When checking HasLeft/HasRight", func() {
			So(l.HasLeft(), ShouldBeTrue)
			So(l.HasRight(), ShouldBeFalse)
			So(r.HasLeft(), ShouldBeFalse)
			So(r.HasRight(), ShouldBeTrue)
		})

		Convey("When unwrapping values", func() {
			So(l.UnwrapLeft(), ShouldEqual, "error")
			So(r.UnwrapRight(), ShouldEqual, 42)
		})

		Convey("When mapping values", func() {
			doubled := MapRight(r, func(x int) int { return x * 2 })
			So(doubled.UnwrapRight(), ShouldEqual, 84)

			uppercased := MapLeft(l, func(s string) string { return "ERROR: " + s })
			So(uppercased.UnwrapLeft(), ShouldEqual, "ERROR: error")
		})

		Convey("When converting to string", func() {
			So(fmt.Sprint(Empty[int, int]()), ShouldEqual, "Empty")
			So(fmt.Sprint(l), ShouldEqual, "Left(error)")
			So(fmt.Sprint(r), ShouldEqual, "Right(42)")
		})

		Convey("When converting to go string", func() {
			So(fmt.Sprintf("%#v", Empty[int, int]()), ShouldEqual, "Either {}")
			So(fmt.Sprintf("%#v", l), ShouldEqual, "Either { Left: error }")
			So(fmt.Sprintf("%#v", r), ShouldEqual, "Either { Right: 42 }")
		})

		Convey("When swapping sides", func() {
			s := r.Flip()
			So(s.HasLeft(), ShouldBeTrue)
			So(s.UnwrapLeft(), ShouldEqual, 42)
		})
	})
}

func TestEitherCombinators(t *testing.T) {
	Convey("Given Either combinators", t, func() {
		Convey("When using RightOr", func() {
			l := Left[string, int]("error")
			r := Right[string](42)

			So(l.RightOr(123), ShouldEqual, 123)
			So(r.RightOr(123), ShouldEqual, 42)
		})

		Convey("When using RightOrElse", func() {
			l := Left[string, int]("error")
			r := Right[string](42)

			f := func() int { return 999 }

			So(l.RightOrElse(f), ShouldEqual, 999)
			So(r.RightOrElse(f), ShouldEqual, 42)
		})

		Convey("When using AndThen", func() {
			r := Right[string](42)
			f := func(x int) Either[string, string] {
				return Right[string, string](fmt.Sprintf("%d", x))
			}

			result := RightAndThen(r, f)
			So(result.UnwrapRight(), ShouldEqual, "42")
		})
	})
}

func ExampleEither_ExpectLeft() {
	defer func() { fmt.Println(recover()) }() // fruits are healthy: test

	fmt.Println(Left[int, string](123).ExpectLeft("fruits are healthy")) // 123
	fmt.Println(Right[int]("test").ExpectLeft("fruits are healthy"))     // panic!

	// Output:
	// 123
	// fruits are healthy: test
}

func ExampleEither_ExpectRight() {
	defer func() { fmt.Println(recover()) }() // fruits are healthy: error

	fmt.Println(Right[string]("value").ExpectRight("fruits are healthy"))        // value
	fmt.Println(Left[string, string]("error").ExpectRight("fruits are healthy")) // panic!

	// Output:
	// value
	// fruits are healthy: error
}

func ExampleEither_Flip() {
	fmt.Println(Left[int, string](123).Flip())           // Right(123)
	fmt.Println(Right[int]("hello").Flip())              // Left(hello)
	fmt.Println(Left[bool, float64](true).Flip())        // Right(true)
	fmt.Println(Right[bool](3.14).Flip())                // Left(3.14)
	fmt.Println(Left[[]int, string]([]int{1, 2}).Flip()) // Right([1 2])

	// Output:
	// Right(123)
	// Left(hello)
	// Right(true)
	// Left(3.14)
	// Right([1 2])
}

func ExampleEither_LeftOr() {
	fmt.Println(Left[int, string](5).LeftOr(10)) // 5
	fmt.Println(Right[int]("foo").LeftOr(10))    // 10

	// Output:
	// 5
	// 10
}

func ExampleEither_LeftOrEmpty() {
	fmt.Println(Left[int, string](42).LeftOrEmpty())      // 42
	fmt.Println(Right[int]("hello").LeftOrEmpty())        // 0
	fmt.Println(Left[string, bool]("test").LeftOrEmpty()) // test
	fmt.Println(Right[string](true).LeftOrEmpty())        // ""
	fmt.Println(Left[float64, int](3.14).LeftOrEmpty())   // 3.14
	fmt.Println(Right[float64](1).LeftOrEmpty())          // 0

	// Output:
	// 42
	// 0
	// test
	//
	// 3.14
	// 0
}

func ExampleEither_LeftOrElse() {
	nobody := func() string { return "nobody" }
	vikings := func() string { return "vikings" }

	fmt.Println(Left[string, int]("barbarians").LeftOrElse(vikings))
	fmt.Println(Right[string](2).LeftOrElse(vikings))

	fmt.Println(Left[string, int]("barbarians").LeftOrElse(nobody))
	fmt.Println(Right[string](2).LeftOrElse(nobody))

	// Output:
	// barbarians
	// vikings
	// barbarians
	// nobody
}

func ExampleEither_RightOr() {
	fmt.Println(Right[string, int](5).RightOr(2))    // 5
	fmt.Println(Left[string, int]("foo").RightOr(2)) // 2

	// Output:
	// 5
	// 2
}

func ExampleEither_RightOrEmpty() {
	fmt.Println(Right[string](42).RightOrEmpty())          // 42
	fmt.Println(Left[string, int]("error").RightOrEmpty()) // 0
	fmt.Println(Right[string](123).RightOrEmpty())         // 123
	fmt.Println(Left[string, int]("fail").RightOrEmpty())  // 0
	fmt.Println(Right[string](0).RightOrEmpty())           // 0

	// Output:
	// 42
	// 0
	// 123
	// 0
	// 0
}

func ExampleEither_RightOrElse() {
	k := 21

	fmt.Println(Right[string](2).RightOrElse(func() int { return 2 * k }))         // 2
	fmt.Println(Left[string, int]("foo").RightOrElse(func() int { return 2 * k })) // 42

	// Output:
	// 2
	// 42
}

func ExampleMapLeft() {
	fmt.Println(MapLeft(Left[int, string](5), func(x int) float64 { return float64(x) * 2.5 })) // Left(12.5)
	fmt.Println(MapLeft(Right[int]("hello"), func(x int) float64 { return float64(x) }))        // Right(hello)

	// Output:
	// Left(12.5)
	// Right(hello)
}

func ExampleMapRight() {
	fmt.Println(MapRight(Left[int, string](42), strings.ToUpper)) // Left(42)
	fmt.Println(MapRight(Right[int]("foo"), strings.ToUpper))     // Right(FOO)

	// Output:
	// Left(42)
	// Right(FOO)
}

func ExampleMapEither() {
	pow := func(n int) float64 { return float64(n * n) }
	fmt.Println(MapEither(Left[int, string](42), pow, strings.ToUpper)) // Left(1764)
	fmt.Println(MapEither(Right[int]("hello"), pow, strings.ToUpper))   // Right(HELLO)

	// Output:
	// Left(1764)
	// Right(HELLO)
}

func ExampleReduce() {
	square := func(n int) int { return n * n }
	negate := func(n int) int { return -n }

	fmt.Println(Reduce(Left[int, int](4), square, negate)) // 16
	fmt.Println(Reduce(Right[int](-4), square, negate))    // 4

	// Output:
	// 16
	// 4
}

func ExampleLeftAndThen() {
	f := func(s string) Either[string, int] {
		return Left[string, int](s + "!")
	}

	fmt.Println(LeftAndThen(Left[string, int]("hello"), f)) // Left(hello!)
	fmt.Println(LeftAndThen(Right[string, int](42), f))     // Right(42)

	// Output:
	// Left(hello!)
	// Right(42)
}

func ExampleRightAndThen() {
	f := func(n int) Either[string, string] {
		return Right[string](strconv.Itoa(n * n))
	}

	fmt.Println(RightAndThen(Right[string](2), f))           // Right(4)
	fmt.Println(RightAndThen(Left[string, int]("error"), f)) // Left(error)

	// Output:
	// Right(4)
	// Left(error)
}
