package guard

import (
	"fmt"
	"runtime"
)

func GetCallstack() string {
	const stackSize = 1024

	// Create a slice to hold the PC addresses
	pc := make([]uintptr, stackSize)

	// Capture the PC addresses for the current goroutine's call stack
	n := runtime.Callers(0, pc)

	// Create a frames slice to hold the frames
	frames := runtime.CallersFrames(pc[:n])

	// Create a string to store the call stack
	callStack := ""

	// Iterate over the frames and append them to the call stack string
	readFrames := 0
	for {
		frame, more := frames.Next()
		readFrames++

		if readFrames > 2 {
			callStack += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		}

		if !more {
			break
		}
	}

	return callStack
}
