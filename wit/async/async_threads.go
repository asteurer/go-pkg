//go:build tinygo

package async

import (
	"fmt"
	"runtime"
	"unsafe"
)

// This is the TinyGo implementation. Exports are stackful and every goroutine
// is a thread that blocks on its own waitable set. A thread that blocks in
// waitable-set.wait does not stop the other threads. There is no callback
// loop, no task state and no scheduler idle hook.

// Run runs the body of an async export on the calling thread. After the body
// has returned its result, it keeps the task alive until the goroutines that
// the body started have finished, so a handler can stream a response body
// from a goroutine.
func Run(closure func()) uint32 {
	base := runtime.NumGoroutine()
	closure()
	for runtime.NumGoroutine() > base {
		runtime.Gosched()
	}
	return CALLBACK_CODE_EXIT
}

// Callback is not used. Exports are stackful and have no callback.
func Callback(event0, event1, event2 uint32) uint32 {
	panic("async: callback exports are not used")
}

// waitFor blocks the calling thread until the waitable has an event and
// returns the event and its payload.
func waitFor(handle uint32) (event, payload uint32) {
	set := waitableSetNew()
	waitableJoin(handle, set)
	var buf [2]uint32
	pinner := runtime.Pinner{}
	pinner.Pin(&buf)
	event = waitableSetWait(set, unsafe.Pointer(&buf[0]))
	pinner.Unpin()
	waitableJoin(handle, 0)
	waitableSetDrop(set)
	return event, buf[1]
}

func SubtaskWait(status uint32) {
	subtask := status >> 4
	status = status & 0xF

	switch status {
	case STATUS_STARTING, STATUS_STARTED:
		event, payload := waitFor(subtask)
		if event != EVENT_SUBTASK || payload != STATUS_RETURNED {
			panic(fmt.Sprintf("unexpected subtask event %v payload %v", event, payload))
		}
		subtaskDrop(subtask)

	case STATUS_RETURNED:

	default:
		panic(fmt.Sprintf("unexpected subtask status: %v", status))
	}
}

func FutureOrStreamWait(code uint32, handle int32) (uint32, uint32) {
	if code == RETURN_CODE_BLOCKED {
		_, code = waitFor(uint32(handle))
	}

	count := code >> 4
	code = code & 0xF

	return code, count
}

func Yield() {
	runtime.Gosched()
}

//go:wasmimport $root [waitable-set-new]
//go:noescape
func waitableSetNew() uint32

//go:wasmimport $root [waitable-set-wait]
//go:noescape
func waitableSetWait(waitableSet uint32, eventPayload unsafe.Pointer) uint32

//go:wasmimport $root [waitable-set-drop]
//go:noescape
func waitableSetDrop(waitableSet uint32)

//go:wasmimport $root [waitable-join]
//go:noescape
func waitableJoin(waitable, waitableSet uint32)

//go:wasmimport $root [subtask-drop]
//go:noescape
func subtaskDrop(subtask uint32)
