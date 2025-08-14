package debug

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Stack is like [runtime/debug.Stack], but with a skip parameter and an
// easier to read format.
func Stack(skip int) string {
	var out strings.Builder

	trace := make([]uintptr, 32)
	for {
		n := runtime.Callers(skip, trace)
		if n < len(trace) {
			trace = trace[:n]
			break
		}
		trace = make([]uintptr, len(trace)*2)
	}

	frames := runtime.CallersFrames(trace)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&out,
			"- %-24v 0x%x+0x%-4x %v:%v\n",
			path.Base(frame.Function)+"()", frame.Entry, frame.PC-frame.Entry,
			frame.File, frame.Line,
		)

		if !more {
			break
		}
	}

	return out.String()
}
