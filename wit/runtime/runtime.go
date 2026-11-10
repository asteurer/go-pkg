package runtime

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Handle struct {
	value int32
}

func (h *Handle) Use() int32 {
	if h.value == 0 {
		panic("nil handle")
	}
	return h.value
}

func (h *Handle) Take() int32 {
	if h.value == 0 {
		panic("nil handle")
	}
	value := h.value
	h.value = 0
	return value
}

func (h *Handle) Set(value int32) {
	if value == 0 {
		panic("nil handle")
	}
	if h.value != 0 {
		panic("handle already set")
	}
	h.value = value
}

func (h *Handle) TakeOrNil() int32 {
	value := h.value
	h.value = 0
	return value
}

func MakeHandle(value int32) *Handle {
	if value == 0 {
		panic("nil handle")
	}
	return &Handle{value}
}

func Allocate(pinner *runtime.Pinner, size, align uintptr) unsafe.Pointer {
	pointer := allocateRaw(size, align)
	pinner.Pin(pointer)
	return pointer
}

func allocateRaw(size, align uintptr) unsafe.Pointer {
	if size == 0 {
		return nil
	}

	if size%align != 0 {
		panic(fmt.Sprintf("size %v is not compatible with alignment %v", size, align))
	}

	switch align {
	case 1:
		return unsafe.Pointer(unsafe.SliceData(make([]uint8, size)))
	case 2:
		return unsafe.Pointer(unsafe.SliceData(make([]uint16, size/align)))
	case 4:
		return unsafe.Pointer(unsafe.SliceData(make([]uint32, size/align)))
	case 8:
		return unsafe.Pointer(unsafe.SliceData(make([]uint64, size/align)))
	default:
		panic(fmt.Sprintf("unsupported alignment: %v", align))
	}
}

var pinner = runtime.Pinner{}

func Unpin() {
	pinner.Unpin()
}
