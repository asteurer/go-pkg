//go:build !tinygo

package runtime

import (
	"unsafe"
)

// NB: `cabi_realloc` may be called before the Go runtime has been initialized,
// in which case we need to use `runtime.sbrk` to do allocations.  The following
// is an abbreviation of [Till's
// efforts](https://github.com/bytecodealliance/go-modules/pull/367).

//go:linkname sbrk runtime.sbrk
func sbrk(n uintptr) unsafe.Pointer

//nolint:unused
var useGCAllocations = false

func init() {
	useGCAllocations = true
}

//nolint:unused
func offset(ptr, align uintptr) uintptr {
	newptr := (ptr + align - 1) &^ (align - 1)
	return newptr - ptr
}

//nolint:unused
//go:wasmimport wasi_snapshot_preview1 adapter_monotonic_clock_set_paused
func adapterMonotonicClockSetPaused(paused bool)

// `procPin` increments the current M's lock count, which prevents the Go
// runtime from starting a GC cycle or performing GC assist work until
// `procUnpin` is called. The runtime explicitly supports linking to these;
// see https://go.dev/issue/67401.

//nolint:unused
//go:linkname procPin runtime.procPin
func procPin() int

//nolint:unused
//go:linkname procUnpin runtime.procUnpin
func procUnpin()

//nolint:unused
//go:wasmexport cabi_realloc
func cabiRealloc(oldPointer unsafe.Pointer, oldSize, align, newSize uintptr) unsafe.Pointer {
	if oldPointer != nil || oldSize != 0 {
		panic("todo")
	}

	if useGCAllocations {
		// Calls to imports from `cabi_realloc` are forbidden by the
		// component model, so we must not let the Go garbage collector
		// run here. A GC cycle may call imports in several ways: it
		// reads the monotonic and realtime clocks, polls for I/O when
		// restarting the world, and may switch to other goroutines
		// (e.g. to park an assist), any of which may then call imports.
		// Pinning the M defers any GC work until the next allocation
		// made outside of `cabi_realloc`.
		//
		// We additionally call `adapter_monotonic_clock_set_paused`
		// before and after allocating in case anything else reads the
		// monotonic clock while allocating.  See
		// https://github.com/bytecodealliance/wasmtime/pull/13563 for
		// details.
		procPin()
		adapterMonotonicClockSetPaused(true)
		pointer := Allocate(&pinner, newSize, align)
		adapterMonotonicClockSetPaused(false)
		procUnpin()
		return pointer
	} else {
		alignedSize := newSize + offset(newSize, align)
		unaligned := sbrk(alignedSize)
		off := offset(uintptr(unaligned), align)
		return unsafe.Add(unaligned, off)
	}
}
