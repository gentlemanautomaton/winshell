// +build windows

package shobjidl

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// ISequentialStreamVtbl represents the component object model virtual
// function table for the ISequentialStream interface.
type ISequentialStreamVtbl struct {
	ole.IUnknownVtbl
	Read  uintptr
	Write uintptr
}

// ISequentialStream represents the component object model interface for
// marshaling binary representations of objects.
type ISequentialStream struct {
	ole.IUnknown
}

// VTable returns the component object model virtual function table for the
// object.
func (v *ISequentialStream) VTable() *ISequentialStreamVtbl {
	return (*ISequentialStreamVtbl)(unsafe.Pointer(v.RawVTable))
}
