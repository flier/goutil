//go:build go1.23

package xiter_test

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleResource() {
	const poem = `Nature's first green is gold,
	Her hardest hue to hold.
	Her early leaf's a flower;
	But only so an hour.`

	f := io.NopCloser(strings.NewReader(poem))

	lines := Resource(
		func() (r *bufio.Reader, err error) {
			return bufio.NewReader(f), nil
		}, func(r *bufio.Reader) (line string, err error) {
			line, err = r.ReadString('\n')

			if line = strings.TrimSpace(line); len(line) > 0 && err == io.EOF {
				err = nil
			}

			return
		}, func(r *bufio.Reader) {
			_ = f.Close()
		})

	for line := range lines {
		fmt.Println(line)
	}

	// Output:
	// Nature's first green is gold,
	// Her hardest hue to hold.
	// Her early leaf's a flower;
	// But only so an hour.
}

func ExampleLines() {
	const poem = `Nature's first green is gold,
	Her hardest hue to hold.
	Her early leaf's a flower;
	But only so an hour.`

	f := io.NopCloser(strings.NewReader(poem))

	for line := range Lines(f) {
		fmt.Println(strings.TrimSpace(line))
	}

	// Output:
	// Nature's first green is gold,
	// Her hardest hue to hold.
	// Her early leaf's a flower;
	// But only so an hour.
}

func TestResource(t *testing.T) {
	Convey("Resource", t, func() {
		Convey("Should iterate over values from a resource", func() {
			counter := 0
			values := []string{"value_0", "value_1", "value_2"}
			index := 0

			start := func() (int, error) { return 0, nil }
			next := func(s int) (string, error) {
				if index >= len(values) {
					return "", io.EOF
				}
				result := values[index]
				index++
				return result, nil
			}
			stop := func(s int) { counter++ }

			seq := Resource(start, next, stop)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldResemble, []string{"value_0", "value_1", "value_2"})
			So(counter, ShouldEqual, 1) // stop should be called once
		})

		Convey("Should handle start function error", func() {
			start := func() (int, error) { return 0, fmt.Errorf("start error") }
			next := func(s int) (string, error) { return "value", nil }
			stop := func(s int) {}

			seq := Resource(start, next, stop)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldBeEmpty)
		})

		Convey("Should handle empty resource", func() {
			counter := 0
			start := func() (int, error) { return 0, nil }
			next := func(s int) (string, error) { return "", io.EOF }
			stop := func(s int) { counter++ }

			seq := Resource(start, next, stop)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldBeEmpty)
			So(counter, ShouldEqual, 1) // stop should be called once
		})

		Convey("Should handle early termination", func() {
			counter := 0
			values := []string{"value_0", "value_1", "value_2", "value_3", "value_4"}
			index := 0

			start := func() (int, error) { return 0, nil }
			next := func(s int) (string, error) {
				if index >= len(values) {
					return "", io.EOF
				}
				result := values[index]
				index++
				return result, nil
			}
			stop := func(s int) { counter++ }

			seq := Resource(start, next, stop)
			result := make([]string, 0)
			count := 0
			for v := range seq {
				result = append(result, v)
				count++
				if count >= 2 { // Early termination
					break
				}
			}

			So(result, ShouldResemble, []string{"value_0", "value_1"})
			So(counter, ShouldEqual, 1) // stop should be called once
		})
	})
}

func TestLines(t *testing.T) {
	Convey("Lines", t, func() {
		Convey("Should iterate over lines from io.ReadCloser", func() {
			content := "line1\nline2\nline3\n"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldResemble, []string{"line1", "line2", "line3"})
		})

		Convey("Should handle single line", func() {
			content := "single line"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldResemble, []string{"single line"})
		})

		Convey("Should handle empty content", func() {
			content := ""
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldBeEmpty)
		})

		Convey("Should handle content with only newlines", func() {
			content := "\n\n\n"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldResemble, []string{"", "", ""})
		})

		Convey("Should handle content without trailing newline", func() {
			content := "line1\nline2\nline3"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldResemble, []string{"line1", "line2", "line3"})
		})

		Convey("Should handle early termination", func() {
			content := "line1\nline2\nline3\nline4\nline5"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			count := 0
			for v := range seq {
				result = append(result, v)
				count++
				if count >= 2 { // Early termination
					break
				}
			}

			So(result, ShouldResemble, []string{"line1", "line2"})
		})

		Convey("Should handle large content", func() {
			lines := make([]string, 1000)
			for i := range lines {
				lines[i] = fmt.Sprintf("line_%d", i)
			}
			content := strings.Join(lines, "\n") + "\n"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(len(result), ShouldEqual, 1000)
			So(result[0], ShouldEqual, "line_0")
			So(result[999], ShouldEqual, "line_999")
		})

		Convey("Should handle special characters in lines", func() {
			content := "line with spaces\nline\twith\ttabs\nline\nwith\nmultiple\nnewlines\n"
			reader := io.NopCloser(strings.NewReader(content))

			seq := Lines(reader)
			result := make([]string, 0)
			for v := range seq {
				result = append(result, v)
			}

			So(result, ShouldResemble, []string{
				"line with spaces",
				"line\twith\ttabs",
				"line",
				"with",
				"multiple",
				"newlines",
			})
		})
	})
}
