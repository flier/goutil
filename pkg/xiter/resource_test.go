//go:build go1.23

package xiter_test

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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
